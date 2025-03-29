package cot

import (
	"encoding/xml"
	"time"
)

// Link represents a link element as defined in link.xsd
type Link struct {
	XMLName    xml.Name `xml:"link" json:"-"`
	UID        string   `xml:"uid,attr,omitempty" json:"uid,omitempty"`
	Type       string   `xml:"type,attr,omitempty" json:"type,omitempty"`
	Relation   string   `xml:"relation,attr,omitempty" json:"relation,omitempty"`
	Point      string   `xml:"point,attr,omitempty" json:"point,omitempty"`
	URL        string   `xml:"url,attr,omitempty" json:"url,omitempty"`
	Remarks    string   `xml:"remarks,attr,omitempty" json:"remarks,omitempty"`
	Production *CotTime `xml:"production_time,attr,omitempty" json:"production_time,omitempty"`
	Version    string   `xml:"version,attr,omitempty" json:"version,omitempty"`
	Parent     string   `xml:"parent,attr,omitempty" json:"parent,omitempty"`
	Medium     string   `xml:"medium,attr,omitempty" json:"medium,omitempty"`
	Style      *Style   `xml:"Style,omitempty" json:"style,omitempty"`
}

// SetUID sets the uid attribute of the link
func (l *Link) SetUID(uid string) *Link {
	l.UID = uid
	return l
}

// SetType sets the type attribute of the link
func (l *Link) SetType(linkType string) *Link {
	l.Type = linkType
	return l
}

// SetRelation sets the relation attribute of the link
func (l *Link) SetRelation(relation string) *Link {
	l.Relation = relation
	return l
}

// SetPoint sets the point attribute of the link
func (l *Link) SetPoint(point string) *Link {
	l.Point = point
	return l
}

// SetURL sets the url attribute of the link
func (l *Link) SetURL(url string) *Link {
	l.URL = url
	return l
}

// SetRemarks sets the remarks attribute of the link
func (l *Link) SetRemarks(remarks string) *Link {
	l.Remarks = remarks
	return l
}

// SetProduction sets the production_time attribute of the link
func (l *Link) SetProduction(production time.Time) *Link {
	cotTime := CotTime(production)
	l.Production = &cotTime
	return l
}

// SetVersion sets the version attribute of the link
func (l *Link) SetVersion(version string) *Link {
	l.Version = version
	return l
}

// SetParent sets the parent attribute of the link
func (l *Link) SetParent(parent string) *Link {
	l.Parent = parent
	return l
}

// SetMedium sets the medium attribute of the link
func (l *Link) SetMedium(medium string) *Link {
	l.Medium = medium
	return l
}

// SetStyle sets the Style element of the link
func (l *Link) SetStyle(style *Style) *Link {
	l.Style = style
	return l
}
