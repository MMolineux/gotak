package cot

import (
	"encoding/xml"
	"testing"
)

func TestDetailMarshalXML(t *testing.T) {
	// Given
	detail := Detail{}
	detail.AddContact("Team Alpha")

	// When
	data, err := xml.Marshal(detail)
	if err != nil {
		t.Fatalf("Failed to marshal Detail to XML: %v", err)
	}

	// Then
	expected := `<detail><contact callsign="Team Alpha"></contact></detail>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestDetailMarshalXMLWithMultipleElements(t *testing.T) {
	// Given
	detail := Detail{}
	detail.AddContact("Team Alpha")
	detail.AddRemarks("Situation report")
	detail.AddStatus().SetBattery(75)
	detail.AddLink(&Link{UID: "ANDROID-355161070128345", Type: "a-f-G-U-C"})

	// When
	data, err := xml.Marshal(detail)
	if err != nil {
		t.Fatalf("Failed to marshal Detail with multiple elements to XML: %v", err)
	}

	// Then
	expected := `<detail><contact callsign="Team Alpha"></contact><remarks>Situation report</remarks><status battery="75"></status><link uid="ANDROID-355161070128345" type="a-f-G-U-C"></link></detail>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestDetailUnmarshalXML(t *testing.T) {
	// Given
	xmlData := `<detail><contact callsign="Team Alpha"></contact><remarks>Situation report</remarks><status battery="75"></status><link uid="ANDROID-355161070128345" type="a-f-G-U-C"></link></detail>`

	// When
	var detail Detail
	err := xml.Unmarshal([]byte(xmlData), &detail)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to Detail: %v", err)
	}

	// Then
	if detail.Contact == nil {
		t.Fatalf("Contact not correctly unmarshaled, got nil")
	}
	if detail.Contact.Callsign != "Team Alpha" {
		t.Errorf("Contact.Callsign not correctly unmarshaled. Got: %s, Expected: %s", detail.Contact.Callsign, "Team Alpha")
	}

	if detail.Remarks == nil {
		t.Fatalf("Remarks not correctly unmarshaled, got nil")
	}
	if detail.Remarks.Text != "Situation report" {
		t.Errorf("Remarks.Text not correctly unmarshaled. Got: %s, Expected: %s", detail.Remarks.Text, "Situation report")
	}

	if detail.Status == nil {
		t.Fatalf("Status not correctly unmarshaled, got nil")
	}
	if detail.Status.Battery != 75 {
		t.Errorf("Status.Battery not correctly unmarshaled. Got: %d, Expected: %d", detail.Status.Battery, 75)
	}

	if len(detail.Links) != 1 {
		t.Fatalf("Links not correctly unmarshaled. Expected 1 link, got %d", len(detail.Links))
	}
	if detail.Links[0].UID != "ANDROID-355161070128345" {
		t.Errorf("Links[0].UID not correctly unmarshaled. Got: %s, Expected: %s", detail.Links[0].UID, "ANDROID-355161070128345")
	}
	if detail.Links[0].Type != "a-f-G-U-C" {
		t.Errorf("Links[0].Type not correctly unmarshaled. Got: %s, Expected: %s", detail.Links[0].Type, "a-f-G-U-C")
	}
}

func TestDetailAddContact(t *testing.T) {
	// Given
	detail := Detail{}
	callsign := "Team Alpha"

	// When
	contact := detail.AddContact(callsign)

	// Then
	if detail.Contact == nil {
		t.Fatalf("AddContact failed to add a contact to the detail")
	}
	if detail.Contact.Callsign != callsign {
		t.Errorf("Contact.Callsign not set correctly. Got: %s, Expected: %s", detail.Contact.Callsign, callsign)
	}
	if contact != detail.Contact {
		t.Errorf("AddContact did not return the created contact")
	}
}

func TestDetailAddTakv(t *testing.T) {
	// Given
	detail := Detail{}
	platform := "ANDROID"
	version := "4.1.0.0"

	// When
	takv := detail.AddTakv(platform, version)

	// Then
	if detail.Takv == nil {
		t.Fatalf("AddTakv failed to add a takv to the detail")
	}
	if detail.Takv.Platform != platform {
		t.Errorf("Takv.Platform not set correctly. Got: %s, Expected: %s", detail.Takv.Platform, platform)
	}
	if detail.Takv.Version != version {
		t.Errorf("Takv.Version not set correctly. Got: %s, Expected: %s", detail.Takv.Version, version)
	}
	if takv != detail.Takv {
		t.Errorf("AddTakv did not return the created takv")
	}
}

func TestDetailAddRemarks(t *testing.T) {
	// Given
	detail := Detail{}
	text := "Situation report for Team Alpha"

	// When
	remarks := detail.AddRemarks(text)

	// Then
	if detail.Remarks == nil {
		t.Fatalf("AddRemarks failed to add remarks to the detail")
	}
	if detail.Remarks.Text != text {
		t.Errorf("Remarks.Text not set correctly. Got: %s, Expected: %s", detail.Remarks.Text, text)
	}
	if remarks != detail.Remarks {
		t.Errorf("AddRemarks did not return the created remarks")
	}
}

func TestDetailAddStatus(t *testing.T) {
	// Given
	detail := Detail{}

	// When
	status := detail.AddStatus()
	status.SetBattery(75)

	// Then
	if detail.Status == nil {
		t.Fatalf("AddStatus failed to add a status to the detail")
	}
	if detail.Status.Battery != 75 {
		t.Errorf("Status.Battery not set correctly. Got: %d, Expected: %d", detail.Status.Battery, 75)
	}
	if status != detail.Status {
		t.Errorf("AddStatus did not return the created status")
	}
}

func TestDetailAddTrack(t *testing.T) {
	// Given
	detail := Detail{}

	// When
	track := detail.AddTrack()
	track.Course = 180.5
	track.Speed = 25.3

	// Then
	if detail.Track == nil {
		t.Fatalf("AddTrack failed to add a track to the detail")
	}
	if detail.Track.Course != 180.5 {
		t.Errorf("Track.Course not set correctly. Got: %f, Expected: %f", detail.Track.Course, 180.5)
	}
	if detail.Track.Speed != 25.3 {
		t.Errorf("Track.Speed not set correctly. Got: %f, Expected: %f", detail.Track.Speed, 25.3)
	}
	if track != detail.Track {
		t.Errorf("AddTrack did not return the created track")
	}
}

func TestDetailAddPrecisionLocation(t *testing.T) {
	// Given
	detail := Detail{}
	altSrc := "DTED2"

	// When
	precisionLocation := detail.AddPrecisionLocation(altSrc)
	precisionLocation.SetGeoPointSrc("GPS")

	// Then
	if detail.PrecisionLocation == nil {
		t.Fatalf("AddPrecisionLocation failed to add a precision location to the detail")
	}
	if detail.PrecisionLocation.AltSrc != altSrc {
		t.Errorf("PrecisionLocation.AltSrc not set correctly. Got: %s, Expected: %s", detail.PrecisionLocation.AltSrc, altSrc)
	}
	if detail.PrecisionLocation.GeoPointSrc != "GPS" {
		t.Errorf("PrecisionLocation.GeoPointSrc not set correctly. Got: %s, Expected: %s", detail.PrecisionLocation.GeoPointSrc, "GPS")
	}
	if precisionLocation != detail.PrecisionLocation {
		t.Errorf("AddPrecisionLocation did not return the created precision location")
	}
}

func TestDetailAddShape(t *testing.T) {
	// Given
	detail := Detail{}

	// When
	shape := detail.AddShape()
	shape.AddEllipse(100.5, 50.5, 45.0)

	// Then
	if detail.Shape == nil {
		t.Fatalf("AddShape failed to add a shape to the detail")
	}
	if detail.Shape.Ellipse == nil {
		t.Fatalf("Shape.Ellipse is nil but expected to be set")
	}
	if detail.Shape.Ellipse.Major != 100.5 {
		t.Errorf("Shape.Ellipse.Major not set correctly. Got: %f, Expected: %f", detail.Shape.Ellipse.Major, 100.5)
	}
	if detail.Shape.Ellipse.Minor != 50.5 {
		t.Errorf("Shape.Ellipse.Minor not set correctly. Got: %f, Expected: %f", detail.Shape.Ellipse.Minor, 50.5)
	}
	if detail.Shape.Ellipse.Angle != 45.0 {
		t.Errorf("Shape.Ellipse.Angle not set correctly. Got: %f, Expected: %f", detail.Shape.Ellipse.Angle, 45.0)
	}
	if shape != detail.Shape {
		t.Errorf("AddShape did not return the created shape")
	}
}

func TestDetailAddLink(t *testing.T) {
	// Given
	detail := Detail{}
	uid := "ANDROID-355161070128345"
	linkType := "a-f-G-U-C"

	// When
	link := detail.AddLink(&Link{UID: uid, Type: linkType})

	// Then
	if len(detail.Links) != 1 {
		t.Fatalf("AddLink failed to add a link to the detail. Expected 1 link, got %d", len(detail.Links))
	}
	if detail.Links[0].UID != uid {
		t.Errorf("Link.UID not set correctly. Got: %s, Expected: %s", detail.Links[0].UID, uid)
	}
	if detail.Links[0].Type != linkType {
		t.Errorf("Link.Type not set correctly. Got: %s, Expected: %s", detail.Links[0].Type, linkType)
	}
	if link != detail.Links[0] {
		t.Errorf("AddLink did not return the created link")
	}
}

func TestDetailAddPointLink(t *testing.T) {
	// Given
	detail := Detail{}
	point := "38.870,-77.055"

	// When
	link := detail.AddPointLink(point)

	// Then
	if len(detail.Links) != 1 {
		t.Fatalf("AddPointLink failed to add a link to the detail. Expected 1 link, got %d", len(detail.Links))
	}
	if detail.Links[0].Point != point {
		t.Errorf("Link.Point not set correctly. Got: %s, Expected: %s", detail.Links[0].Point, point)
	}
	if link != detail.Links[0] {
		t.Errorf("AddPointLink did not return the created link")
	}
}

func TestDetailSetColor(t *testing.T) {
	// Given
	detail := Detail{}
	value := int64(0xff0000ff)

	// When
	color := detail.SetColor(value)

	// Then
	if detail.Color == nil {
		t.Fatalf("SetColor failed to set a color for the detail")
	}
	if detail.Color.Value != value {
		t.Errorf("Color.Value not set correctly. Got: %d, Expected: %d", detail.Color.Value, value)
	}
	if color != detail.Color {
		t.Errorf("SetColor did not return the created color")
	}
}

func TestDetailSetStrokeColor(t *testing.T) {
	// Given
	detail := Detail{}
	value := int64(0xff0000ff)

	// When
	strokeColor := detail.SetStrokeColor(value)

	// Then
	if detail.StrokeColor == nil {
		t.Fatalf("SetStrokeColor failed to set a stroke color for the detail")
	}
	if detail.StrokeColor.Value != value {
		t.Errorf("StrokeColor.Value not set correctly. Got: %d, Expected: %d", detail.StrokeColor.Value, value)
	}
	if strokeColor != detail.StrokeColor {
		t.Errorf("SetStrokeColor did not return the created stroke color")
	}
}

func TestDetailSetStrokeWeight(t *testing.T) {
	// Given
	detail := Detail{}
	value := 2.0

	// When
	strokeWeight := detail.SetStrokeWeight(value)

	// Then
	if detail.StrokeWeight == nil {
		t.Fatalf("SetStrokeWeight failed to set a stroke weight for the detail")
	}
	if detail.StrokeWeight.Value != value {
		t.Errorf("StrokeWeight.Value not set correctly. Got: %f, Expected: %f", detail.StrokeWeight.Value, value)
	}
	if strokeWeight != detail.StrokeWeight {
		t.Errorf("SetStrokeWeight did not return the created stroke weight")
	}
}

func TestDetailSetFillColor(t *testing.T) {
	// Given
	detail := Detail{}
	value := int64(0x7f00ff00)

	// When
	fillColor := detail.SetFillColor(value)

	// Then
	if detail.FillColor == nil {
		t.Fatalf("SetFillColor failed to set a fill color for the detail")
	}
	if detail.FillColor.Value != value {
		t.Errorf("FillColor.Value not set correctly. Got: %d, Expected: %d", detail.FillColor.Value, value)
	}
	if fillColor != detail.FillColor {
		t.Errorf("SetFillColor did not return the created fill color")
	}
}

func TestDetailSetLabelsOn(t *testing.T) {
	// Given
	detail := Detail{}
	value := true

	// When
	labelsOn := detail.SetLabelsOn(value)

	// Then
	if detail.LabelsOn == nil {
		t.Fatalf("SetLabelsOn failed to set a labels_on for the detail")
	}
	if detail.LabelsOn.Value != value {
		t.Errorf("LabelsOn.Value not set correctly. Got: %t, Expected: %t", detail.LabelsOn.Value, value)
	}
	if labelsOn != detail.LabelsOn {
		t.Errorf("SetLabelsOn did not return the created labels_on")
	}
}

func TestDetailAddArchive(t *testing.T) {
	// Given
	detail := Detail{}

	// When
	archive := detail.AddArchive()

	// Then
	if detail.Archive == nil {
		t.Fatalf("AddArchive failed to add an archive to the detail")
	}
	if archive != detail.Archive {
		t.Errorf("AddArchive did not return the created archive")
	}
}

func TestDetailSetTog(t *testing.T) {
	// Given
	detail := Detail{}
	value := true

	// When
	tog := detail.SetTog(value)

	// Then
	if detail.Tog == nil {
		t.Fatalf("SetTog failed to set a tog for the detail")
	}
	if detail.Tog.Value != value {
		t.Errorf("Tog.Value not set correctly. Got: %t, Expected: %t", detail.Tog.Value, value)
	}
	if tog != detail.Tog {
		t.Errorf("SetTog did not return the created tog")
	}
}
