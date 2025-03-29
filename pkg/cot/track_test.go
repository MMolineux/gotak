package cot

import (
	"encoding/xml"
	"testing"
	"time"
)

func TestTrackMarshalXML(t *testing.T) {
	// Given
	timestamp := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	track := Track{
		Course:    180.5,
		Speed:     25.3,
		Slope:     10.2,
		Etype:     "t-x-d",
		TimeStamp: timestamp,
	}

	// When
	data, err := xml.Marshal(track)
	if err != nil {
		t.Fatalf("Failed to marshal Track to XML: %v", err)
	}

	// Then
	expected := `<track course="180.5" speed="25.3" slope="10.2" etype="t-x-d" timeStamp="2023-05-15T10:30:00Z"></track>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestTrackUnmarshalXML(t *testing.T) {
	// Given
	xmlData := `<track course="180.5" speed="25.3" slope="10.2" etype="t-x-d" timeStamp="2023-05-15T10:30:00Z"></track>`

	// When
	var track Track
	err := xml.Unmarshal([]byte(xmlData), &track)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to Track: %v", err)
	}

	// Then
	expectedTimestamp := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	if track.Course != 180.5 {
		t.Errorf("Course not correctly unmarshaled. Got: %f, Expected: %f", track.Course, 180.5)
	}
	if track.Speed != 25.3 {
		t.Errorf("Speed not correctly unmarshaled. Got: %f, Expected: %f", track.Speed, 25.3)
	}
	if track.Slope != 10.2 {
		t.Errorf("Slope not correctly unmarshaled. Got: %f, Expected: %f", track.Slope, 10.2)
	}
	if track.Etype != "t-x-d" {
		t.Errorf("Etype not correctly unmarshaled. Got: %s, Expected: %s", track.Etype, "t-x-d")
	}
	if !track.TimeStamp.Equal(expectedTimestamp) {
		t.Errorf("TimeStamp not correctly unmarshaled. Got: %v, Expected: %v", track.TimeStamp, expectedTimestamp)
	}
}

func TestTrackSetters(t *testing.T) {
	// Create a Track with setters (if implemented)
	track := Track{}

	// Test setting values
	track.Course = 180.5
	track.Speed = 25.3
	track.Slope = 10.2
	track.Etype = "t-x-d"
	timestamp := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	track.TimeStamp = timestamp

	// Verify values
	if track.Course != 180.5 {
		t.Errorf("Course not set correctly. Got: %f, Expected: %f", track.Course, 180.5)
	}
	if track.Speed != 25.3 {
		t.Errorf("Speed not set correctly. Got: %f, Expected: %f", track.Speed, 25.3)
	}
	if track.Slope != 10.2 {
		t.Errorf("Slope not set correctly. Got: %f, Expected: %f", track.Slope, 10.2)
	}
	if track.Etype != "t-x-d" {
		t.Errorf("Etype not set correctly. Got: %s, Expected: %s", track.Etype, "t-x-d")
	}
	if !track.TimeStamp.Equal(timestamp) {
		t.Errorf("TimeStamp not set correctly. Got: %v, Expected: %v", track.TimeStamp, timestamp)
	}
}
