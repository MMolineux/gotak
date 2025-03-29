package tak

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/angry-kivi/gotak/pkg/parser"
	"github.com/sirupsen/logrus"
)

// Define a new error type for skipped messages
var ErrMessageSkipped = errors.New("message skipped due to flow tag rules")

// DefaultLogger returns a logger configured with reasonable defaults
func DefaultLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	return logger
}

type udpConnInterface interface {
	Close() error
	ReadFromUDP(b []byte) (int, *net.UDPAddr, error)
	WriteTo(b []byte, addr net.Addr) (int, error)
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

// MulticastClient implements the Client interface for multicast communication
type MulticastClient struct {
	config        ClientConfig
	conn          udpConnInterface
	multicastAddr *net.UDPAddr
	seenMessages  map[string]uint64  // Map of "from:msgID" -> messageID
	mutex         sync.RWMutex       // Protects seenMessages
	ctx           context.Context    // For cancellation
	cancel        context.CancelFunc // For cleanup
}

// NewMulticastClient creates a new client for TAK multicast communication
func NewMulticastClient(config ClientConfig) (Client, error) {
	// Ensure logger exists
	if config.Logger == nil {
		logger := logrus.New()
		logger.SetLevel(logrus.InfoLevel)
		config.Logger = logger
	}

	// If multicast address is not set, use the default
	if config.MulticastAddr == "" {
		config.MulticastAddr = DefaultMulticastAddr
	}

	// If multicast port is not set, use the default
	if config.MulticastPort == 0 {
		config.MulticastPort = DefaultMulticastPort
	}

	client := &MulticastClient{
		config:       config,
		seenMessages: make(map[string]uint64),
	}

	client.config.Logger.WithFields(logrus.Fields{
		"address": config.MulticastAddr,
		"port":    config.MulticastPort,
	}).Debug("Multicast client created")

	return client, nil
}

// Connect establishes a multicast connection
func (c *MulticastClient) Connect(ctx context.Context) error {
	if c.conn != nil {
		return errors.New("client is already connected")
	}

	c.config.Logger.WithFields(logrus.Fields{
		"address": c.config.MulticastAddr,
		"port":    c.config.MulticastPort,
	}).Debug("Connecting to multicast group")

	// Create context with cancellation for client operations
	c.ctx, c.cancel = context.WithCancel(ctx)
	// Start background pruning of seen messages
	go c.pruneSeenMessages()

	// Resolve the multicast address
	var err error
	c.multicastAddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", c.config.MulticastAddr, c.config.MulticastPort))
	if err != nil {
		c.config.Logger.WithError(err).Error("Failed to resolve multicast address")
		return err
	}

	// Find an appropriate interface to use for multicast
	interfaces, err := net.Interfaces()
	if err != nil {
		c.config.Logger.WithError(err).Error("Failed to get network interfaces")
		return err
	}

	var iface *net.Interface
	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && i.Flags&net.FlagMulticast != 0 && i.Flags&net.FlagLoopback == 0 {
			iface = &i
			break
		}
	}

	if iface == nil {
		return errors.New("no suitable multicast interface found")
	}

	c.config.Logger.WithField("interface", iface.Name).Debug("Using interface for multicast")

	// Create the UDP connection
	c.conn, err = net.ListenMulticastUDP("udp", iface, c.multicastAddr)
	if err != nil {
		c.config.Logger.WithError(err).Error("Failed to create multicast connection")
		return err
	}

	// Set the read and write deadlines if configured
	if c.config.ReadTimeout > 0 {
		err = c.conn.SetReadDeadline(time.Now().Add(c.config.ReadTimeout))
		if err != nil {
			c.config.Logger.WithError(err).Error("Failed to set read deadline")
			return err
		}
	}

	if c.config.WriteTimeout > 0 {
		err = c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout))
		if err != nil {
			c.config.Logger.WithError(err).Error("Failed to set write deadline")
			return err
		}
	}

	c.config.Logger.Info("Connected to multicast group successfully")
	return nil
}

// Disconnect closes the multicast connection
func (c *MulticastClient) Disconnect() error {
	if c.conn == nil {
		return errors.New("client is not connected")
	}

	// Cancel any background operations
	if c.cancel != nil {
		c.cancel()
	}

	c.config.Logger.Debug("Disconnecting multicast client")
	err := c.conn.Close()
	if err == nil {
		c.conn = nil
		c.multicastAddr = nil
	} else {
		c.config.Logger.WithError(err).Error("Error closing multicast connection")
	}
	return err
}

// Send transmits data to the multicast group
func (c *MulticastClient) Send(data []byte) error {
	if c.conn == nil || c.multicastAddr == nil {
		return errors.New("client is not connected")
	}
	// Try to parse and enrich with flow tags if it's XML
	enrichedData, err := c.enrichWithFlowTags(data)
	if err != nil {
		// Not CoT XML or other error, just send as-is
		c.config.Logger.WithField("reason", err.Error()).
			Debug("Sending raw data without flow tag processing")
		return c.sendRaw(data)
	}

	// Send the enriched data with flow tags
	return c.sendRaw(enrichedData)
}

// Receive reads data from the multicast group
func (c *MulticastClient) Receive() ([]byte, error) {
	if c.conn == nil {
		return nil, errors.New("client is not connected")
	}

	c.config.Logger.Debug("Waiting to receive data via multicast")

	if c.config.ReadTimeout > 0 {
		deadline := time.Now().Add(c.config.ReadTimeout)
		if err := c.conn.SetReadDeadline(deadline); err != nil {
			c.config.Logger.WithError(err).Error("Failed to set read deadline")
			return nil, err
		}
	}

	// Buffer for incoming data
	buffer := make([]byte, 8192) // 8K buffer should be enough for most TAK messages
	n, _, err := c.conn.ReadFromUDP(buffer)
	if err != nil {
		c.config.Logger.WithError(err).Error("Failed to receive multicast data")
		return nil, err
	}

	data := buffer[:n]
	c.config.Logger.WithField("bytes", n).Debug("Data received via multicast")

	// Process flow tags if this is CoT XML
	shouldProcess, err := c.handleMessage(data)
	if err != nil {
		// Not CoT XML or parsing error, just return the data
		return data, nil
	}

	if !shouldProcess {
		c.config.Logger.Debug("Skipping message due to flow tags (already seen or from self)")
		return nil, ErrMessageSkipped
	}

	return data, nil
}

// IsConnected returns true if the client has an active connection
func (c *MulticastClient) IsConnected() bool {
	return c.conn != nil && c.multicastAddr != nil
}

// pruneSeenMessages periodically cleans up the message history to prevent memory leaks
func (c *MulticastClient) pruneSeenMessages() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.mutex.Lock()

			// If we have too many entries, keep only the most recent ones
			if len(c.seenMessages) > 1000 {
				c.config.Logger.WithField("count", len(c.seenMessages)).
					Debug("Pruning seen messages")

				// Simple approach: just reset if we get too big
				// For a more sophisticated approach, you could keep only recent entries
				c.seenMessages = make(map[string]uint64)
			}

			c.mutex.Unlock()
		}
	}
}

// handleMessage checks if a message should be processed based on flow tags
// Returns: (shouldProcess bool, err error)
// - shouldProcess: true if the message should be processed, false if it should be skipped
// - err: any error that occurred during parsing (if nil, flow tag processing succeeded)
func (c *MulticastClient) handleMessage(data []byte) (bool, error) {
	// Parse the message to extract flow tags
	xmlParser := parser.NewXMLParser()
	event, err := xmlParser.ParseCoT(data)
	if err != nil {
		// Not a valid CoT XML message, just process it without flow tag handling
		return true, err
	}

	// If no flow tags, always process
	if event.Detail.FlowTags == nil {
		return true, nil
	}

	// Don't process messages that originated from us
	if event.Detail.FlowTags.From == c.config.ClientID {
		return false, nil
	}

	// Check if we've seen this message before (based on sender + message ID)
	key := fmt.Sprintf("%s:%d", event.Detail.FlowTags.From, event.Detail.FlowTags.MessageID)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if lastSeen, exists := c.seenMessages[key]; exists && lastSeen >= event.Detail.FlowTags.MessageID {
		// We've already processed this message
		return false, nil
	}

	// Update the seen messages map
	c.seenMessages[key] = event.Detail.FlowTags.MessageID

	// Process the message
	return true, nil
}

// sendRaw sends raw data without any processing
func (c *MulticastClient) sendRaw(data []byte) error {
	c.config.Logger.WithField("bytes", len(data)).Debug("Sending data via multicast")

	if c.config.WriteTimeout > 0 {
		deadline := time.Now().Add(c.config.WriteTimeout)
		if err := c.conn.SetWriteDeadline(deadline); err != nil {
			c.config.Logger.WithError(err).Error("Failed to set write deadline")
			return err
		}
	}

	_, err := c.conn.WriteTo(data, c.multicastAddr)
	if err != nil {
		c.config.Logger.WithError(err).Error("Failed to send multicast data")
	}
	return err
}

// enrichWithFlowTags tries to parse CoT XML and add flow tags
func (c *MulticastClient) enrichWithFlowTags(data []byte) ([]byte, error) {
	// Try to parse as CoT XML
	xmlParser := parser.NewXMLParser()

	event, err := xmlParser.ParseCoT(data)
	if err != nil {
		return nil, fmt.Errorf("not valid CoT XML: %w", err)
	}

	// Add flow tags if they don't exist
	if event.Detail.FlowTags == nil {
		// Add new flow tags
		c.config.Logger.Debug("Adding flow tags to outgoing message")
		event.Detail.AddFlowTags(c.config.ClientID)
	} else if event.Detail.FlowTags.From != c.config.ClientID {
		// This is a message we're forwarding - add our client ID to the hop list
		c.config.Logger.Debug("Updating flow tags for forwarded message")
		event.Detail.FlowTags.AddHop(c.config.ClientID)
	}

	// Serialize back to XML

	enrichedData, err := xmlParser.SerializeCoT(event)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize enriched XML: %w", err)
	}

	return enrichedData, nil
}
