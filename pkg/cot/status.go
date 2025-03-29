package cot

import "encoding/xml"

// Status represents the status element in CoT messages
type Status struct {
	XMLName   xml.Name `xml:"status" json:"-"`
	Battery   int      `xml:"battery,attr,omitempty" json:"battery,omitempty"`
	Readiness bool     `xml:"readiness,attr,omitempty" json:"readiness,omitempty"`
}

// SetBattery sets the battery percentage value
func (s *Status) SetBattery(batteryPercent int) *Status {
	s.Battery = batteryPercent
	return s
}

// SetReadiness sets the readiness status
func (s *Status) SetReadiness(readiness bool) *Status {
	s.Readiness = readiness
	return s
}
