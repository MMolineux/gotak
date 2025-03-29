package cot

import (
	"encoding/xml"
	"testing"
)

func TestShapeMarshalXML(t *testing.T) {
	// Given
	shape := Shape{}
	shape.AddEllipse(100.5, 50.5, 45.0)

	// When
	data, err := xml.Marshal(shape)
	if err != nil {
		t.Fatalf("Failed to marshal Shape to XML: %v", err)
	}

	// Then
	expected := `<shape><ellipse major="100.5" minor="50.5" angle="45"></ellipse></shape>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestShapeUnmarshalXML(t *testing.T) {
	// Given
	xmlData := `<shape><ellipse major="100.5" minor="50.5" angle="45"></ellipse></shape>`

	// When
	var shape Shape
	err := xml.Unmarshal([]byte(xmlData), &shape)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to Shape: %v", err)
	}

	// Then
	if shape.Ellipse == nil {
		t.Fatalf("Ellipse not correctly unmarshaled, got nil")
	}
	if shape.Ellipse.Major != 100.5 {
		t.Errorf("Ellipse.Major not correctly unmarshaled. Got: %f, Expected: %f", shape.Ellipse.Major, 100.5)
	}
	if shape.Ellipse.Minor != 50.5 {
		t.Errorf("Ellipse.Minor not correctly unmarshaled. Got: %f, Expected: %f", shape.Ellipse.Minor, 50.5)
	}
	if shape.Ellipse.Angle != 45.0 {
		t.Errorf("Ellipse.Angle not correctly unmarshaled. Got: %f, Expected: %f", shape.Ellipse.Angle, 45.0)
	}
}

func TestShapeWithPolylineMarshalXML(t *testing.T) {
	// Given
	shape := Shape{}
	shape.AddPolyline("38.870,-77.055 38.869,-77.059 38.868,-77.058")

	// When
	data, err := xml.Marshal(shape)
	if err != nil {
		t.Fatalf("Failed to marshal Shape with Polyline to XML: %v", err)
	}

	// Then
	expected := `<shape><polyline points="38.870,-77.055 38.869,-77.059 38.868,-77.058"></polyline></shape>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestShapeWithPolylineUnmarshalXML(t *testing.T) {
	// Given
	xmlData := `<shape><polyline points="38.870,-77.055 38.869,-77.059 38.868,-77.058"></polyline></shape>`

	// When
	var shape Shape
	err := xml.Unmarshal([]byte(xmlData), &shape)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to Shape with Polyline: %v", err)
	}

	// Then
	if shape.Polyline == nil {
		t.Fatalf("Polyline not correctly unmarshaled, got nil")
	}
	expected := "38.870,-77.055 38.869,-77.059 38.868,-77.058"
	if shape.Polyline.Points != expected {
		t.Errorf("Polyline.Points not correctly unmarshaled. Got: %s, Expected: %s", shape.Polyline.Points, expected)
	}
}

func TestShapeAddEllipse(t *testing.T) {
	// Given
	shape := Shape{}
	major := 100.5
	minor := 50.5
	angle := 45.0

	// When
	ellipse := shape.AddEllipse(major, minor, angle)

	// Then
	if shape.Ellipse == nil {
		t.Fatalf("AddEllipse failed to add an ellipse to the shape")
	}
	if shape.Ellipse.Major != major {
		t.Errorf("Ellipse.Major not set correctly. Got: %f, Expected: %f", shape.Ellipse.Major, major)
	}
	if shape.Ellipse.Minor != minor {
		t.Errorf("Ellipse.Minor not set correctly. Got: %f, Expected: %f", shape.Ellipse.Minor, minor)
	}
	if shape.Ellipse.Angle != angle {
		t.Errorf("Ellipse.Angle not set correctly. Got: %f, Expected: %f", shape.Ellipse.Angle, angle)
	}
	if ellipse != shape.Ellipse {
		t.Errorf("AddEllipse did not return the created ellipse")
	}
}

func TestShapeSetEllipse(t *testing.T) {
	// Given
	shape := Shape{}
	ellipse := &Ellipse{
		Major: 100.5,
		Minor: 50.5,
		Angle: 45.0,
	}

	// When
	result := shape.SetEllipse(ellipse)

	// Then
	if shape.Ellipse != ellipse {
		t.Errorf("Ellipse not set correctly. Expected the same pointer")
	}
	if result != &shape {
		t.Errorf("SetEllipse did not return the shape instance for chaining")
	}
}

func TestShapeAddPolyline(t *testing.T) {
	// Given
	shape := Shape{}
	points := "38.870,-77.055 38.869,-77.059 38.868,-77.058"

	// When
	polyline := shape.AddPolyline(points)

	// Then
	if shape.Polyline == nil {
		t.Fatalf("AddPolyline failed to add a polyline to the shape")
	}
	if shape.Polyline.Points != points {
		t.Errorf("Polyline.Points not set correctly. Got: %s, Expected: %s", shape.Polyline.Points, points)
	}
	if polyline != shape.Polyline {
		t.Errorf("AddPolyline did not return the created polyline")
	}
}

func TestShapeSetPolyline(t *testing.T) {
	// Given
	shape := Shape{}
	polyline := &Polyline{
		Points: "38.870,-77.055 38.869,-77.059 38.868,-77.058",
	}

	// When
	result := shape.SetPolyline(polyline)

	// Then
	if shape.Polyline != polyline {
		t.Errorf("Polyline not set correctly. Expected the same pointer")
	}
	if result != &shape {
		t.Errorf("SetPolyline did not return the shape instance for chaining")
	}
}

func TestShapeSetLink(t *testing.T) {
	// Given
	shape := Shape{}
	link := &Link{} // Assuming Link is defined somewhere

	// When
	result := shape.SetLink(link)

	// Then
	if shape.Link != link {
		t.Errorf("Link not set correctly. Expected the same pointer")
	}
	if result != &shape {
		t.Errorf("SetLink did not return the shape instance for chaining")
	}
}

func TestEllipseSetters(t *testing.T) {
	// Given
	ellipse := Ellipse{}
	major := 100.5
	minor := 50.5
	angle := 45.0

	// When
	result := ellipse.SetMajor(major).SetMinor(minor).SetAngle(angle)

	// Then
	if ellipse.Major != major {
		t.Errorf("Major not set correctly. Got: %f, Expected: %f", ellipse.Major, major)
	}
	if ellipse.Minor != minor {
		t.Errorf("Minor not set correctly. Got: %f, Expected: %f", ellipse.Minor, minor)
	}
	if ellipse.Angle != angle {
		t.Errorf("Angle not set correctly. Got: %f, Expected: %f", ellipse.Angle, angle)
	}
	if result != &ellipse {
		t.Errorf("Setter methods did not return the ellipse instance for chaining")
	}
}

func TestPolylineSetters(t *testing.T) {
	// Given
	polyline := Polyline{}
	points := "38.870,-77.055 38.869,-77.059 38.868,-77.058"

	// When
	result := polyline.SetPoints(points)

	// Then
	if polyline.Points != points {
		t.Errorf("Points not set correctly. Got: %s, Expected: %s", polyline.Points, points)
	}
	if result != &polyline {
		t.Errorf("SetPoints did not return the polyline instance for chaining")
	}
}
