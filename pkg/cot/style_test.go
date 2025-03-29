package cot

import (
	"encoding/xml"
	"testing"
)

func TestStyleMarshalXML(t *testing.T) {
	// Given
	style := Style{}
	style.SetLineStyle(&LineStyle{Color: "ff0000ff", Width: 2.0})
	style.SetPolyStyle(&PolyStyle{Color: "7f00ff00"})

	// When
	data, err := xml.Marshal(style)
	if err != nil {
		t.Fatalf("Failed to marshal Style to XML: %v", err)
	}

	// Then
	expected := `<Style><LineStyle><color>ff0000ff</color><width>2</width></LineStyle><PolyStyle><color>7f00ff00</color></PolyStyle></Style>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestStyleUnmarshalXML(t *testing.T) {
	// Given
	xmlData := `<Style><LineStyle><color>ff0000ff</color><width>2</width></LineStyle><PolyStyle><color>7f00ff00</color></PolyStyle></Style>`

	// When
	var style Style
	err := xml.Unmarshal([]byte(xmlData), &style)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to Style: %v", err)
	}

	// Then
	if style.LineStyle == nil {
		t.Fatalf("LineStyle not correctly unmarshaled, got nil")
	}
	if style.LineStyle.Color != "ff0000ff" {
		t.Errorf("LineStyle.Color not correctly unmarshaled. Got: %s, Expected: %s", style.LineStyle.Color, "ff0000ff")
	}
	if style.LineStyle.Width != 2.0 {
		t.Errorf("LineStyle.Width not correctly unmarshaled. Got: %f, Expected: %f", style.LineStyle.Width, 2.0)
	}
	if style.PolyStyle == nil {
		t.Fatalf("PolyStyle not correctly unmarshaled, got nil")
	}
	if style.PolyStyle.Color != "7f00ff00" {
		t.Errorf("PolyStyle.Color not correctly unmarshaled. Got: %s, Expected: %s", style.PolyStyle.Color, "7f00ff00")
	}
}

func TestStyleSetters(t *testing.T) {
	// Given
	style := Style{}
	lineStyle := &LineStyle{Color: "ff0000ff", Width: 2.0}
	polyStyle := &PolyStyle{Color: "7f00ff00"}

	// When
	style.SetLineStyle(lineStyle)
	style.SetPolyStyle(polyStyle)

	// Then
	if style.LineStyle != lineStyle {
		t.Errorf("LineStyle not set correctly. Expected the same pointer")
	}
	if style.PolyStyle != polyStyle {
		t.Errorf("PolyStyle not set correctly. Expected the same pointer")
	}
}

func TestLineStyleMarshalXML(t *testing.T) {
	// Given
	lineStyle := LineStyle{
		Color: "ff0000ff",
		Width: 2.0,
	}

	// When
	data, err := xml.Marshal(lineStyle)
	if err != nil {
		t.Fatalf("Failed to marshal LineStyle to XML: %v", err)
	}

	// Then
	expected := `<LineStyle><color>ff0000ff</color><width>2</width></LineStyle>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestLineStyleSetters(t *testing.T) {
	// Given
	lineStyle := LineStyle{}
	color := "ff0000ff"
	width := 2.0

	// When
	result := lineStyle.SetColor(color).SetWidth(width)

	// Then
	if lineStyle.Color != color {
		t.Errorf("Color not set correctly. Got: %s, Expected: %s", lineStyle.Color, color)
	}
	if lineStyle.Width != width {
		t.Errorf("Width not set correctly. Got: %f, Expected: %f", lineStyle.Width, width)
	}
	if result != &lineStyle {
		t.Errorf("Setter method did not return the instance for chaining")
	}
}

func TestPolyStyleMarshalXML(t *testing.T) {
	// Given
	polyStyle := PolyStyle{
		Color: "7f00ff00",
	}

	// When
	data, err := xml.Marshal(polyStyle)
	if err != nil {
		t.Fatalf("Failed to marshal PolyStyle to XML: %v", err)
	}

	// Then
	expected := `<PolyStyle><color>7f00ff00</color></PolyStyle>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestPolyStyleSetters(t *testing.T) {
	// Given
	polyStyle := PolyStyle{}
	color := "7f00ff00"

	// When
	result := polyStyle.SetColor(color)

	// Then
	if polyStyle.Color != color {
		t.Errorf("Color not set correctly. Got: %s, Expected: %s", polyStyle.Color, color)
	}
	if result != &polyStyle {
		t.Errorf("Setter method did not return the instance for chaining")
	}
}

func TestColorMarshalXML(t *testing.T) {
	// Given
	color := Color{
		Value: 0xff0000ff,
	}

	// When
	data, err := xml.Marshal(color)
	if err != nil {
		t.Fatalf("Failed to marshal Color to XML: %v", err)
	}

	// Then
	expected := `<color argb="4278190335"></color>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestColorSetters(t *testing.T) {
	// Given
	color := Color{}
	value := int64(0xff0000ff)

	// When
	result := color.SetValue(value)

	// Then
	if color.Value != value {
		t.Errorf("Value not set correctly. Got: %d, Expected: %d", color.Value, value)
	}
	if result != &color {
		t.Errorf("Setter method did not return the instance for chaining")
	}
}

func TestStrokeColorSetters(t *testing.T) {
	// Given
	strokeColor := StrokeColor{}
	value := int64(0xff0000ff)

	// When
	result := strokeColor.SetValue(value)

	// Then
	if strokeColor.Value != value {
		t.Errorf("Value not set correctly. Got: %d, Expected: %d", strokeColor.Value, value)
	}
	if result != &strokeColor {
		t.Errorf("Setter method did not return the instance for chaining")
	}
}

func TestStrokeWeightSetters(t *testing.T) {
	// Given
	strokeWeight := StrokeWeight{}
	value := 3.5

	// When
	result := strokeWeight.SetValue(value)

	// Then
	if strokeWeight.Value != value {
		t.Errorf("Value not set correctly. Got: %f, Expected: %f", strokeWeight.Value, value)
	}
	if result != &strokeWeight {
		t.Errorf("Setter method did not return the instance for chaining")
	}
}

func TestFillColorSetters(t *testing.T) {
	// Given
	fillColor := FillColor{}
	value := int64(0x7f00ff00)

	// When
	result := fillColor.SetValue(value)

	// Then
	if fillColor.Value != value {
		t.Errorf("Value not set correctly. Got: %d, Expected: %d", fillColor.Value, value)
	}
	if result != &fillColor {
		t.Errorf("Setter method did not return the instance for chaining")
	}
}

func TestLabelsOnSetters(t *testing.T) {
	// Given
	labelsOn := LabelsOn{}
	value := true

	// When
	result := labelsOn.SetValue(value)

	// Then
	if labelsOn.Value != value {
		t.Errorf("Value not set correctly. Got: %t, Expected: %t", labelsOn.Value, value)
	}
	if result != &labelsOn {
		t.Errorf("Setter method did not return the instance for chaining")
	}
}

func TestStyleMethodChaining(t *testing.T) {
	// Given
	style := Style{}
	lineStyle := LineStyle{}
	polyStyle := PolyStyle{}

	// When
	style.SetLineStyle(lineStyle.SetColor("ff0000ff").SetWidth(2.0)).
		SetPolyStyle(polyStyle.SetColor("7f00ff00"))

	// Then
	if style.LineStyle.Color != "ff0000ff" {
		t.Errorf("LineStyle.Color not set correctly. Got: %s, Expected: %s", style.LineStyle.Color, "ff0000ff")
	}
	if style.LineStyle.Width != 2.0 {
		t.Errorf("LineStyle.Width not set correctly. Got: %f, Expected: %f", style.LineStyle.Width, 2.0)
	}
	if style.PolyStyle.Color != "7f00ff00" {
		t.Errorf("PolyStyle.Color not set correctly. Got: %s, Expected: %s", style.PolyStyle.Color, "7f00ff00")
	}
}
