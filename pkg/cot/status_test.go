package cot

import (
	"encoding/xml"
	"testing"
)

func TestStatusMarshalXML(t *testing.T) {
	// Given
	status := Status{
		Battery: 75,
	}

	// When
	data, err := xml.Marshal(status)
	if err != nil {
		t.Fatalf("Failed to marshal Status to XML: %v", err)
	}

	// Then
	expected := `<status battery="75"></status>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestStatusUnmarshalXML(t *testing.T) {
	// Given
	xmlData := `<status battery="75"></status>`

	// When
	var status Status
	err := xml.Unmarshal([]byte(xmlData), &status)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to Status: %v", err)
	}

	// Then
	if status.Battery != 75 {
		t.Errorf("Battery not correctly unmarshaled. Got: %d, Expected: %d", status.Battery, 75)
	}
}

func TestStatusSetBattery(t *testing.T) {
	// Given
	status := Status{}
	batteryPercent := 85

	// When
	result := status.SetBattery(batteryPercent)

	// Then
	if status.Battery != batteryPercent {
		t.Errorf("Battery not set correctly. Got: %d, Expected: %d", status.Battery, batteryPercent)
	}
	if result != &status {
		t.Errorf("SetBattery did not return the status instance for chaining")
	}
}

func TestStatusSetReadiness(t *testing.T) {
	// Given
	status := Status{}
	readiness := true

	// When
	result := status.SetReadiness(readiness)

	// Then
	if status.Readiness != readiness {
		t.Errorf("Readiness not set correctly. Got: %t, Expected: %t", status.Readiness, readiness)
	}
	if result != &status {
		t.Errorf("SetReadiness did not return the status instance for chaining")
	}
}

func TestStatusComplete(t *testing.T) {
	// Given
	status := Status{}

	// When
	result := status.SetBattery(75).SetReadiness(true)

	// Then
	if status.Battery != 75 {
		t.Errorf("Battery not set correctly. Got: %d, Expected: %d", status.Battery, 75)
	}
	if status.Readiness != true {
		t.Errorf("Readiness not set correctly. Got: %t, Expected: %t", status.Readiness, true)
	}
	if result != &status {
		t.Errorf("Method chaining did not return the status instance")
	}
}

func TestStatusMarshalXMLWithReadiness(t *testing.T) {
	// Given
	status := Status{
		Battery:   75,
		Readiness: true,
	}

	// When
	data, err := xml.Marshal(status)
	if err != nil {
		t.Fatalf("Failed to marshal Status to XML: %v", err)
	}

	// Then
	expected := `<status battery="75" readiness="true"></status>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestStatusUnmarshalXMLWithReadiness(t *testing.T) {
	// Given
	xmlData := `<status battery="75" readiness="true"></status>`

	// When
	var status Status
	err := xml.Unmarshal([]byte(xmlData), &status)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to Status: %v", err)
	}

	// Then
	if status.Battery != 75 {
		t.Errorf("Battery not correctly unmarshaled. Got: %d, Expected: %d", status.Battery, 75)
	}
	if status.Readiness != true {
		t.Errorf("Readiness not correctly unmarshaled. Got: %t, Expected: %t", status.Readiness, true)
	}
}
