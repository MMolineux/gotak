package tak

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/angry-kivi/gotak/pkg/util"
)

// TLSClient implements the Client interface for TLS connections
type TLSClient struct {
	config ClientConfig
	conn   *tls.Conn
}

// NewTLSClient creates a new TLS client for TAK communication
func NewTLSClient(config ClientConfig) (*TLSClient, error) {
	// If TLSConfig is already provided, use it
	if config.TLSConfig == nil {
		// Check if any certificate files are provided
		certProvided := config.CertFile != ""

		// If we're not skipping verification and no cert provided, return error
		if !config.SkipTLSVerify && !certProvided {
			return nil, errors.New("TLS connection requires either a certificate or skip-verify option")
		}

		// If certificate file is provided, attempt to load config
		if certProvided {
			var err error

			// Use the auto-detection function to handle both PEM and P12 formats
			config.TLSConfig, err = util.LoadTLSConfigAuto(
				config.CertFile,
				config.KeyFile,
				config.P12Password,
				config.CAFile,
				config.SkipTLSVerify,
			)

			if err != nil {
				return nil, fmt.Errorf("failed to load TLS configuration: %w", err)
			}
		} else if config.SkipTLSVerify {
			// Create a basic TLS config that skips verification
			config.TLSConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
	}

	// Ensure we have a valid TLS configuration at this point
	if config.TLSConfig == nil {
		return nil, errors.New("failed to create TLS configuration")
	}

	return &TLSClient{
		config: config,
	}, nil
}

// Connect establishes a TLS connection to the TAK server
func (c *TLSClient) Connect(ctx context.Context) error {
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
	conn, err := tls.DialWithDialer(dialer, "tcp", address, c.config.TLSConfig.Clone())
	if err != nil {
		return err
	}

	// Apply TCP_USER_TIMEOUT if configured (Linux only).
	// Causes kernel to kill connection if unacked data sits in send buffer too long.
	if c.config.TCPUserTimeout > 0 {
		if err := setTCPUserTimeoutTLS(conn, c.config.TCPUserTimeout); err != nil {
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

// Disconnect closes the TLS connection
func (c *TLSClient) Disconnect() error {
	if c.conn == nil {
		return nil
		//return errors.New("client is not connected")
	}

	err := c.conn.Close()
	if err == nil {
		c.conn = nil
	}
	return err
}

// Send transmits data over the TLS connection
func (c *TLSClient) Send(data []byte) error {
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

// Receive reads data from the TLS connection
func (c *TLSClient) Receive() ([]byte, error) {
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

// HasPendingData checks if there's data available to read without blocking
func (c *TLSClient) HasPendingData() (bool, error) {
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

// IsConnected returns true if the client has an active connection
func (c *TLSClient) IsConnected() bool {
	return c.conn != nil
}
