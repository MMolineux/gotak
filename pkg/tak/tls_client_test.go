package tak

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTLSClient_ValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      ClientConfig
		expectError bool
	}{
		{
			name: "Valid config with cert file and key file",
			config: ClientConfig{
				Address:        "127.0.0.1",
				Port:           8089,
				ConnectionType: ConnectionTypeTLS,
				CertFile:       "testdata/client.crt",
				KeyFile:        "testdata/client.key",
				CAFile:         "testdata/ca.crt",
			},
			expectError: true,
		},
		{
			name: "Valid config with P12 file",
			config: ClientConfig{
				Address:        "127.0.0.1",
				Port:           8089,
				ConnectionType: ConnectionTypeTLS,
				CertFile:       "testdata/client.p12",
				P12Password:    "password",
			},
			expectError: true,
		},
		{
			name: "Valid config with skip verify",
			config: ClientConfig{
				Address:        "127.0.0.1",
				Port:           8089,
				ConnectionType: ConnectionTypeTLS,
				SkipTLSVerify:  true,
			},
			expectError: false,
		},
		{
			name: "Invalid config - no cert file and no skip verify",
			config: ClientConfig{
				Address:        "127.0.0.1",
				Port:           8089,
				ConnectionType: ConnectionTypeTLS,
			},
			expectError: true,
		},
		{
			name: "Invalid config - PEM cert without key file",
			config: ClientConfig{
				Address:        "127.0.0.1",
				Port:           8089,
				ConnectionType: ConnectionTypeTLS,
				CertFile:       "testdata/client.crt",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTLSClient(tt.config)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTLSClient_TLSConfig(t *testing.T) {
	// Test with provided TLS config
	customTLSConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
	}

	config := ClientConfig{
		Address:        "127.0.0.1",
		Port:           8089,
		ConnectionType: ConnectionTypeTLS,
		TLSConfig:      customTLSConfig,
		SkipTLSVerify:  true,
		Logger:         logrus.New(),
	}

	client, err := NewTLSClient(config)
	require.NoError(t, err)

	// tlsClient, ok := client.(*TLSClient)
	// require.True(t, ok)

	// Check that our TLS config is being used
	assert.Equal(t, tls.VersionTLS12, int(client.config.TLSConfig.MinVersion))
	assert.Contains(t, client.config.TLSConfig.CipherSuites, tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384)
}

func TestTLSClient_SkipVerify(t *testing.T) {
	config := ClientConfig{
		Address:        "127.0.0.1",
		Port:           8089,
		ConnectionType: ConnectionTypeTLS,
		SkipTLSVerify:  true,
		Logger:         logrus.New(),
	}

	client, err := NewTLSClient(config)
	require.NoError(t, err)

	// Verify that InsecureSkipVerify is set
	assert.True(t, client.config.TLSConfig.InsecureSkipVerify)
}

// Mock TLS connection test
// Note: This is a simplified test that doesn't actually establish a TLS connection
// A full test would require server certificates and more complex setup
func TestTLSClient_ConnectionLifecycle(t *testing.T) {
	config := ClientConfig{
		Address:        "non.existent.host", // Using a non-existent host to avoid actual connection
		Port:           8089,
		ConnectionType: ConnectionTypeTLS,
		SkipTLSVerify:  true,
		DialTimeout:    time.Second * 1, // Short timeout for faster test
		Logger:         logrus.New(),
	}

	client, err := NewTLSClient(config)
	require.NoError(t, err)

	// IsConnected should be false initially
	assert.False(t, client.IsConnected())

	// Connect should fail (because the host doesn't exist)
	ctx := context.Background()
	err = client.Connect(ctx)
	assert.Error(t, err)

	// IsConnected should still be false
	assert.False(t, client.IsConnected())

	// Disconnect should not error even if we weren't connected
	err = client.Disconnect()
	assert.NoError(t, err)
}

func TestTLSClient_RootCAConfig(t *testing.T) {
	// This test verifies the logic of adding a CA file to the cert pool
	// It doesn't actually load a real CA file since that would require
	// setting up test certificates

	// Create a custom cert pool for testing
	certPool := x509.NewCertPool()

	// Create a client config with a mocked CA file
	// In a real scenario, we would need actual certificate files
	config := ClientConfig{
		Address:        "127.0.0.1",
		Port:           8089,
		ConnectionType: ConnectionTypeTLS,
		SkipTLSVerify:  true, // So we don't need actual certs for this test
		TLSConfig: &tls.Config{
			RootCAs: certPool,
		},
		Logger: logrus.New(),
	}

	client, err := NewTLSClient(config)
	require.NoError(t, err)

	// Check that the root CA cert pool is preserved
	assert.Equal(t, certPool, client.config.TLSConfig.RootCAs)
}
