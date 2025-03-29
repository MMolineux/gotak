package parser

import (
	"encoding/xml"

	"github.com/angry-kivi/gotak/pkg/cot"
)

// XMLParser handles parsing and serialization of CoT XML messages
type XMLParser struct{}

// ParseCoT converts XML data to a CoT Event
func (p *XMLParser) ParseCoT(data []byte) (*cot.Event, error) {
	var event cot.Event
	err := xml.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// SerializeCoT converts a CoT Event to XML data
func (p *XMLParser) SerializeCoT(event *cot.Event) ([]byte, error) {
	// First marshal the event to XML
	xmlData, err := xml.Marshal(event)
	if err != nil {
		return nil, err
	}

	// Prepend the XML declaration
	xmlDeclaration := []byte("<?xml version='1.0' encoding='UTF-8' standalone='yes'?>\n")
	result := append(xmlDeclaration, xmlData...)

	return result, nil
}

// SerializePrettyCoT converts a CoT Event to indented XML data with XML declaration
func (p *XMLParser) SerializePrettyCoT(event *cot.Event) ([]byte, error) {
	// Marshal the event to indented XML
	xmlData, err := xml.MarshalIndent(event, "", "  ")
	if err != nil {
		return nil, err
	}

	// Prepend the XML declaration
	xmlDeclaration := []byte("<?xml version='1.0' encoding='UTF-8' standalone='yes'?>\n")
	result := append(xmlDeclaration, xmlData...)

	return result, nil
}

// NewXMLParser creates a new XML parser for CoT messages
func NewXMLParser() *XMLParser {
	return &XMLParser{}
}
