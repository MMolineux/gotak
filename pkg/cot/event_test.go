package cot

import (
	"encoding/json"
	"encoding/xml"
	"reflect"
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {
	// When
	event := NewEvent("a-f-G-U-C", "TEST-UID")

	// Then
	if event.Version != "2.0" {
		t.Errorf("Expected version to be 2.0, got %s", event.Version)
	}
	if event.UID != "TEST-UID" {
		t.Errorf("Expected UID to be TEST-UID, got %s", event.UID)
	}
	if event.Type != "a-f-G-U-C" {
		t.Errorf("Expected type to be a-f-G-U-C, got %s", event.Type)
	}
	if event.How != "m-g" {
		t.Errorf("Expected how to be m-g, got %s", event.How)
	}

	// Check point values
	if event.Point.Lat != 0.0 || event.Point.Lon != 0.0 {
		t.Errorf("Expected default point coordinates to be 0, got lat=%f, lon=%f",
			event.Point.Lat, event.Point.Lon)
	}
	if event.Point.Hae != nil || event.Point.Ce != nil || event.Point.Le != nil {
		t.Errorf("Expected default HAE/CE/LE values to be %f, got hae=%f, CE=%f, LE=%f",
			DefaultValue, *event.Point.Hae, *event.Point.Ce, *event.Point.Le)
	}

	// Time checks
	now := time.Now().In(time.UTC)
	eventTime := time.Time(event.Time).In(time.UTC)
	timeDelta := now.Sub(eventTime).Abs()

	if timeDelta > time.Second*5 {
		t.Errorf("Event time not set correctly. Current time: %v, Event time: %v, Difference: %v",
			now, eventTime, timeDelta)
	}

	// Check stale time is ~10 minutes in the future
	staleTime := time.Time(event.Stale).In(time.UTC)
	expectedStaleTime := eventTime.Add(10 * time.Minute)
	staleDelta := staleTime.Sub(expectedStaleTime).Abs()

	if staleDelta > time.Second*5 {
		t.Errorf("Stale time not set correctly. Expected: %v, Got: %v, Difference: %v",
			expectedStaleTime, staleTime, staleDelta)
	}
}

func TestNewPingEvent(t *testing.T) {
	// When
	event := NewPingEvent("PING-UID")

	// Then
	if event.Version != "2.0" {
		t.Errorf("Expected version to be 2.0, got %s", event.Version)
	}
	if event.UID != "PING-UID" {
		t.Errorf("Expected UID to be PING-UID, got %s", event.UID)
	}
	if event.Type != "t-x-c-t" {
		t.Errorf("Expected type to be t-x-c-t, got %s", event.Type)
	}
	if event.How != "h-g-i-g-o" {
		t.Errorf("Expected how to be h-g-i-g-o, got %s", event.How)
	}

	// Check that detail and takv are set
	if event.Detail.Takv == nil {
		t.Errorf("Expected Takv to be set in ping event")
	} else {
		if event.Detail.Takv.Platform != "GoTAK" {
			t.Errorf("Expected platform to be GoTAK, got %s", event.Detail.Takv.Platform)
		}
		if event.Detail.Takv.Version != "1.0" {
			t.Errorf("Expected version to be 1.0, got %s", event.Detail.Takv.Version)
		}
	}
}

func TestEventSetters(t *testing.T) {
	// Given
	event := &Event{}
	testTime := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	hae := 100.0
	ce := 50.0
	le := 25.0
	// Create a point with optional fields
	// Set Hae, Ce and Le to some values
	point := Point{Lat: 12.345, Lon: 67.890, Hae: &hae, Ce: &ce, Le: &le}
	detail := Detail{
		Takv: &Takv{
			Platform: "TestPlatform",
			Version:  "2.0",
		},
	}

	// When: Chain all setters
	event.SetVersion("2.1").
		SetUID("SETTER-TEST").
		SetType("a-f-G-U-C-I").
		SetTime(testTime).
		SetStart(testTime).
		SetStale(testTime.Add(time.Hour)).
		SetHow("m-g-test").
		SetAccess("RESTRICTED").
		SetQos("test-qos").
		SetOpex("test-opex").
		SetPoint(point).
		SetDetail(detail)

	// Then
	if event.Version != "2.1" {
		t.Errorf("SetVersion failed, got %s", event.Version)
	}
	if event.UID != "SETTER-TEST" {
		t.Errorf("SetUID failed, got %s", event.UID)
	}
	if event.Type != "a-f-G-U-C-I" {
		t.Errorf("SetType failed, got %s", event.Type)
	}
	if !time.Time(event.Time).Equal(testTime) {
		t.Errorf("SetTime failed, got %v", time.Time(event.Time))
	}
	if !time.Time(event.Start).Equal(testTime) {
		t.Errorf("SetStart failed, got %v", time.Time(event.Start))
	}
	if !time.Time(event.Stale).Equal(testTime.Add(time.Hour)) {
		t.Errorf("SetStale failed, got %v", time.Time(event.Stale))
	}
	if event.How != "m-g-test" {
		t.Errorf("SetHow failed, got %s", event.How)
	}
	if event.Access != "RESTRICTED" {
		t.Errorf("SetAccess failed, got %s", event.Access)
	}
	if event.Qos != "test-qos" {
		t.Errorf("SetQos failed, got %s", event.Qos)
	}
	if event.Opex != "test-opex" {
		t.Errorf("SetOpex failed, got %s", event.Opex)
	}
	if !reflect.DeepEqual(event.Point, point) {
		t.Errorf("SetPoint failed, got %+v", event.Point)
	}
	if !reflect.DeepEqual(event.Detail, detail) {
		t.Errorf("SetDetail failed, got %+v", event.Detail)
	}
}

func TestEventMarshalXML(t *testing.T) {
	// Given
	testTime := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	hae := 100.0
	ce := 50.0
	le := 25.0

	event := &Event{
		Version: "2.0",
		UID:     "XML-TEST-UID",
		Type:    "a-f-G-U-C",
		Time:    CotTime(testTime),
		Start:   CotTime(testTime),
		Stale:   CotTime(testTime.Add(time.Hour)),
		How:     "m-g",
		Point: Point{
			Lat: 12.345,
			Lon: 67.890,
			Hae: &hae,
			Ce:  &ce,
			Le:  &le,
		},
		Detail: Detail{
			Takv: &Takv{
				Platform: "TestPlatform",
				Version:  "2.0",
			},
		},
	}

	// When
	xmlData, err := xml.MarshalIndent(event, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal event to XML: %v", err)
	}

	// Then
	expectedXML := `<event version="2.0" uid="XML-TEST-UID" type="a-f-G-U-C" time="2023-05-15T10:30:00.000Z" start="2023-05-15T10:30:00.000Z" stale="2023-05-15T11:30:00.000Z" how="m-g">
  <point lat="12.345" lon="67.89" hae="100" ce="50" le="25"></point>
  <detail>
    <takv platform="TestPlatform" version="2.0"></takv>
  </detail>
</event>`

	if string(xmlData) != expectedXML {
		t.Errorf("Marshaled XML does not match expected.\nGot: %s\nExpected: %s", string(xmlData), expectedXML)
	}

	// When - unmarshal back to verify round trip
	var unmarshaledEvent Event
	err = xml.Unmarshal(xmlData, &unmarshaledEvent)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to event: %v", err)
	}

	// Then - check a few key fields
	if unmarshaledEvent.UID != "XML-TEST-UID" {
		t.Errorf("Unmarshaled UID incorrect. Got: %s, Expected: %s", unmarshaledEvent.UID, "XML-TEST-UID")
	}
	if unmarshaledEvent.Point.Lat != 12.345 {
		t.Errorf("Unmarshaled Point.Lat incorrect. Got: %f, Expected: %f", unmarshaledEvent.Point.Lat, 12.345)
	}
	if unmarshaledEvent.Point.Lon != 67.89 {
		t.Errorf("Unmarshaled Point.Lon incorrect. Got: %f, Expected: %f", unmarshaledEvent.Point.Lon, 67.89)
	}
	if unmarshaledEvent.Detail.Takv == nil || unmarshaledEvent.Detail.Takv.Platform != "TestPlatform" {
		t.Errorf("Unmarshaled Detail.Takv incorrect")
	}
}

func TestEventMarshalJSON(t *testing.T) {
	// Given
	testTime := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	hae := 100.0
	ce := 50.0
	le := 25.0
	event := &Event{
		Version: "2.0",
		UID:     "JSON-TEST-UID",
		Type:    "a-f-G-U-C",
		Time:    CotTime(testTime),
		Start:   CotTime(testTime),
		Stale:   CotTime(testTime.Add(time.Hour)),
		How:     "m-g",
		Point: Point{
			Lat: 12.345,
			Lon: 67.890,
			Hae: &hae,
			Ce:  &ce,
			Le:  &le,
		},
		Detail: Detail{
			Takv: &Takv{
				Platform: "TestPlatform",
				Version:  "2.0",
			},
		},
	}

	// When
	jsonData, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal event to JSON: %v", err)
	}

	// Then
	// Partial check - just verify key fields
	expectedJSON := `{
  "version": "2.0",
  "uid": "JSON-TEST-UID",
  "type": "a-f-G-U-C",
  "time": "2023-05-15T10:30:00.000Z",
  "start": "2023-05-15T10:30:00.000Z",
  "stale": "2023-05-15T11:30:00.000Z",
  "how": "m-g",
  "point": {
    "lat": 12.345,
    "lon": 67.89,
    "hae": 100,
    "ce": 50,
    "le": 25
  },
  "detail": {
    "takv": {
      "platform": "TestPlatform",
      "version": "2.0"
    }
  }
}`

	if string(jsonData) != expectedJSON {
		t.Errorf("Marshaled JSON does not match expected.\nGot: %s\nExpected: %s", string(jsonData), expectedJSON)
	}

	// When - unmarshal back to verify round trip
	var unmarshaledEvent Event
	err = json.Unmarshal(jsonData, &unmarshaledEvent)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to event: %v", err)
	}

	// Then - check a few key fields
	if unmarshaledEvent.UID != "JSON-TEST-UID" {
		t.Errorf("Unmarshaled UID incorrect. Got: %s, Expected: %s", unmarshaledEvent.UID, "JSON-TEST-UID")
	}
	if unmarshaledEvent.Point.Lat != 12.345 {
		t.Errorf("Unmarshaled Point.Lat incorrect. Got: %f, Expected: %f", unmarshaledEvent.Point.Lat, 12.345)
	}
	if unmarshaledEvent.Point.Lon != 67.89 {
		t.Errorf("Unmarshaled Point.Lon incorrect. Got: %f, Expected: %f", unmarshaledEvent.Point.Lon, 67.89)
	}
	if unmarshaledEvent.Detail.Takv == nil || unmarshaledEvent.Detail.Takv.Platform != "TestPlatform" {
		t.Errorf("Unmarshaled Detail.Takv incorrect")
	}
}

func TestEventWithOptionalFields(t *testing.T) {
	// Given
	testTime := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	hae := 100.0
	ce := 50.0
	le := 25.0
	event := &Event{
		Version: "2.0",
		UID:     "OPTIONAL-TEST-UID",
		Type:    "a-f-G-U-C",
		Time:    CotTime(testTime),
		Start:   CotTime(testTime),
		Stale:   CotTime(testTime.Add(time.Hour)),
		How:     "m-g",
		// Setting optional fields
		Access: "SECRET",
		Qos:    "high-quality",
		Opex:   "operational",
		Point: Point{
			Lat: 12.345,
			Lon: 67.890,
			Hae: &hae,
			Ce:  &ce,
			Le:  &le,
		},
	}

	// When
	xmlData, err := xml.MarshalIndent(event, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal event with optional fields to XML: %v", err)
	}

	// Then - verify optional fields are included
	expectedXML := `<event version="2.0" uid="OPTIONAL-TEST-UID" type="a-f-G-U-C" time="2023-05-15T10:30:00.000Z" start="2023-05-15T10:30:00.000Z" stale="2023-05-15T11:30:00.000Z" how="m-g" access="SECRET" qos="high-quality" opex="operational">
  <point lat="12.345" lon="67.89" hae="100" ce="50" le="25"></point>
  <detail></detail>
</event>`

	if string(xmlData) != expectedXML {
		t.Errorf("Marshaled XML with optional fields does not match expected.\nGot: %s\nExpected: %s", string(xmlData), expectedXML)
	}
}
