package cot

import (
	"encoding/xml"
	"testing"
	"time"
)

func TestRemarksMarshalXML(t *testing.T) {
	// Given
	timestamp := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	remarks := Remarks{
		Source:   "Tactical Device",
		SourceID: "ANDROID-355161070128345",
		Time:     &timestamp,
		To:       "All",
		Text:     "Situation report for Team Alpha",
	}

	// When
	data, err := xml.Marshal(remarks)
	if err != nil {
		t.Fatalf("Failed to marshal Remarks to XML: %v", err)
	}

	// Then
	expected := `<remarks source="Tactical Device" sourceID="ANDROID-355161070128345" time="2023-05-15T10:30:00Z" to="All">Situation report for Team Alpha</remarks>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestRemarksMarshalXMLWithMinimalFields(t *testing.T) {
	// Given
	remarks := Remarks{
		Text: "Situation report for Team Alpha",
	}

	// When
	data, err := xml.Marshal(remarks)
	if err != nil {
		t.Fatalf("Failed to marshal Remarks to XML: %v", err)
	}

	// Then
	expected := `<remarks>Situation report for Team Alpha</remarks>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestRemarksUnmarshalXML(t *testing.T) {
	// Given
	xmlData := `<remarks source="Tactical Device" sourceID="ANDROID-355161070128345" time="2023-05-15T10:30:00Z" to="All">Situation report for Team Alpha</remarks>`

	// When
	var remarks Remarks
	err := xml.Unmarshal([]byte(xmlData), &remarks)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to Remarks: %v", err)
	}

	// Then
	expectedTimestamp := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	if remarks.Source != "Tactical Device" {
		t.Errorf("Source not correctly unmarshaled. Got: %s, Expected: %s", remarks.Source, "Tactical Device")
	}
	if remarks.SourceID != "ANDROID-355161070128345" {
		t.Errorf("SourceID not correctly unmarshaled. Got: %s, Expected: %s", remarks.SourceID, "ANDROID-355161070128345")
	}
	if !remarks.Time.Equal(expectedTimestamp) {
		t.Errorf("Time not correctly unmarshaled. Got: %v, Expected: %v", remarks.Time, expectedTimestamp)
	}
	if remarks.To != "All" {
		t.Errorf("To not correctly unmarshaled. Got: %s, Expected: %s", remarks.To, "All")
	}
	if remarks.Text != "Situation report for Team Alpha" {
		t.Errorf("Text not correctly unmarshaled. Got: %s, Expected: %s", remarks.Text, "Situation report for Team Alpha")
	}
}

func TestRemarksSetters(t *testing.T) {
	// Given
	remarks := Remarks{}
	source := "Tactical Device"
	sourceID := "ANDROID-355161070128345"
	timestamp := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	to := "All"
	text := "Situation report for Team Alpha"

	// When
	result := remarks.
		SetSource(source).
		SetSourceID(sourceID).
		SetTime(&timestamp).
		SetTo(to).
		SetText(text)

	// Then
	if remarks.Source != source {
		t.Errorf("Source not set correctly. Got: %s, Expected: %s", remarks.Source, source)
	}
	if remarks.SourceID != sourceID {
		t.Errorf("SourceID not set correctly. Got: %s, Expected: %s", remarks.SourceID, sourceID)
	}
	if !remarks.Time.Equal(timestamp) {
		t.Errorf("Time not set correctly. Got: %v, Expected: %v", remarks.Time, timestamp)
	}
	if remarks.To != to {
		t.Errorf("To not set correctly. Got: %s, Expected: %s", remarks.To, to)
	}
	if remarks.Text != text {
		t.Errorf("Text not set correctly. Got: %s, Expected: %s", remarks.Text, text)
	}
	if result != &remarks {
		t.Errorf("Setter methods did not return the remarks instance for chaining")
	}
}
