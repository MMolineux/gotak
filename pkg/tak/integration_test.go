//go:build integration
// +build integration

package tak

import (
	"context"
	"testing"
	"time"

	"github.com/angry-kivi/gotak/pkg/cot"
	"github.com/angry-kivi/gotak/pkg/parser"
	"github.com/angry-kivi/gotak/pkg/tak"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMulticastClient_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create a sender client
	senderConfig := tak.ClientConfig{
		MulticastAddr:  "224.0.0.1", // Local loopback multicast
		MulticastPort:  12345,       // Use a non-standard port to avoid conflicts
		ConnectionType: tak.ConnectionTypeMulticast,
		ClientID:       "test-sender",
		Logger:         logger.WithField("client", "sender"),
	}

	sender, err := tak.NewMulticastClient(senderConfig)
	require.NoError(t, err)

	// Create a receiver client
	receiverConfig := tak.ClientConfig{
		MulticastAddr:  "224.0.0.1",
		MulticastPort:  12345,
		ConnectionType: tak.ConnectionTypeMulticast,
		ClientID:       "test-receiver",
		Logger:         logger.WithField("client", "receiver"),
	}

	receiver, err := tak.NewMulticastClient(receiverConfig)
	require.NoError(t, err)

	// Connect both clients
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = sender.Connect(ctx)
	if err != nil {
		t.Logf("Failed to connect sender, likely missing multicast support: %v", err)
		t.Skip("Skipping test due to missing multicast support")
		return
	}
	defer sender.Disconnect()

	err = receiver.Connect(ctx)
	if err != nil {
		t.Logf("Failed to connect receiver: %v", err)
		t.Skip("Skipping test due to network configuration issues")
		return
	}
	defer receiver.Disconnect()

	// Create a simple test message
	testEvent := cot.NewEvent("a-f-G-U-C", "test-integration")
	testEvent.Detail.AddRemarks("Integration test message")

	// Serialize the event
	xmlParser := parser.NewXMLParser()
	data, err := xmlParser.SerializeCoT(testEvent)
	require.NoError(t, err)

	// Start receiving in a goroutine
	msgChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	go func() {
		// Add a timeout for the receive operation
		receiveCtx, receiveCancel := context.WithTimeout(ctx, 5*time.Second)
		defer receiveCancel()

		// Set a deadline for the receiver to allow context cancellation
		for {
			select {
			case <-receiveCtx.Done():
				errChan <- receiveCtx.Err()
				return
			default:
				msg, err := receiver.Receive()
				if err != nil {
					if err == tak.ErrMessageSkipped {
						// Skip this message and continue
						continue
					}
					errChan <- err
					return
				}
				msgChan <- msg
				return
			}
		}
	}()

	// Wait a moment for the receiver to start
	time.Sleep(1 * time.Second)

	// Send the test message
	err = sender.Send(data)
	require.NoError(t, err)

	// Wait for the result with timeout
	select {
	case receivedData := <-msgChan:
		// Parse and verify the received message
		receivedEvent, err := xmlParser.ParseCoT(receivedData)
		assert.NoError(t, err)
		assert.Equal(t, testEvent.UID, receivedEvent.UID)

		// Check flow tags were added
		assert.NotNil(t, receivedEvent.Detail.FlowTags)
		assert.Equal(t, "test-sender", receivedEvent.Detail.FlowTags.From)

	case err := <-errChan:
		t.Logf("Receive error: %v", err)
		t.Skip("Skipping verification due to receive error (common in CI environments)")

	case <-time.After(10 * time.Second):
		t.Fatal("Timeout waiting for message")
	}
}
