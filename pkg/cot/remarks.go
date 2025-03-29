package cot

import (
	"encoding/xml"
	"time"
)

// Remarks represents a remarks element as defined in remarks.xsd
type Remarks struct {
	XMLName  xml.Name   `xml:"remarks" json:"-"`
	Source   string     `xml:"source,attr,omitempty" json:"source,omitempty"`
	SourceID string     `xml:"sourceID,attr,omitempty" json:"source_id,omitempty"`
	Time     *time.Time `xml:"time,attr,omitempty" json:"time,omitempty"`
	To       string     `xml:"to,attr,omitempty" json:"to,omitempty"`
	Text     string     `xml:",chardata" json:"text"` // For the mixed content
}

// SetSource sets the source attribute of the remarks
func (r *Remarks) SetSource(source string) *Remarks {
	r.Source = source
	return r
}

// SetSourceID sets the sourceID attribute of the remarks
func (r *Remarks) SetSourceID(sourceID string) *Remarks {
	r.SourceID = sourceID
	return r
}

// SetTime sets the time attribute of the remarks
func (r *Remarks) SetTime(t *time.Time) *Remarks {
	r.Time = t
	return r
}

// SetTo sets the to attribute of the remarks
func (r *Remarks) SetTo(to string) *Remarks {
	r.To = to
	return r
}

// SetText sets the text content of the remarks
func (r *Remarks) SetText(text string) *Remarks {
	r.Text = text
	return r
}
