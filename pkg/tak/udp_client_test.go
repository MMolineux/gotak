package tak

import (
	"context"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUDPClient_Connect(t *testing.T) {
	// Start a mock UDP server
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0})
	require.NoError(t, err)
	defer conn.Close()

	_, portStr, err := net.SplitHostPort(conn.LocalAddr().String())
	require.NoError(t, err)
	port, err := strconv.Atoi(portStr)
	require.NoError(t, err)

	// Create a client with the test server address
	config := ClientConfig{
		Address:        "127.0.0.1",
		Port:           port,
		ConnectionType: ConnectionTypeUDP,
		Logger:         logrus.New(),
	}

	client, err := NewUDPClient(config)
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

func TestUDPClient_SendReceive(t *testing.T) {
	// Start a mock UDP server
	serverConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0})
	require.NoError(t, err)
	defer serverConn.Close()

	_, portStr, err := net.SplitHostPort(serverConn.LocalAddr().String())
	require.NoError(t, err)
	port, err := strconv.Atoi(portStr)
	require.NoError(t, err)

	// Start a goroutine to echo data back
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := make([]byte, 1024)
		n, addr, err := serverConn.ReadFromUDP(buf)
		if err != nil {
			return
		}

		// Echo the data back
		_, err = serverConn.WriteToUDP(buf[:n], addr)
		if err != nil {
			return
		}
	}()

	// Create a client with the test server address
	config := ClientConfig{
		Address:        "127.0.0.1",
		Port:           port,
		ConnectionType: ConnectionTypeUDP,
		Logger:         logrus.New(),
	}

	client, err := NewUDPClient(config)
	require.NoError(t, err)

	// Test Connect
	ctx := context.Background()
	err = client.Connect(ctx)
	require.NoError(t, err)
	defer client.Disconnect()

	// Test Send
	testData := []byte("Hello, TAK Server!")
	err = client.Send(testData)
	require.NoError(t, err)

	// Set a read deadline on the client to avoid test hanging

	client.conn.SetReadDeadline(time.Now().Add(3 * time.Second))

	// Test Receive
	receivedData, err := client.Receive()
	require.NoError(t, err)
	assert.Equal(t, testData, receivedData)

	wg.Wait()
}

func TestUDPClient_ConnectFailure(t *testing.T) {
	// Use a reserved-but-unlikely-to-be-assigned port
	config := ClientConfig{
		Address:        "127.0.0.1",
		Port:           65534, // Unlikely to be open, but not impossible
		ConnectionType: ConnectionTypeUDP,
		Logger:         logrus.New(),
	}

	client, err := NewUDPClient(config)
	require.NoError(t, err)

	// Test Connect - in UDP, this generally won't fail unless there's a binding issue
	ctx := context.Background()
	err = client.Connect(ctx)

	// Whether this succeeds or fails depends on the environment.
	// If it succeeds, make sure we can disconnect.
	if err == nil {
		defer client.Disconnect()
		assert.True(t, client.IsConnected())
	} else {
		assert.False(t, client.IsConnected())
	}
}

func TestUDPClient_SendWithoutConnect(t *testing.T) {
	// Create a client without connecting
	config := ClientConfig{
		Address:        "127.0.0.1",
		Port:           12345,
		ConnectionType: ConnectionTypeUDP,
		Logger:         logrus.New(),
	}

	client, err := NewUDPClient(config)
	require.NoError(t, err)

	// Test Send without connecting - should fail
	testData := []byte("Hello, TAK Server!")
	err = client.Send(testData)
	assert.Error(t, err)
}

func TestUDPClient_ReceiveWithoutConnect(t *testing.T) {
	// Create a client without connecting
	config := ClientConfig{
		Address:        "127.0.0.1",
		Port:           12345,
		ConnectionType: ConnectionTypeUDP,
		Logger:         logrus.New(),
	}

	client, err := NewUDPClient(config)
	require.NoError(t, err)

	// Test Receive without connecting - should fail
	_, err = client.Receive()
	assert.Error(t, err)
}
