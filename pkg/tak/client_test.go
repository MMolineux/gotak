package tak

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name         string
		config       ClientConfig
		expectedType string
		expectError  bool
	}{
		{
			name: "TCP Client",
			config: ClientConfig{
				Address:        "127.0.0.1",
				Port:           8089,
				ConnectionType: ConnectionTypeTCP,
				Logger:         logrus.New(),
			},
			expectedType: "*tak.TCPClient",
			expectError:  false,
		},
		{
			name: "TLS Client with SkipVerify",
			config: ClientConfig{
				Address:        "127.0.0.1",
				Port:           8089,
				ConnectionType: ConnectionTypeTLS,
				SkipTLSVerify:  true,
				Logger:         logrus.New(),
			},
			expectedType: "*tak.TLSClient",
			expectError:  false,
		},
		{
			name: "TLS Client without cert or skipverify",
			config: ClientConfig{
				Address:        "127.0.0.1",
				Port:           8089,
				ConnectionType: ConnectionTypeTLS,
				Logger:         logrus.New(),
			},
			expectedType: "",
			expectError:  true,
		},
		{
			name: "UDP Client",
			config: ClientConfig{
				Address:        "127.0.0.1",
				Port:           8089,
				ConnectionType: ConnectionTypeUDP,
				Logger:         logrus.New(),
			},
			expectedType: "*tak.UDPClient",
			expectError:  false,
		},
		{
			name: "Multicast Client",
			config: ClientConfig{
				MulticastAddr:  "239.2.3.1",
				MulticastPort:  6969,
				ConnectionType: ConnectionTypeMulticast,
				Logger:         logrus.New(),
			},
			expectedType: "*tak.MulticastClient",
			expectError:  false,
		},
		{
			name: "Invalid Connection Type",
			config: ClientConfig{
				Address:        "127.0.0.1",
				Port:           8089,
				ConnectionType: "invalid",
				Logger:         logrus.New(),
			},
			expectedType: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, tt.expectedType, GetType(client))
			}
		})
	}
}

// Helper function to get the type of a client as a string
func GetType(client Client) string {
	if client == nil {
		return ""
	}
	return GetTypeName(client)
}

// GetTypeName returns the type of the interface as a string
func GetTypeName(v interface{}) string {
	if v == nil {
		return ""
	}
	return TypeName(v)
}

// TypeName gets the name of the type as a string
func TypeName(v interface{}) string {
	if _, ok := v.(*TCPClient); ok {
		return "*tak.TCPClient"
	} else if _, ok := v.(*TLSClient); ok {
		return "*tak.TLSClient"
	} else if _, ok := v.(*UDPClient); ok {
		return "*tak.UDPClient"
	} else if _, ok := v.(*MulticastClient); ok {
		return "*tak.MulticastClient"
	}
	return ""
}
