package cot

import (
	"encoding/xml"
	"fmt"
	"sync/atomic"
	"time"
)

// FlowTags represents the _flow-tags_ element in CoT messages
type FlowTags struct {
	XMLName   xml.Name `xml:"_flow-tags_" json:"-"`
	Version   float64  `xml:"version,attr,omitempty" json:"version,omitempty"` // Version of the flow tags
	From      string   `xml:"f,attr" json:"from"`                              // Client that injected the message
	MessageID uint64   `xml:"m,attr" json:"message_id,string"`                 // Message sequence number
	Timestamp int64    `xml:"t,attr" json:"timestamp,string"`                  // Timestamp of creation
	Hops      []string `xml:"h,attr,omitempty" json:"hops,omitempty"`          // Optional hop list
}

// GlobalSequence is an atomic counter for message IDs
var GlobalSequence uint64

// NewFlowTags creates a new flow tag for outgoing messages
func NewFlowTags(clientID string) *FlowTags {
	return &FlowTags{
		From:      clientID,
		MessageID: atomic.AddUint64(&GlobalSequence, 1),
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}
}

// ShouldProcess determines if a message with these flow tags should be processed
// based on previously seen messages - returns true if should process
func ShouldProcess(flowTags *FlowTags, clientID string, seenMessages map[string]uint64) bool {
	if flowTags == nil {
		return true // No flow tags, always process
	}

	// Don't process messages that originated from us
	if flowTags.From == clientID {
		return false
	}

	// Check if we've seen this message before (based on sender + message ID)
	key := fmt.Sprintf("%s:%d", flowTags.From, flowTags.MessageID)
	if lastSeen, exists := seenMessages[key]; exists && lastSeen >= flowTags.MessageID {
		return false // Already processed this message
	}

	// Update the seen messages map
	seenMessages[key] = flowTags.MessageID

	// Prune old entries periodically (implement this based on your needs)

	return true
}

// AddHop adds this client as a hop in the message path
func (f *FlowTags) AddHop(clientID string) {
	if f.Hops == nil {
		f.Hops = make([]string, 0)
	}
	f.Hops = append(f.Hops, clientID)
}
