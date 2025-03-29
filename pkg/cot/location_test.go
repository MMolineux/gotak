package cot

import (
	"encoding/xml"
	"testing"
)

func TestPrecisionLocationMarshalXML(t *testing.T) {
	// Given
	precisionLocation := PrecisionLocation{
		GeoPointSrc:       "GPS",
		AltSrc:            "DTED2",
		PreciseImageFile:  "image.jpg",
		PreciseImageFileX: 150.5,
		PreciseImageFileY: 200.75,
	}

	// When
	data, err := xml.Marshal(precisionLocation)
	if err != nil {
		t.Fatalf("Failed to marshal PrecisionLocation to XML: %v", err)
	}

	// Then
	expected := `<precisionlocation geopointsrc="GPS" altsrc="DTED2" PRECISE_IMAGE_FILE="image.jpg" PRECISE_IMAGE_FILE_X="150.5" PRECISE_IMAGE_FILE_Y="200.75"></precisionlocation>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestPrecisionLocationMarshalXMLWithMinimalFields(t *testing.T) {
	// Given
	precisionLocation := PrecisionLocation{
		AltSrc: "DTED2",
	}

	// When
	data, err := xml.Marshal(precisionLocation)
	if err != nil {
		t.Fatalf("Failed to marshal PrecisionLocation to XML: %v", err)
	}

	// Then
	expected := `<precisionlocation altsrc="DTED2"></precisionlocation>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestPrecisionLocationUnmarshalXML(t *testing.T) {
	// Given
	xmlData := `<precisionlocation geopointsrc="GPS" altsrc="DTED2" PRECISE_IMAGE_FILE="image.jpg" PRECISE_IMAGE_FILE_X="150.5" PRECISE_IMAGE_FILE_Y="200.75"></precisionlocation>`

	// When
	var precisionLocation PrecisionLocation
	err := xml.Unmarshal([]byte(xmlData), &precisionLocation)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to PrecisionLocation: %v", err)
	}

	// Then
	if precisionLocation.GeoPointSrc != "GPS" {
		t.Errorf("GeoPointSrc not correctly unmarshaled. Got: %s, Expected: %s", precisionLocation.GeoPointSrc, "GPS")
	}
	if precisionLocation.AltSrc != "DTED2" {
		t.Errorf("AltSrc not correctly unmarshaled. Got: %s, Expected: %s", precisionLocation.AltSrc, "DTED2")
	}
	if precisionLocation.PreciseImageFile != "image.jpg" {
		t.Errorf("PreciseImageFile not correctly unmarshaled. Got: %s, Expected: %s", precisionLocation.PreciseImageFile, "image.jpg")
	}
	if precisionLocation.PreciseImageFileX != 150.5 {
		t.Errorf("PreciseImageFileX not correctly unmarshaled. Got: %f, Expected: %f", precisionLocation.PreciseImageFileX, 150.5)
	}
	if precisionLocation.PreciseImageFileY != 200.75 {
		t.Errorf("PreciseImageFileY not correctly unmarshaled. Got: %f, Expected: %f", precisionLocation.PreciseImageFileY, 200.75)
	}
}

func TestPrecisionLocationSetters(t *testing.T) {
	// Given
	precisionLocation := PrecisionLocation{}
	geoPointSrc := "GPS"
	altSrc := "DTED2"
	preciseImageFile := "image.jpg"
	preciseImageFileX := 150.5
	preciseImageFileY := 200.75

	// When
	result := precisionLocation.
		SetGeoPointSrc(geoPointSrc).
		SetAltSrc(altSrc).
		SetPreciseImageFile(preciseImageFile).
		SetPreciseImageFileX(preciseImageFileX).
		SetPreciseImageFileY(preciseImageFileY)

	// Then
	if precisionLocation.GeoPointSrc != geoPointSrc {
		t.Errorf("GeoPointSrc not set correctly. Got: %s, Expected: %s", precisionLocation.GeoPointSrc, geoPointSrc)
	}
	if precisionLocation.AltSrc != altSrc {
		t.Errorf("AltSrc not set correctly. Got: %s, Expected: %s", precisionLocation.AltSrc, altSrc)
	}
	if precisionLocation.PreciseImageFile != preciseImageFile {
		t.Errorf("PreciseImageFile not set correctly. Got: %s, Expected: %s", precisionLocation.PreciseImageFile, preciseImageFile)
	}
	if precisionLocation.PreciseImageFileX != preciseImageFileX {
		t.Errorf("PreciseImageFileX not set correctly. Got: %f, Expected: %f", precisionLocation.PreciseImageFileX, preciseImageFileX)
	}
	if precisionLocation.PreciseImageFileY != preciseImageFileY {
		t.Errorf("PreciseImageFileY not set correctly. Got: %f, Expected: %f", precisionLocation.PreciseImageFileY, preciseImageFileY)
	}
	if result != &precisionLocation {
		t.Errorf("Setter methods did not return the precisionLocation instance for chaining")
	}
}
