package tak

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTCPClient_Connect(t *testing.T) {
	// Start a mock TCP server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	_, portStr, err := net.SplitHostPort(listener.Addr().String())
	require.NoError(t, err)
	port, err := strconv.Atoi(portStr)
	require.NoError(t, err)

	// Accept connections in the background
	connCh := make(chan net.Conn, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		connCh <- conn
	}()

	// Create a client with the test server address
	config := ClientConfig{
		Address:        "127.0.0.1",
		Port:           port,
		ConnectionType: ConnectionTypeTCP,
		DialTimeout:    time.Second * 5,
		Logger:         logrus.New(),
	}

	client, err := NewTCPClient(config)
	require.NoError(t, err)

	// Test Connect
	ctx := context.Background()
	err = client.Connect(ctx)
	require.NoError(t, err)
	assert.True(t, client.IsConnected())

	// Close the connection
	err = client.Disconnect()
	require.NoError(t, err)
	assert.False(t, client.IsConnected())
}

func TestTCPClient_SendReceive(t *testing.T) {
	// Start a mock TCP server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	_, portStr, err := net.SplitHostPort(listener.Addr().String())
	require.NoError(t, err)
	port, err := strconv.Atoi(portStr)
	require.NoError(t, err)

	// Handle connections and respond to client data
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		// Echo the data back
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		_, err = conn.Write(buf[:n])
		if err != nil {
			return
		}
	}()

	// Create a client with the test server address
	config := ClientConfig{
		Address:        "127.0.0.1",
		Port:           port,
		ConnectionType: ConnectionTypeTCP,
		DialTimeout:    time.Second * 5,
		Logger:         logrus.New(),
	}

	client, err := NewTCPClient(config)
	require.NoError(t, err)

	// Connect to the server
	ctx := context.Background()
	err = client.Connect(ctx)
	require.NoError(t, err)
	defer client.Disconnect()

	// Test Send
	testData := []byte("Hello, TAK Server!")
	err = client.Send(testData)
	require.NoError(t, err)

	// Test Receive
	receivedData, err := client.Receive()
	require.NoError(t, err)
	assert.Equal(t, testData, receivedData)
}

func TestTCPClient_ConnectError(t *testing.T) {
	// Create a client with an invalid address
	config := ClientConfig{
		Address:        "non.existent.host",
		Port:           12345,
		ConnectionType: ConnectionTypeTCP,
		DialTimeout:    time.Second * 1, // Short timeout for faster test
		Logger:         logrus.New(),
	}

	client, err := NewTCPClient(config)
	require.NoError(t, err)

	// Test Connect - should fail
	ctx := context.Background()
	err = client.Connect(ctx)
	assert.Error(t, err)
	assert.False(t, client.IsConnected())
}

func TestTCPClient_SendWithoutConnect(t *testing.T) {
	// Create a client without connecting
	config := ClientConfig{
		Address:        "127.0.0.1",
		Port:           12345,
		ConnectionType: ConnectionTypeTCP,
		Logger:         logrus.New(),
	}

	client, err := NewTCPClient(config)
	require.NoError(t, err)

	// Test Send without connecting - should fail
	testData := []byte("Hello, TAK Server!")
	err = client.Send(testData)
	assert.Error(t, err)
}

func TestTCPClient_ReceiveWithoutConnect(t *testing.T) {
	// Create a client without connecting
	config := ClientConfig{
		Address:        "127.0.0.1",
		Port:           12345,
		ConnectionType: ConnectionTypeTCP,
		Logger:         logrus.New(),
	}

	client, err := NewTCPClient(config)
	require.NoError(t, err)

	// Test Receive without connecting - should fail
	_, err = client.Receive()
	assert.Error(t, err)
}
