package cot

import (
	"encoding/xml"
	"testing"
	"time"
)

func TestLinkMarshalXML(t *testing.T) {
	testTime := CotTime(time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC))
	// Given
	link := Link{
		UID:        "ANDROID-355161070128345",
		Type:       "a-f-G-U-C",
		Relation:   "p-p",
		Point:      "38.870,-77.055",
		URL:        "https://example.com/image.jpg",
		Remarks:    "Team Alpha",
		Production: &testTime,
		Version:    "2.0",
		Parent:     "Mission-123",
		Medium:     "chat",
	}

	// When
	data, err := xml.Marshal(link)
	if err != nil {
		t.Fatalf("Failed to marshal Link to XML: %v", err)
	}

	// Then
	expected := `<link uid="ANDROID-355161070128345" type="a-f-G-U-C" relation="p-p" point="38.870,-77.055" url="https://example.com/image.jpg" remarks="Team Alpha" production_time="2023-05-15T10:30:00.000Z" version="2.0" parent="Mission-123" medium="chat"></link>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestLinkMarshalXMLWithStyle(t *testing.T) {
	// Given
	link := Link{
		UID:  "ANDROID-355161070128345",
		Type: "a-f-G-U-C",
		Style: &Style{
			LineStyle: &LineStyle{
				Color: "ff0000ff",
				Width: 2.0,
			},
		},
	}

	// When
	data, err := xml.Marshal(link)
	if err != nil {
		t.Fatalf("Failed to marshal Link with Style to XML: %v", err)
	}

	// Then
	expected := `<link uid="ANDROID-355161070128345" type="a-f-G-U-C"><Style><LineStyle><color>ff0000ff</color><width>2</width></LineStyle></Style></link>`
	if string(data) != expected {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestLinkUnmarshalXML(t *testing.T) {
	// Given
	xmlData := `<link uid="ANDROID-355161070128345" type="a-f-G-U-C" relation="p-p" point="38.870,-77.055" url="https://example.com/image.jpg" remarks="Team Alpha" production_time="2023-05-15T10:30:00.000Z" version="2.0" parent="Mission-123" medium="chat"></link>`

	// When
	var link Link
	err := xml.Unmarshal([]byte(xmlData), &link)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to Link: %v", err)
	}

	// Then
	if link.UID != "ANDROID-355161070128345" {
		t.Errorf("UID not correctly unmarshaled. Got: %s, Expected: %s", link.UID, "ANDROID-355161070128345")
	}
	if link.Type != "a-f-G-U-C" {
		t.Errorf("Type not correctly unmarshaled. Got: %s, Expected: %s", link.Type, "a-f-G-U-C")
	}
	if link.Relation != "p-p" {
		t.Errorf("Relation not correctly unmarshaled. Got: %s, Expected: %s", link.Relation, "p-p")
	}
	if link.Point != "38.870,-77.055" {
		t.Errorf("Point not correctly unmarshaled. Got: %s, Expected: %s", link.Point, "38.870,-77.055")
	}
	if link.URL != "https://example.com/image.jpg" {
		t.Errorf("URL not correctly unmarshaled. Got: %s, Expected: %s", link.URL, "https://example.com/image.jpg")
	}
	if link.Remarks != "Team Alpha" {
		t.Errorf("Remarks not correctly unmarshaled. Got: %s, Expected: %s", link.Remarks, "Team Alpha")
	}
	if link.Production.String() != "2023-05-15T10:30:00.000Z" {
		t.Errorf("Production not correctly unmarshaled. Got: %s, Expected: %s", link.Production, "2023-05-15T10:30:00.000Z")
	}
	if link.Version != "2.0" {
		t.Errorf("Version not correctly unmarshaled. Got: %s, Expected: %s", link.Version, "2.0")
	}
	if link.Parent != "Mission-123" {
		t.Errorf("Parent not correctly unmarshaled. Got: %s, Expected: %s", link.Parent, "Mission-123")
	}
	if link.Medium != "chat" {
		t.Errorf("Medium not correctly unmarshaled. Got: %s, Expected: %s", link.Medium, "chat")
	}
}

func TestLinkUnmarshalXMLWithStyle(t *testing.T) {
	// Given
	xmlData := `<link uid="ANDROID-355161070128345" type="a-f-G-U-C"><Style><LineStyle><color>ff0000ff</color><width>2</width></LineStyle></Style></link>`

	// When
	var link Link
	err := xml.Unmarshal([]byte(xmlData), &link)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to Link with Style: %v", err)
	}

	// Then
	if link.UID != "ANDROID-355161070128345" {
		t.Errorf("UID not correctly unmarshaled. Got: %s, Expected: %s", link.UID, "ANDROID-355161070128345")
	}
	if link.Type != "a-f-G-U-C" {
		t.Errorf("Type not correctly unmarshaled. Got: %s, Expected: %s", link.Type, "a-f-G-U-C")
	}
	if link.Style == nil {
		t.Fatalf("Style not correctly unmarshaled, got nil")
	}
	if link.Style.LineStyle == nil {
		t.Fatalf("LineStyle not correctly unmarshaled, got nil")
	}
	if link.Style.LineStyle.Color != "ff0000ff" {
		t.Errorf("Style.LineStyle.Color not correctly unmarshaled. Got: %s, Expected: %s", link.Style.LineStyle.Color, "ff0000ff")
	}
	if link.Style.LineStyle.Width != 2.0 {
		t.Errorf("Style.LineStyle.Width not correctly unmarshaled. Got: %f, Expected: %f", link.Style.LineStyle.Width, 2.0)
	}
}

func TestLinkSetters(t *testing.T) {
	// Given
	link := Link{}
	uid := "ANDROID-355161070128345"
	linkType := "a-f-G-U-C"
	relation := "p-p"
	point := "38.870,-77.055"
	url := "https://example.com/image.jpg"
	remarks := "Team Alpha"
	//production := "2023-05-15T10:30:00Z"
	production := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	version := "2.0"
	parent := "Mission-123"
	medium := "chat"
	style := &Style{
		LineStyle: &LineStyle{
			Color: "ff0000ff",
			Width: 2.0,
		},
	}

	// When
	result := link.
		SetUID(uid).
		SetType(linkType).
		SetRelation(relation).
		SetPoint(point).
		SetURL(url).
		SetRemarks(remarks).
		SetProduction(production).
		SetVersion(version).
		SetParent(parent).
		SetMedium(medium).
		SetStyle(style)

	// Then
	if link.UID != uid {
		t.Errorf("UID not set correctly. Got: %s, Expected: %s", link.UID, uid)
	}
	if link.Type != linkType {
		t.Errorf("Type not set correctly. Got: %s, Expected: %s", link.Type, linkType)
	}
	if link.Relation != relation {
		t.Errorf("Relation not set correctly. Got: %s, Expected: %s", link.Relation, relation)
	}
	if link.Point != point {
		t.Errorf("Point not set correctly. Got: %s, Expected: %s", link.Point, point)
	}
	if link.URL != url {
		t.Errorf("URL not set correctly. Got: %s, Expected: %s", link.URL, url)
	}
	if link.Remarks != remarks {
		t.Errorf("Remarks not set correctly. Got: %s, Expected: %s", link.Remarks, remarks)
	}
	if link.Production.Time() != production {
		t.Errorf("Production not set correctly. Got: %s, Expected: %s", link.Production, production)
	}
	if link.Version != version {
		t.Errorf("Version not set correctly. Got: %s, Expected: %s", link.Version, version)
	}
	if link.Parent != parent {
		t.Errorf("Parent not set correctly. Got: %s, Expected: %s", link.Parent, parent)
	}
	if link.Medium != medium {
		t.Errorf("Medium not set correctly. Got: %s, Expected: %s", link.Medium, medium)
	}
	if link.Style != style {
		t.Errorf("Style not set correctly. Expected the same pointer")
	}
	if result != &link {
		t.Errorf("Setter methods did not return the link instance for chaining")
	}
}
