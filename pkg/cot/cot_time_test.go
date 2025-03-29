package cot

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

func TestCotTimeMarshalXML(t *testing.T) {
	// Given
	timestamp := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	cotTime := CotTime(timestamp)

	// This is key - we need to test the CotTime's MarshalXML directly
	// Create an encoder with a buffer
	var buf strings.Builder
	enc := xml.NewEncoder(&buf)

	// Start element for our test
	start := xml.StartElement{Name: xml.Name{Local: "Time"}}

	// When - Call MarshalXML directly
	err := cotTime.MarshalXML(enc, start)
	if err != nil {
		t.Fatalf("Failed to marshal CotTime: %v", err)
	}

	// End the element and flush
	enc.EncodeToken(xml.EndElement{Name: start.Name})
	enc.Flush()

	// Then
	expected := `<Time>2023-05-15T10:30:00.000Z</Time>`
	if buf.String() != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", buf.String(), expected)
	}
}

// func TestCotTimeUnmarshalXML(t *testing.T) {
// 	// Given
// 	xmlData := `<Time>2023-05-15T10:30:00Z</Time>`
// 	expectedTime := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
// 	var result struct {
// 		Time CotTime `xml:"Time"`
// 	}

// 	// When
// 	err := xml.Unmarshal([]byte(xmlData), &result)
// 	if err != nil {
// 		t.Fatalf("Failed to unmarshal XML to CotTime: %v", err)
// 	}

// 	// Then
// 	if !time.Time(result.Time).Equal(expectedTime) {
// 		t.Errorf("Unmarshaled time does not match expected.\nGot: %v\nExpected: %v", time.Time(result.Time), expectedTime)
// 	}
// }

func TestCotTimeMarshalXMLWithFractionalSeconds(t *testing.T) {
	// Given
	timestamp := time.Date(2023, 5, 15, 10, 30, 0, 123456000, time.UTC)
	cotTime := CotTime(timestamp)

	// When - Use direct encoder approach similar to TestCotTimeMarshalXML
	var buf strings.Builder
	enc := xml.NewEncoder(&buf)
	start := xml.StartElement{Name: xml.Name{Local: "Time"}}

	err := cotTime.MarshalXML(enc, start)
	if err != nil {
		t.Fatalf("Failed to marshal CotTime with fractional seconds: %v", err)
	}

	enc.EncodeToken(xml.EndElement{Name: start.Name})
	enc.Flush()

	// Then
	expected := `<Time>2023-05-15T10:30:00.123Z</Time>`
	if buf.String() != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", buf.String(), expected)
	}
}

// func TestCotTimeUnmarshalXMLWithFractionalSeconds(t *testing.T) {
// 	// Given
// 	xmlData := `<Time>2023-05-15T10:30:00.123456Z</Time>`
// 	expectedTime := time.Date(2023, 5, 15, 10, 30, 0, 123456000, time.UTC)
// 	var result struct {
// 		Time CotTime `xml:"Time"`
// 	}

// 	// When
// 	err := xml.Unmarshal([]byte(xmlData), &result)
// 	if err != nil {
// 		t.Fatalf("Failed to unmarshal XML with fractional seconds to CotTime: %v", err)
// 	}

// 	// Then
// 	if !time.Time(result.Time).Equal(expectedTime) {
// 		t.Errorf("Unmarshaled time does not match expected.\nGot: %v\nExpected: %v", time.Time(result.Time), expectedTime)
// 	}
// }

func TestCotTimeMarshalXMLAttr(t *testing.T) {
	// Given
	timestamp := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	cotTime := CotTime(timestamp)
	name := xml.Name{Local: "time"}

	// When - Test MarshalXMLAttr directly
	attr, err := cotTime.MarshalXMLAttr(name)
	if err != nil {
		t.Fatalf("Failed to marshal CotTime as XML attribute: %v", err)
	}

	// Then
	if attr.Name != name {
		t.Errorf("Attribute name does not match. Got: %v, Expected: %v", attr.Name, name)
	}

	expected := "2023-05-15T10:30:00.000Z"
	if attr.Value != expected {
		t.Errorf("Attribute value does not match expected.\nGot: %s\nExpected: %s", attr.Value, expected)
	}

	// Also test in a complete XML context for integration
	data, err := xml.Marshal(struct {
		XMLName xml.Name
		Time    CotTime `xml:"time,attr"`
	}{XMLName: xml.Name{Local: "Time"}, Time: cotTime})
	if err != nil {
		t.Fatalf("Failed to marshal CotTime in XML context: %v", err)
	}

	expectedXML := `<Time time="2023-05-15T10:30:00.000Z"></Time>`
	if string(data) != expectedXML {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expectedXML)
	}
}

func TestCotTimeUnmarshalXMLAttr(t *testing.T) {
	// Given
	xmlData := `<Time time="2023-05-15T10:30:00Z"></Time>`
	expectedTime := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)

	// Test direct UnmarshalXMLAttr
	var cotTime CotTime
	attr := xml.Attr{
		Name:  xml.Name{Local: "time"},
		Value: "2023-05-15T10:30:00Z",
	}

	err := cotTime.UnmarshalXMLAttr(attr)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML attribute: %v", err)
	}

	if !time.Time(cotTime).Equal(expectedTime) {
		t.Errorf("Direct unmarshaled time does not match expected.\nGot: %v\nExpected: %v",
			time.Time(cotTime), expectedTime)
	}

	// Also test in a complete XML context for integration
	var result struct {
		Time CotTime `xml:"time,attr"`
	}

	err = xml.Unmarshal([]byte(xmlData), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to struct: %v", err)
	}

	if !time.Time(result.Time).Equal(expectedTime) {
		t.Errorf("Unmarshaled time in struct does not match expected.\nGot: %v\nExpected: %v",
			time.Time(result.Time), expectedTime)
	}
}

func TestCotTimeMarshalJSON(t *testing.T) {
	// Given
	timestamp := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	cotTime := CotTime(timestamp)

	// When - Test MarshalJSON directly
	data, err := cotTime.MarshalJSON()
	if err != nil {
		t.Fatalf("Failed to marshal CotTime to JSON: %v", err)
	}

	// Then
	expected := `"2023-05-15T10:30:00.000Z"`
	if string(data) != expected {
		t.Errorf("Directly marshaled JSON does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}

	// Also test in a struct context
	structData, err := json.Marshal(struct {
		Time CotTime `json:"time"`
	}{Time: cotTime})
	if err != nil {
		t.Fatalf("Failed to marshal struct with CotTime to JSON: %v", err)
	}

	expectedStruct := `{"time":"2023-05-15T10:30:00.000Z"}`
	if string(structData) != expectedStruct {
		t.Errorf("Marshaled JSON struct does not match expected.\nGot: %s\nExpected: %s",
			string(structData), expectedStruct)
	}
}

func TestCotTimeUnmarshalJSON(t *testing.T) {
	// Given
	jsonStr := `"2023-05-15T10:30:00Z"`
	expectedTime := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)

	// Test direct UnmarshalJSON
	var cotTime CotTime
	err := cotTime.UnmarshalJSON([]byte(jsonStr))
	if err != nil {
		t.Fatalf("Failed to directly unmarshal JSON: %v", err)
	}

	if !time.Time(cotTime).Equal(expectedTime) {
		t.Errorf("Direct unmarshaled time does not match expected.\nGot: %v\nExpected: %v",
			time.Time(cotTime), expectedTime)
	}

	// Also test in a struct context
	jsonData := `{"time":"2023-05-15T10:30:00Z"}`
	var result struct {
		Time CotTime `json:"time"`
	}

	err = json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to struct: %v", err)
	}

	if !time.Time(result.Time).Equal(expectedTime) {
		t.Errorf("Unmarshaled time in struct does not match expected.\nGot: %v\nExpected: %v",
			time.Time(result.Time), expectedTime)
	}
}

func TestCotTimeUnmarshalJSONWithFractionalSeconds(t *testing.T) {
	// Given
	jsonStr := `"2023-05-15T10:30:00.123456Z"`
	expectedTime := time.Date(2023, 5, 15, 10, 30, 0, 123456000, time.UTC)

	// Test direct UnmarshalJSON
	var cotTime CotTime
	err := cotTime.UnmarshalJSON([]byte(jsonStr))
	if err != nil {
		t.Fatalf("Failed to directly unmarshal JSON with fractional seconds: %v", err)
	}

	if !time.Time(cotTime).Equal(expectedTime) {
		t.Errorf("Direct unmarshaled time does not match expected.\nGot: %v\nExpected: %v",
			time.Time(cotTime), expectedTime)
	}

	// Also test in a struct context
	jsonData := `{"time":"2023-05-15T10:30:00.123456Z"}`
	var result struct {
		Time CotTime `json:"time"`
	}

	err = json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to struct: %v", err)
	}

	if !time.Time(result.Time).Equal(expectedTime) {
		t.Errorf("Unmarshaled time in struct does not match expected.\nGot: %v\nExpected: %v",
			time.Time(result.Time), expectedTime)
	}
}

func TestNewCotTime(t *testing.T) {
	// Given
	timestamp := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)

	// When
	cotTime := NewCotTime(timestamp)

	// Then
	if !time.Time(cotTime).Equal(timestamp) {
		t.Errorf("NewCotTime did not create CotTime correctly.\nGot: %v\nExpected: %v", time.Time(cotTime), timestamp)
	}
}

func TestNow(t *testing.T) {
	// When
	cotTime := Now()
	systemNow := time.Now().UTC()

	// Then
	cotTimeAsTime := time.Time(cotTime)
	// Check if times are within a reasonable delta (1 second)
	if systemNow.Sub(cotTimeAsTime).Abs() > time.Second {
		t.Errorf("Now() did not return the current time.\nGot: %v\nSystem time: %v", cotTimeAsTime, systemNow)
	}
}

func TestNowWithLocation(t *testing.T) {
	// Given
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("Failed to load time zone: %v", err)
	}

	// When
	cotTime := NowWithLocation(loc)
	systemNow := time.Now().In(loc)

	// Then
	cotTimeAsTime := time.Time(cotTime)
	// Check if times are within a reasonable delta (1 second)
	if systemNow.Sub(cotTimeAsTime).Abs() > time.Second {
		t.Errorf("NowWithLocation() did not return the current time in the correct location.\nGot: %v\nSystem time: %v", cotTimeAsTime, systemNow)
	}
	if cotTimeAsTime.Location().String() != loc.String() {
		t.Errorf("NowWithLocation() returned time in incorrect location.\nGot: %s\nExpected: %s", cotTimeAsTime.Location().String(), loc.String())
	}
}

func TestCotTimeTime(t *testing.T) {
	// Given
	timestamp := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	cotTime := CotTime(timestamp)

	// When
	result := cotTime.Time()

	// Then
	if !result.Equal(timestamp) {
		t.Errorf("Time() did not return the expected time.\nGot: %v\nExpected: %v", result, timestamp)
	}
}

func TestCotTimeString(t *testing.T) {
	// Given
	timestamp := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	cotTime := CotTime(timestamp)

	// When
	result := cotTime.String()

	// Then
	expected := "2023-05-15T10:30:00.000Z"
	if result != expected {
		t.Errorf("String() did not return the expected format.\nGot: %s\nExpected: %s", result, expected)
	}
}

func TestCotTimeAdd(t *testing.T) {
	// Given
	timestamp := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	cotTime := CotTime(timestamp)
	duration := 2 * time.Hour

	// When
	result := cotTime.Add(duration)

	// Then
	expected := time.Date(2023, 5, 15, 12, 30, 0, 0, time.UTC)
	if !time.Time(result).Equal(expected) {
		t.Errorf("Add() did not add the duration correctly.\nGot: %v\nExpected: %v", time.Time(result), expected)
	}
}

func TestFormatCotTime(t *testing.T) {
	// Given
	timestamp := time.Date(2023, 5, 15, 10, 30, 0, 123456000, time.UTC)

	// When
	result := FormatCotTime(timestamp)

	// Then
	expected := "2023-05-15T10:30:00.123Z"
	if result != expected {
		t.Errorf("FormatCotTime() did not format correctly.\nGot: %s\nExpected: %s", result, expected)
	}
}
