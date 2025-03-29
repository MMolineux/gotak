package tak

import (
	"context"
	"errors"
	"net"
	"strconv"
	"time"
)

// UDPClient implements the Client interface for UDP connections
type UDPClient struct {
	config     ClientConfig
	conn       *net.UDPConn
	serverAddr *net.UDPAddr
}

// NewUDPClient creates a new UDP client for TAK communication
func NewUDPClient(config ClientConfig) (*UDPClient, error) {
	return &UDPClient{
		config: config,
	}, nil
}

// Connect establishes a UDP connection for communication with the TAK server
func (c *UDPClient) Connect(ctx context.Context) error {
	if c.conn != nil {
		return errors.New("client is already connected")
	}

	// Resolve the server address
	var err error
	c.serverAddr, err = net.ResolveUDPAddr("udp", net.JoinHostPort(c.config.Address,
		strconv.Itoa(c.config.Port)))
	if err != nil {
		return err
	}

	// Create a local address to listen on
	localAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	if err != nil {
		return err
	}

	// Create the UDP connection
	c.conn, err = net.DialUDP("udp", localAddr, c.serverAddr)
	if err != nil {
		return err
	}

	// Set the read and write deadlines if configured
	if c.config.ReadTimeout > 0 {
		err = c.conn.SetReadDeadline(time.Now().Add(c.config.ReadTimeout))
		if err != nil {
			return err
		}
	}

	if c.config.WriteTimeout > 0 {
		err = c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout))
		if err != nil {
			return err
		}
	}

	return nil
}

// Disconnect closes the UDP connection
func (c *UDPClient) Disconnect() error {
	if c.conn == nil {
		return errors.New("client is not connected")
	}

	err := c.conn.Close()
	if err == nil {
		c.conn = nil
		c.serverAddr = nil
	}
	return err
}

// Send transmits data over the UDP connection
func (c *UDPClient) Send(data []byte) error {
	if c.conn == nil || c.serverAddr == nil {
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

// Receive reads data from the UDP connection
func (c *UDPClient) Receive() ([]byte, error) {
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
	n, _, err := c.conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

// IsConnected returns true if the client has an active connection
func (c *UDPClient) IsConnected() bool {
	return c.conn != nil && c.serverAddr != nil
}
