package cot

import (
	"encoding/xml"
	"time"
)

// Track represents a track element as defined in track.xsd
type Track struct {
	XMLName   xml.Name  `xml:"track" json:"track,omitempty"`
	Course    float64   `xml:"course,attr,omitempty" json:"course,omitempty"`
	Speed     float64   `xml:"speed,attr,omitempty" json:"speed,omitempty"`          // Speed in m/s
	Slope     float64   `xml:"slope,attr,omitempty" json:"slope,omitempty"`          // Vertical path angle in degrees
	Etype     string    `xml:"etype,attr,omitempty" json:"e_type,omitempty"`         // Using underscore for JSON
	TimeStamp time.Time `xml:"timeStamp,attr,omitempty" json:"time_stamp,omitempty"` // Using underscore for JSON
}

// SetCourse sets the course in degrees
func (t *Track) SetCourse(course float64) *Track {
	t.Course = course
	return t
}

// SetSpeed sets the speed in meters per second
func (t *Track) SetSpeed(speed float64) *Track {
	t.Speed = speed
	return t
}

// SetSlope sets the vertical path angle in degrees
func (t *Track) SetSlope(slope float64) *Track {
	t.Slope = slope
	return t
}

// SetEtype sets the entity type
func (t *Track) SetEtype(etype string) *Track {
	t.Etype = etype
	return t
}

// SetTimeStamp sets the timestamp
func (t *Track) SetTimeStamp(timestamp time.Time) *Track {
	t.TimeStamp = timestamp
	return t
}

// SetCurrentTimeStamp sets the timestamp to the current time
func (t *Track) SetCurrentTimeStamp() *Track {
	t.TimeStamp = time.Now()
	return t
}
