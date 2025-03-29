package tak

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockUDPConn implements udpConnInterface for testing
type mockUDPConn struct {
	closed         bool
	readData       []byte
	readError      error
	writeError     error
	deadlineError  error
	writtenData    [][]byte
	writeAddresses []net.Addr
}

func (m *mockUDPConn) Close() error {
	m.closed = true
	return nil
}

func (m *mockUDPConn) ReadFromUDP(b []byte) (int, *net.UDPAddr, error) {
	if m.readError != nil {
		return 0, nil, m.readError
	}
	n := copy(b, m.readData)
	return n, &net.UDPAddr{IP: net.ParseIP("239.2.3.1"), Port: 6969}, nil
}

func (m *mockUDPConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	if m.writeError != nil {
		return 0, m.writeError
	}
	m.writtenData = append(m.writtenData, append([]byte{}, b...))
	m.writeAddresses = append(m.writeAddresses, addr)
	return len(b), nil
}

func (m *mockUDPConn) SetReadDeadline(t time.Time) error {
	return m.deadlineError
}

func (m *mockUDPConn) SetWriteDeadline(t time.Time) error {
	return m.deadlineError
}

/*func TestMulticastClient_Config(t *testing.T) {
	tests := []struct {
		name        string
		config      ClientConfig
		expectError bool
	}{
		{
			name: "Valid config",
			config: ClientConfig{
				MulticastAddr:  "239.2.3.1", // Valid multicast address
				MulticastPort:  6969,
				ConnectionType: ConnectionTypeMulticast,
			},
			expectError: false,
		},
		{
			name: "Invalid multicast address",
			config: ClientConfig{
				MulticastAddr:  "127.0.0.1", // Not a multicast address
				MulticastPort:  6969,
				ConnectionType: ConnectionTypeMulticast,
			},
			expectError: true,
		},
		{
			name: "Empty multicast address",
			config: ClientConfig{
				MulticastAddr:  "",
				MulticastPort:  6969,
				ConnectionType: ConnectionTypeMulticast,
			},
			expectError: true,
		},
		{
			name: "Invalid port number",
			config: ClientConfig{
				MulticastAddr:  "239.2.3.1",
				MulticastPort:  -1, // Invalid port
				ConnectionType: ConnectionTypeMulticast,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMulticastClient(tt.config)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}*/

// Test the multicast connection lifecycle
// Note: This doesn't actually join a multicast group since that requires
// network configuration, but we test the interface methods
func TestMulticastClient_ConnectionLifecycle(t *testing.T) {
	// Use a valid multicast address
	config := ClientConfig{
		MulticastAddr:  "239.2.3.1",
		MulticastPort:  6969,
		ConnectionType: ConnectionTypeMulticast,
		Logger:         logrus.New(),
	}

	client, err := NewMulticastClient(config)
	require.NoError(t, err)

	// IsConnected should be false initially
	assert.False(t, client.IsConnected())

	// Connect might succeed or fail depending on network configuration
	// We just want to ensure it doesn't panic
	ctx := context.Background()
	err = client.Connect(ctx)

	// Whether this succeeds or fails depends on the network environment
	if err == nil {
		// Successfully connected
		assert.True(t, client.IsConnected())

		// Test disconnect
		err = client.Disconnect()
		assert.NoError(t, err)
		assert.False(t, client.IsConnected())
	} else {
		// Failed to connect, which is expected in most CI environments
		t.Logf("Connect failed as expected in this environment: %v", err)
		assert.False(t, client.IsConnected())
	}
}

func TestMulticastClient_SendWithoutConnect(t *testing.T) {
	// Create a client without connecting
	config := ClientConfig{
		MulticastAddr:  "239.2.3.1",
		MulticastPort:  6969,
		ConnectionType: ConnectionTypeMulticast,
		Logger:         logrus.New(),
	}

	client, err := NewMulticastClient(config)
	require.NoError(t, err)

	// Test Send without connecting - should fail
	testData := []byte("Hello, TAK Multicast!")
	err = client.Send(testData)
	assert.Error(t, err)
}

func TestMulticastClient_ReceiveWithoutConnect(t *testing.T) {
	// Create a client without connecting
	config := ClientConfig{
		MulticastAddr:  "239.2.3.1",
		MulticastPort:  6969,
		ConnectionType: ConnectionTypeMulticast,
		Logger:         logrus.New(),
	}

	client, err := NewMulticastClient(config)
	require.NoError(t, err)

	// Test Receive without connecting - should fail
	_, err = client.Receive()
	assert.Error(t, err)
}
func TestMulticastClient_SendReceive(t *testing.T) {
	// Create a client with mocked connections
	config := ClientConfig{
		MulticastAddr:  "239.2.3.1",
		MulticastPort:  6969,
		ConnectionType: ConnectionTypeMulticast,
		ClientID:       "test-client",
		Logger:         logrus.New(),
	}

	client := &MulticastClient{
		config:       config,
		seenMessages: make(map[string]uint64),
	}

	// Setup mocks
	mockConn := &mockUDPConn{
		readData: []byte("<event><detail></detail></event>"),
	}
	client.conn = mockConn
	client.multicastAddr = &net.UDPAddr{IP: net.ParseIP("239.2.3.1"), Port: 6969}
	client.ctx, client.cancel = context.WithCancel(context.Background())

	// Test Send
	err := client.Send([]byte("<event><detail></detail></event>"))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockConn.writtenData))

	// Test Receive
	data, err := client.Receive()
	assert.NoError(t, err)
	assert.Equal(t, "<event><detail></detail></event>", string(data))

	// Cleanup
	client.Disconnect()
	assert.True(t, mockConn.closed)
}

// Add tests for flow tag handling
func TestMulticastClient_FlowTagHandling(t *testing.T) {
	// Create a client with mocked connections
	config := ClientConfig{
		MulticastAddr:  "239.2.3.1",
		MulticastPort:  6969,
		ConnectionType: ConnectionTypeMulticast,
		ClientID:       "test-client",
		Logger:         logrus.New(),
	}

	client := &MulticastClient{
		config:       config,
		seenMessages: make(map[string]uint64),
	}

	// Setup mocks
	mockConn := &mockUDPConn{}
	client.conn = mockConn
	client.multicastAddr = &net.UDPAddr{IP: net.ParseIP("239.2.3.1"), Port: 6969}
	client.ctx, client.cancel = context.WithCancel(context.Background())

	// Test with real XML messages instead of mock parsers

	// Test 1: Message with no flow tags
	noTagsData := []byte("<event><detail></detail></event>")

	// Send should add flow tags
	err := client.Send(noTagsData)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockConn.writtenData))
	// Check that output contains flow tags
	assert.Contains(t, string(mockConn.writtenData[0]), "_flow-tags_")
	assert.Contains(t, string(mockConn.writtenData[0]), "f=\"test-client\"")

	// Test 2: Self-originated message
	selfMsg := []byte("<event><detail><_flow-tags_ f=\"test-client\" m=\"123\" t=\"1000\"/></detail></event>")
	mockConn.readData = selfMsg

	// Should be skipped because it's from us
	_, err = client.Receive()
	assert.Equal(t, ErrMessageSkipped, err)

	// Test 3: Message from another client (first time)
	otherMsg := []byte("<event><detail><_flow-tags_ f=\"other-client\" m=\"456\" t=\"1000\"/></detail></event>")
	mockConn.readData = otherMsg

	// First time should be processed
	data, err := client.Receive()
	assert.NoError(t, err)
	assert.Equal(t, string(otherMsg), string(data))

	// Test 4: Same message again (duplicate)
	// Second time should be skipped
	_, err = client.Receive()
	assert.Equal(t, ErrMessageSkipped, err)

	// Test 5: Different message from same client
	newMsg := []byte("<event><detail><_flow-tags_ f=\"other-client\" m=\"789\" t=\"1000\"/></detail></event>")
	mockConn.readData = newMsg

	// Should be processed (different message ID)
	data, err = client.Receive()
	assert.NoError(t, err)
	assert.Equal(t, string(newMsg), string(data))

	// Test 6: Forwarding a message
	forwardMsg := []byte("<event><detail><_flow-tags_ f=\"third-client\" m=\"555\" t=\"1000\"/></detail></event>")

	// Reset written data
	mockConn.writtenData = nil

	// Forward should add our client to hops
	err = client.Send(forwardMsg)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockConn.writtenData))

	// Check that output contains our client ID in hops
	sentData := string(mockConn.writtenData[0])
	assert.Contains(t, sentData, "_flow-tags_")
	assert.Contains(t, sentData, "f=\"third-client\"")
	assert.Contains(t, sentData, "h=\"test-client\"")

	// Cleanup
	client.Disconnect()
}

// Test for handling invalid XML
func TestMulticastClient_InvalidXml(t *testing.T) {
	// Create a client with mocked connections
	config := ClientConfig{
		MulticastAddr:  "239.2.3.1",
		MulticastPort:  6969,
		ConnectionType: ConnectionTypeMulticast,
		ClientID:       "test-client",
		Logger:         logrus.New(),
	}

	client := &MulticastClient{
		config:       config,
		seenMessages: make(map[string]uint64),
	}

	// Setup mocks
	mockConn := &mockUDPConn{}
	client.conn = mockConn
	client.multicastAddr = &net.UDPAddr{IP: net.ParseIP("239.2.3.1"), Port: 6969}
	client.ctx, client.cancel = context.WithCancel(context.Background())

	// Invalid XML
	invalidXml := []byte("This is not XML")

	// Send should not add flow tags, just send as-is
	err := client.Send(invalidXml)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockConn.writtenData))
	assert.Equal(t, string(invalidXml), string(mockConn.writtenData[0]))

	// Test receiving invalid XML
	mockConn.readData = invalidXml

	// Should just pass through without flow tag processing
	data, err := client.Receive()
	assert.NoError(t, err)
	assert.Equal(t, string(invalidXml), string(data))

	// Cleanup
	client.Disconnect()
}

// Test for network errors
func TestMulticastClient_NetworkErrors(t *testing.T) {
	// Create a client with mocked connections
	config := ClientConfig{
		MulticastAddr:  "239.2.3.1",
		MulticastPort:  6969,
		ConnectionType: ConnectionTypeMulticast,
		ClientID:       "test-client",
		Logger:         logrus.New(),
	}

	client := &MulticastClient{
		config:       config,
		seenMessages: make(map[string]uint64),
	}

	// Setup mocks with errors
	mockConn := &mockUDPConn{
		readError:     errors.New("read error"),
		writeError:    errors.New("write error"),
		deadlineError: errors.New("deadline error"),
	}
	client.conn = mockConn
	client.multicastAddr = &net.UDPAddr{IP: net.ParseIP("239.2.3.1"), Port: 6969}
	client.ctx, client.cancel = context.WithCancel(context.Background())

	// Test read error
	_, err := client.Receive()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read error")

	// Test write error
	err = client.Send([]byte("test"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "write error")

	// Test with timeout setting
	client.config.ReadTimeout = 1 * time.Second
	client.config.WriteTimeout = 1 * time.Second

	// Should get deadline errors
	_, err = client.Receive()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "deadline error")

	err = client.Send([]byte("test"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "deadline error")

	// Cleanup
	client.Disconnect()
}

// Test handleMessage method directly
func TestMulticastClient_HandleMessage(t *testing.T) {
	// Create a client for testing
	config := ClientConfig{
		ClientID: "test-client",
		Logger:   logrus.New(),
	}

	client := &MulticastClient{
		config:       config,
		seenMessages: make(map[string]uint64),
	}

	// Test cases
	testCases := []struct {
		name          string
		message       []byte
		shouldProcess bool
		expectError   bool
	}{
		{
			name:          "Not XML",
			message:       []byte("This is not XML"),
			shouldProcess: true,
			expectError:   true,
		},
		{
			name:          "XML without flow tags",
			message:       []byte("<event><detail></detail></event>"),
			shouldProcess: true,
			expectError:   false,
		},
		{
			name:          "XML with flow tags from self",
			message:       []byte("<event><detail><_flow-tags_ f=\"test-client\" m=\"123\" t=\"1000\"/></detail></event>"),
			shouldProcess: false,
			expectError:   false,
		},
		{
			name:          "XML with flow tags from other",
			message:       []byte("<event><detail><_flow-tags_ f=\"other-client\" m=\"456\" t=\"1000\"/></detail></event>"),
			shouldProcess: true,
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shouldProcess, err := client.handleMessage(tc.message)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.shouldProcess, shouldProcess)
		})
	}
}
