package cot

import (
	"encoding/xml"
	"time"
)

const DefaultValue = 9999999.0

// Event represents the main CoT event structure
type Event struct {
	XMLName xml.Name `xml:"event" json:"-"`

	// Event attributes
	Version string  `xml:"version,attr" json:"version"`
	UID     string  `xml:"uid,attr" json:"uid"`
	Type    string  `xml:"type,attr" json:"type"`
	Time    CotTime `xml:"time,attr" json:"time"`
	Start   CotTime `xml:"start,attr" json:"start"`
	Stale   CotTime `xml:"stale,attr" json:"stale"`
	How     string  `xml:"how,attr" json:"how"`

	Access string `xml:"access,attr,omitempty" json:"access,omitempty"`
	Qos    string `xml:"qos,attr,omitempty" json:"qos,omitempty"`
	Opex   string `xml:"opex,attr,omitempty" json:"opex,omitempty"`

	// Event elements
	Point  Point  `xml:"point" json:"point"`
	Detail Detail `xml:"detail" json:"detail"`
}

// SetVersion sets the version attribute of the event
func (e *Event) SetVersion(version string) *Event {
	e.Version = version
	return e
}

// SetUID sets the uid attribute of the event
func (e *Event) SetUID(uid string) *Event {
	e.UID = uid
	return e
}

// SetType sets the type attribute of the event
func (e *Event) SetType(eventType string) *Event {
	e.Type = eventType
	return e
}

// SetTime sets the time attribute of the event
func (e *Event) SetTime(t time.Time) *Event {
	e.Time = CotTime(t)
	return e
}

// SetStart sets the start attribute of the event
func (e *Event) SetStart(t time.Time) *Event {
	e.Start = CotTime(t)
	return e
}

// SetStale sets the stale attribute of the event
func (e *Event) SetStale(t time.Time) *Event {
	e.Stale = CotTime(t)
	return e
}

// SetHow sets the how attribute of the event
func (e *Event) SetHow(how string) *Event {
	e.How = how
	return e
}

// SetAccess sets the access attribute of the event
func (e *Event) SetAccess(access string) *Event {
	e.Access = access
	return e
}

// SetQos sets the qos attribute of the event
func (e *Event) SetQos(qos string) *Event {
	e.Qos = qos
	return e
}

// SetOpex sets the opex attribute of the event
func (e *Event) SetOpex(opex string) *Event {
	e.Opex = opex
	return e
}

// SetPoint sets the point element of the event
func (e *Event) SetPoint(point Point) *Event {
	e.Point = point
	return e
}

// SetDetail sets the detail element of the event
func (e *Event) SetDetail(detail Detail) *Event {
	e.Detail = detail
	return e
}

// NewEvent creates a new CoT event with default values
func NewEvent(eventType, uid string) *Event {
	// Set the Europe/Stockholm timezone
	loc, _ := time.LoadLocation("Europe/Stockholm")
	now := time.Now().In(loc)
	staleTime := now.Add(10 * time.Minute)

	return &Event{
		Version: "2.0",
		UID:     uid,
		Type:    eventType,
		Time:    CotTime(now),
		Start:   CotTime(now),
		Stale:   CotTime(staleTime),
		How:     "m-g",
		Point: Point{
			Lat: 0.0,
			Lon: 0.0,
		},
	}
}

func NewPingEvent(uid string) *Event {
	loc, _ := time.LoadLocation("Europe/Stockholm")
	now := time.Now().In(loc)
	staleTime := now.Add(10 * time.Minute)

	return &Event{
		Version: "2.0",
		UID:     uid,
		Type:    "t-x-c-t",
		Time:    CotTime(now),
		Start:   CotTime(now),
		Stale:   CotTime(staleTime),
		How:     "h-g-i-g-o",
		Point: Point{
			Lat: 0.0,
			Lon: 0.0,
		},
		Detail: Detail{
			Takv: &Takv{
				Platform: "GoTAK",
				Version:  "1.0",
			},
		},
	}
}
