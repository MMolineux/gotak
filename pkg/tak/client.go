package tak

import (
	"context"
	"crypto/tls"
	"errors"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ConnectionType defines the connection type for a TAK client
type ConnectionType string

const (
	// ConnectionTypeTCP indicates a TCP connection
	ConnectionTypeTCP ConnectionType = "tcp"
	// ConnectionTypeTLS indicates a TLS-secured TCP connection
	ConnectionTypeTLS ConnectionType = "tls"
	// ConnectionTypeUDP indicates a UDP connection
	ConnectionTypeUDP ConnectionType = "udp"
	// ConnectionTypeMulticast indicates a multicast UDP connection
	ConnectionTypeMulticast ConnectionType = "multicast"
)

// Client is the interface for TAK clients
type Client interface {
	// Connect establishes a connection to the TAK server
	Connect(ctx context.Context) error

	// Disconnect closes the connection to the TAK server
	Disconnect() error

	// Send transmits data to the TAK server
	Send(data []byte) error

	// Receive waits for and returns data from the TAK server
	Receive() ([]byte, error)

	// IsConnected returns true if the client is connected
	IsConnected() bool
}

// ClientConfig holds configuration for a TAK client
type ClientConfig struct {
	Address        string
	Port           int
	ClientID       string
	ConnectionType ConnectionType

	// TLS configuration
	TLSConfig     *tls.Config
	CertFile      string
	KeyFile       string
	CAFile        string
	P12Password   string
	SkipTLSVerify bool

	// Timeout configurations
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// Multicast-specific configuration
	MulticastAddr string
	MulticastPort int

	// Logging
	Logger logrus.FieldLogger
}

// NewClient creates a new TAK client with the given configuration
func NewClient(config ClientConfig) (Client, error) {
	switch config.ConnectionType {
	case ConnectionTypeTCP:
		return NewTCPClient(config)
	case ConnectionTypeTLS:
		// Check if any certificate files are provided
		if config.CertFile == "" {
			if !config.SkipTLSVerify {
				return nil, errors.New("TLS connection requires either a certificate file or skip-verify option")
			}

			// Log warning if logger is provided
			if config.Logger != nil {
				config.Logger.Warn("TLS enabled with skip-verify but no certificate file provided")
				config.Logger.Warn("This connection will be encrypted but not authenticated")
			}
		} else {
			// Check certificate format
			isPKCS12 := false
			if len(config.CertFile) > 4 {
				ext := strings.ToLower(config.CertFile[len(config.CertFile)-4:])
				isPKCS12 = ext == ".p12" || ext == ".pfx"
			}

			if isPKCS12 {
				if config.Logger != nil {
					config.Logger.Info("Using PKCS#12 (.p12) certificate format")
					if config.P12Password == "" {
						config.Logger.Warn("No password provided for .p12 file. If your certificate is password-protected, the connection may fail")
					}
				}
			} else {
				// For PEM format certificates, ensure key file is provided
				if config.KeyFile == "" {
					return nil, errors.New("PEM certificate format requires a key file")
				}

				if config.Logger != nil {
					config.Logger.Info("Using PEM certificate format")
				}
			}
		}
		return NewTLSClient(config)
	case ConnectionTypeUDP:
		return NewUDPClient(config)
	case ConnectionTypeMulticast:
		return NewMulticastClient(config)
	default:
		return nil, errors.New("unsupported connection type: " + string(config.ConnectionType))
	}
}
