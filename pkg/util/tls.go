package util

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/pkcs12"
)

// LoadTLSConfigAuto automatically determines the certificate type (PEM or P12)
// and creates a TLS configuration accordingly
func LoadTLSConfigAuto(certFile, keyFile, password, caFile string, skipVerify bool) (*tls.Config, error) {
	// Check if the certFile has a .p12 or .pfx extension
	if len(certFile) > 4 && (certFile[len(certFile)-4:] == ".p12" || certFile[len(certFile)-4:] == ".pfx") {
		// It's a P12/PFX file
		return LoadTLSConfigFromP12(certFile, password, caFile, skipVerify)
	}

	// Otherwise, assume it's PEM format
	return LoadTLSConfig(certFile, keyFile, caFile, skipVerify)
}

// LoadTLSConfig creates a TLS configuration from certificate and key files
func LoadTLSConfig(certFile, keyFile, caFile string, skipVerify bool) (*tls.Config, error) {
	config := &tls.Config{
		InsecureSkipVerify: skipVerify,
	}
	// Check if key file is provided
	if keyFile == "" {
		return nil, errors.New("PEM certificate format requires a key file")
	}

	// Load client certificate if provided
	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		config.Certificates = []tls.Certificate{cert}
	}

	// Load CA certificate if provided
	if caFile != "" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, errors.New("failed to append CA certificate to pool")
		}

		config.RootCAs = caCertPool
	}

	return config, nil
}

// LoadTLSConfigFromP12 creates a TLS configuration from a PKCS#12 (.p12) file
func LoadTLSConfigFromP12(p12File string, password string, caFile string, skipVerify bool) (*tls.Config, error) {
	config := &tls.Config{
		InsecureSkipVerify: skipVerify,
	}

	// Read the .p12 file
	p12Data, err := os.ReadFile(p12File)
	if err != nil {
		return nil, fmt.Errorf("failed to read P12 file: %w", err)
	}

	// Try to decode all certificates from the P12 file
	blocks, err := pkcs12.ToPEM(p12Data, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decode P12 data: %w", err)
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("no certificates or keys found in P12 file")
	}

	// Variables to hold our certificate components
	var certDERBlock []byte
	var keyDERBlock []byte
	var caCertDERBlocks [][]byte

	// Extract the certificate, private key, and CA certificates
	for _, block := range blocks {
		switch block.Type {
		case "CERTIFICATE":
			// The first certificate is assumed to be the client certificate,
			// any subsequent ones are assumed to be CA certificates
			if certDERBlock == nil {
				certDERBlock = block.Bytes
			} else {
				caCertDERBlocks = append(caCertDERBlocks, block.Bytes)
			}
		case "PRIVATE KEY", "RSA PRIVATE KEY", "EC PRIVATE KEY":
			keyDERBlock = block.Bytes
		}
	}

	// Ensure we have both the certificate and private key
	if certDERBlock == nil || keyDERBlock == nil {
		return nil, fmt.Errorf("P12 file must contain both a certificate and a private key")
	}

	// Create a certificate from the PEM blocks
	cert, err := tls.X509KeyPair(
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDERBlock}),
		pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDERBlock}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create X509 key pair: %w", err)
	}

	// Parse the client certificate for additional information
	clientCert, err := x509.ParseCertificate(certDERBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to parse client certificate: %w", err)
	}
	cert.Leaf = clientCert

	config.Certificates = []tls.Certificate{cert}

	// Create a certificate pool for CA certificates
	certPool := x509.NewCertPool()

	// Add the client certificate to the pool
	certPool.AddCert(clientCert)

	// Add any CA certificates from the P12 file
	for _, caCertDER := range caCertDERBlocks {
		caCert, err := x509.ParseCertificate(caCertDER)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CA certificate: %w", err)
		}
		certPool.AddCert(caCert)
	}

	// Load additional CA certificate if provided
	if caFile != "" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		if !certPool.AppendCertsFromPEM(caCert) {
			return nil, errors.New("failed to append CA certificate to pool")
		}
	}

	// Set the RootCAs to our certPool
	config.RootCAs = certPool

	return config, nil
}
