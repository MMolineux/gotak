package cot

import (
	"encoding/xml"
)

// Detail contains extended information for CoT events
type Detail struct {
	XMLName xml.Name `xml:"detail" json:"-"`

	// Individual detail elements
	Contact           *Contact           `xml:"contact,omitempty" json:"contact,omitempty"`
	Remarks           *Remarks           `xml:"remarks,omitempty" json:"remarks,omitempty"`
	Status            *Status            `xml:"status,omitempty" json:"status,omitempty"`
	Takv              *Takv              `xml:"takv,omitempty" json:"takv,omitempty"`
	Track             *Track             `xml:"track,omitempty" json:"track,omitempty"`
	PrecisionLocation *PrecisionLocation `xml:"precisionlocation,omitempty" json:"precisionlocation,omitempty"`
	Shape             *Shape             `xml:"shape,omitempty" json:"shape,omitempty"`

	// Links can appear multiple times, especially for polygon points
	Links []*Link `xml:"link,omitempty" json:"link,omitempty"`

	// Drawing style elements - these are direct children of detail, not shape
	Color        *Color        `xml:"color,omitempty" json:"color,omitempty"`
	StrokeColor  *StrokeColor  `xml:"strokeColor,omitempty" json:"stroke_color,omitempty"`
	StrokeWeight *StrokeWeight `xml:"strokeWeight,omitempty" json:"stroke_weight,omitempty"`
	FillColor    *FillColor    `xml:"fillColor,omitempty" json:"fill_color,omitempty"`
	LabelsOn     *LabelsOn     `xml:"labels_on,omitempty" json:"labels_on,omitempty"`
	Archive      *Archive      `xml:"archive,omitempty" json:"archive,omitempty"`
	Tog          *Tog          `xml:"tog,omitempty" json:"tog,omitempty"`

	// Flow tags for mesh networking
	FlowTags *FlowTags `xml:"_flow-tags_,omitempty" json:"flow_tags,omitempty"`
	// Raw XML for any unmapped elements
	RawXML []byte `xml:",innerxml" json:"-"`
}

func (d Detail) SetPrecisionLocation(s string) {
	d.PrecisionLocation = &PrecisionLocation{
		AltSrc: s,
	}
}

// AddContact adds a contact element to the detail
func (d *Detail) AddContact(callsign string) *Contact {
	d.Contact = &Contact{
		Callsign: callsign,
	}
	return d.Contact
}

// AddTakv adds a takv element to the detail
func (d *Detail) AddTakv(platform, version string) *Takv {
	d.Takv = &Takv{
		Platform: platform,
		Version:  version,
	}
	return d.Takv
}

// AddRemarks adds a remarks element to the detail
func (d *Detail) AddRemarks(text string) *Remarks {
	d.Remarks = &Remarks{
		Text: text,
	}
	return d.Remarks
}

// AddStatus adds a status element to the detail
func (d *Detail) AddStatus() *Status {
	d.Status = &Status{}
	return d.Status
}

// AddTrack adds a track element to the detail
func (d *Detail) AddTrack() *Track {
	d.Track = &Track{}
	return d.Track
}

// AddPrecisionLocation adds a precision location element to the detail
func (d *Detail) AddPrecisionLocation(altSrc string) *PrecisionLocation {
	d.PrecisionLocation = &PrecisionLocation{
		AltSrc: altSrc,
	}
	return d.PrecisionLocation
}

// AddShape adds a shape element to the detail
func (d *Detail) AddShape() *Shape {
	d.Shape = &Shape{}
	return d.Shape
}

// AddLink adds a new link element to the detail and returns it
func (d *Detail) AddLink(link *Link) *Link {
	d.Links = append(d.Links, link)
	return link
}

// AddPointLink adds a link element with a point attribute
func (d *Detail) AddPointLink(point string) *Link {
	link := &Link{
		Point: point,
	}
	d.Links = append(d.Links, link)
	return link
}

// SetColor sets the color for the detail
func (d *Detail) SetColor(value int64) *Color {
	d.Color = &Color{
		Value: value,
	}
	return d.Color
}

// SetStrokeColor sets the stroke color for the detail
func (d *Detail) SetStrokeColor(value int64) *StrokeColor {
	d.StrokeColor = &StrokeColor{
		Value: value,
	}
	return d.StrokeColor
}

// SetStrokeWeight sets the stroke weight for the detail
func (d *Detail) SetStrokeWeight(value float64) *StrokeWeight {
	d.StrokeWeight = &StrokeWeight{
		Value: value,
	}
	return d.StrokeWeight
}

// SetFillColor sets the fill color for the detail
func (d *Detail) SetFillColor(value int64) *FillColor {
	d.FillColor = &FillColor{
		Value: value,
	}
	return d.FillColor
}

// SetLabelsOn sets whether labels should be displayed
func (d *Detail) SetLabelsOn(value bool) *LabelsOn {
	d.LabelsOn = &LabelsOn{
		Value: value,
	}
	return d.LabelsOn
}

// SetTog sets the tog value for the detail
func (d *Detail) SetTog(value bool) *Tog {
	d.Tog = &Tog{
		Value: value,
	}
	return d.Tog
}

// AddArchive adds an empty archive element to the detail
func (d *Detail) AddArchive() *Archive {
	d.Archive = &Archive{}
	return d.Archive
}

func (d *Detail) AddFlowTags(clientID string) *FlowTags {
	d.FlowTags = NewFlowTags(clientID)
	return d.FlowTags
}

// Archive represents an empty archive element
type Archive struct {
	XMLName xml.Name `xml:"archive" json:"-"`
}

type Tog struct {
	XMLName xml.Name `xml:"tog" json:"-"`
	Value   bool     `xml:"enabled,attr" json:"enabled"`
}

// UID contains user identification
type UID struct {
	XMLName xml.Name `xml:"__uid" json:"-"`

	Droid string `xml:"Droid,attr,omitempty" json:"Droid,omitempty"`
}

// Group contains group information
type Group struct {
	XMLName xml.Name `xml:"group" json:"-"`

	Name string `xml:"name,attr,omitempty" json:"name,omitempty"`
	Role string `xml:"role,attr,omitempty" json:"role,omitempty"`
}

// Chat contains chat information
type Chat struct {
	XMLName xml.Name `xml:"__chat" json:"-"`

	SenderCallsign string `xml:"senderCallsign,attr,omitempty" json:"senderCallsign,omitempty"`
	ChatRoom       string `xml:"chatroom,attr,omitempty" json:"chatroom,omitempty"`
	GroupOwner     string `xml:"groupOwner,attr,omitempty" json:"groupOwner,omitempty"`
	ID             string `xml:"id,attr,omitempty" json:"id,omitempty"`
	Parent         string `xml:"parent,attr,omitempty" json:"parent,omitempty"`
	MessageID      string `xml:"messageID,attr,omitempty" json:"messageID,omitempty"`
	ChatContent    string `xml:",chardata" json:"content,omitempty"`
}

// Sensor contains sensor information
type Sensor struct {
	XMLName xml.Name `xml:"sensor" json:"-"`

	Type          string `xml:"type,attr,omitempty" json:"type,omitempty"`
	Vfov          string `xml:"vfov,attr,omitempty" json:"vfov,omitempty"`
	Hfov          string `xml:"hfov,attr,omitempty" json:"hfov,omitempty"`
	Range         string `xml:"range,attr,omitempty" json:"range,omitempty"`
	Azimuth       string `xml:"azimuth,attr,omitempty" json:"azimuth,omitempty"`
	DisplayMagery string `xml:"displayMagery,attr,omitempty" json:"displayMagery,omitempty"`
}

// Video contains video information
type Video struct {
	XMLName xml.Name `xml:"video" json:"-"`

	Url             string `xml:"url,attr,omitempty" json:"url,omitempty"`
	NetworkTimeout  string `xml:"networkTimeout,attr,omitempty" json:"networkTimeout,omitempty"`
	BufferTime      string `xml:"buffer,attr,omitempty" json:"buffer,omitempty"`
	ConnectionToken string `xml:"connectionToken,attr,omitempty" json:"connectionToken,omitempty"`
}

// Height contains height information
type Height struct {
	XMLName xml.Name `xml:"height" json:"-"`

	Value     float64 `xml:"value,attr,omitempty" json:"value,omitempty"`
	Unit      string  `xml:"unit,attr,omitempty" json:"unit,omitempty"`
	Reference string  `xml:"reference,attr,omitempty" json:"reference,omitempty"`
}
