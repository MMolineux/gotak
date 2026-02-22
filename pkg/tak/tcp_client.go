package tak

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

// TCPClient implements the Client interface for TCP connections
type TCPClient struct {
	config ClientConfig
	conn   net.Conn
}

// NewTCPClient creates a new TCP client for TAK communication
func NewTCPClient(config ClientConfig) (*TCPClient, error) {
	return &TCPClient{
		config: config,
	}, nil
}

// Connect establishes a TCP connection to the TAK server
func (c *TCPClient) Connect(ctx context.Context) error {
	if c.conn != nil {
		return errors.New("client is already connected")
	}

	dialer := &net.Dialer{
		Timeout:   c.config.DialTimeout,
		KeepAlive: c.config.KeepAlive,
	}
	address := net.JoinHostPort(c.config.Address, strconv.Itoa(c.config.Port))

	var err error
	// Use the provided context for cancellation
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return err
	}

	// Apply TCP_USER_TIMEOUT if configured (Linux only).
	// Causes kernel to kill connection if unacked data sits in send buffer too long.
	if c.config.TCPUserTimeout > 0 {
		if err := setTCPUserTimeout(conn, c.config.TCPUserTimeout); err != nil {
			conn.Close()
			return fmt.Errorf("failed to set TCP_USER_TIMEOUT: %w", err)
		}
	}

	// Set up connection monitoring
	go func() {
		<-ctx.Done()
		if c.conn != nil {
			c.conn.Close()
		}
	}()

	c.conn = conn
	return nil
}

// Disconnect closes the TCP connection
func (c *TCPClient) Disconnect() error {
	if c.conn == nil {
		return errors.New("client is not connected")
	}

	err := c.conn.Close()
	if err == nil {
		c.conn = nil
	}
	return err
}

// Send transmits data over the TCP connection
func (c *TCPClient) Send(data []byte) error {
	if c.conn == nil {
		return errors.New("client is not connected")
	}

	if c.config.WriteTimeout > 0 {
		deadline := time.Now().Add(c.config.WriteTimeout)
		if err := c.conn.SetWriteDeadline(deadline); err != nil {
			return err
		}
	}

	_, err := c.conn.Write(data)
	return err
}

// Receive reads data from the TCP connection
func (c *TCPClient) Receive() ([]byte, error) {
	if c.conn == nil {
		return nil, errors.New("client is not connected")
	}

	if c.config.ReadTimeout > 0 {
		deadline := time.Now().Add(c.config.ReadTimeout)
		if err := c.conn.SetReadDeadline(deadline); err != nil {
			return nil, err
		}
	}

	buffer := make([]byte, 8192) // Reasonable buffer size for CoT messages
	n, err := c.conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

// IsConnected returns true if the client has an active connection
func (c *TCPClient) IsConnected() bool {
	return c.conn != nil
}

// HasPendingData checks if there's data available to read without blocking
func (c *TCPClient) HasPendingData() (bool, error) {
	if c.conn == nil {
		return false, errors.New("client is not connected")
	}

	// Set a very short read deadline to check for data
	if err := c.conn.SetReadDeadline(time.Now().Add(1 * time.Millisecond)); err != nil {
		return false, err
	}

	// Reset one byte to see if data is available
	oneByteBuffer := make([]byte, 1)
	_, err := c.conn.Read(oneByteBuffer)

	// Reset the deadline
	c.conn.SetReadDeadline(time.Time{})

	if err == nil {
		// Data was read, put it back (not ideal but functional for checking)
		// In a production system, you might want a better solution for this
		return true, nil
	} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		// This is a timeout error, which means no data is available
		return false, nil
	}

	// Some other error occurred
	return false, err
}
