package cot

import (
	"encoding/xml"
	"testing"
)

func TestPointMarshalXML(t *testing.T) {
	// Given
	point := Point{
		Lat: 38.8977,
		Lon: -77.0365,
	}
	point.SetCe(50.0).SetLe(75.0).SetHae(100.5)

	// When
	data, err := xml.Marshal(point)
	if err != nil {
		t.Fatalf("Failed to marshal Point to XML: %v", err)
	}

	// Then
	expected := `<point lat="38.8977" lon="-77.0365" hae="100.5" ce="50" le="75"></point>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestPointMarshalXMLWithoutOptionalFields(t *testing.T) {
	// Given
	point := Point{
		Lat: 38.8977,
		Lon: -77.0365,
	}

	// When
	data, err := xml.Marshal(point)
	if err != nil {
		t.Fatalf("Failed to marshal Point to XML: %v", err)
	}

	// Then
	expected := `<point lat="38.8977" lon="-77.0365" hae="9999999" ce="9999999" le="9999999"></point>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestPointMarshalXMLWithDefaultValue(t *testing.T) {
	dValue := DefaultValue
	// Given
	point := Point{
		Lat: 38.8977,
		Lon: -77.0365,
		Ce:  &dValue,
		Le:  &dValue,
	}

	// When
	data, err := xml.Marshal(point)
	if err != nil {
		t.Fatalf("Failed to marshal Point to XML: %v", err)
	}

	// Then
	// Assuming DefaultValue is 9999999 based on the MarshalXML implementation
	expected := `<point lat="38.8977" lon="-77.0365" hae="9999999" ce="9999999" le="9999999"></point>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestPointUnmarshalXML(t *testing.T) {
	// Given
	xmlData := `<point lat="38.8977" lon="-77.0365" hae="100.5" ce="50" le="75"></point>`

	// When
	var point Point
	err := xml.Unmarshal([]byte(xmlData), &point)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to Point: %v", err)
	}

	// Then
	if point.Lat != 38.8977 {
		t.Errorf("Lat not correctly unmarshaled. Got: %f, Expected: %f", point.Lat, 38.8977)
	}
	if point.Lon != -77.0365 {
		t.Errorf("Lon not correctly unmarshaled. Got: %f, Expected: %f", point.Lon, -77.0365)
	}
	if *point.Hae != 100.5 {
		t.Errorf("Hae not correctly unmarshaled. Got: %f, Expected: %f", *point.Hae, 100.5)
	}
	if *point.Ce != 50.0 {
		t.Errorf("Ce not correctly unmarshaled. Got: %f, Expected: %f", *point.Ce, 50.0)
	}
	if *point.Le != 75.0 {
		t.Errorf("Le not correctly unmarshaled. Got: %f, Expected: %f", *point.Le, 75.0)
	}
}

func TestPointSetters(t *testing.T) {
	// Given
	point := Point{}
	lat := 38.8977
	lon := -77.0365
	hae := 100.5
	ce := 50.0
	le := 75.0

	// When
	result := point.SetLat(lat).SetLon(lon).SetHae(hae).SetCe(ce).SetLe(le)

	// Then
	if point.Lat != lat {
		t.Errorf("Lat not set correctly. Got: %f, Expected: %f", point.Lat, lat)
	}
	if point.Lon != lon {
		t.Errorf("Lon not set correctly. Got: %f, Expected: %f", point.Lon, lon)
	}
	if *point.Hae != hae {
		t.Errorf("Hae not set correctly. Got: %f, Expected: %f", *point.Hae, hae)
	}
	if *point.Ce != ce {
		t.Errorf("Ce not set correctly. Got: %f, Expected: %f", *point.Ce, ce)
	}
	if *point.Le != le {
		t.Errorf("Le not set correctly. Got: %f, Expected: %f", *point.Le, le)
	}
	if result != &point {
		t.Errorf("Setter methods did not return the point instance for chaining")
	}
}
