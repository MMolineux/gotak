package cot

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

// Point represents the spatial information for a CoT event
type Point struct {
	XMLName xml.Name `xml:"point" json:"-"`

	Lat float64  `xml:"lat,attr" json:"lat"`
	Lon float64  `xml:"lon,attr" json:"lon"`
	Hae *float64 `xml:"hae,attr" json:"hae,omitempty"`
	Ce  *float64 `xml:"ce,attr" json:"ce,omitempty"`
	Le  *float64 `xml:"le,attr" json:"le,omitempty"`
}

func NewPoint(lat, lon float64) Point {
	return Point{
		Lat: lat,
		Lon: lon,
	}
}

// SetLat sets the latitude of the point
func (p *Point) SetLat(lat float64) *Point {
	p.Lat = lat
	return p
}

// SetLon sets the longitude of the point
func (p *Point) SetLon(lon float64) *Point {
	p.Lon = lon
	return p
}

// SetHae sets the height above ellipsoid of the point
func (p *Point) SetHae(hae float64) *Point {
	p.Hae = &hae
	return p
}

// SetCe sets the circular error of the point
func (p *Point) SetCe(ce float64) *Point {
	p.Ce = &ce
	return p
}

// SetLe sets the linear error of the point
func (p *Point) SetLe(le float64) *Point {
	p.Le = &le
	return p
}

func (p Point) String() string {
	return fmt.Sprintf("Lat: %f, Lon: %f, Hae: %f, Ce: %f, Le: %f", p.Lat, p.Lon, *p.Hae, *p.Ce, *p.Le)
}

// MarshalXML customizes XML marshaling for Point
func (p Point) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// Create our custom attributes
	attrs := []xml.Attr{
		{Name: xml.Name{Local: "lat"}, Value: strconv.FormatFloat(p.Lat, 'f', -1, 64)},
		{Name: xml.Name{Local: "lon"}, Value: strconv.FormatFloat(p.Lon, 'f', -1, 64)},
	}

	// Format Hae, Ce and Le as integers if they're DefaultValue, otherwise use regular float format

	var haeValue string
	if p.Hae == nil {
		haeValue = strconv.FormatFloat(DefaultValue, 'f', -1, 64)
	} else {
		haeValue = strconv.FormatFloat(*p.Hae, 'f', -1, 64)
	}
	attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "hae"}, Value: haeValue})

	var ceValue string
	if p.Ce == nil {
		ceValue = strconv.FormatFloat(DefaultValue, 'f', -1, 64)
	} else {
		ceValue = strconv.FormatFloat(*p.Ce, 'f', -1, 64)
	}
	attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "ce"}, Value: ceValue})

	var leValue string
	if p.Le == nil {
		leValue = strconv.FormatFloat(DefaultValue, 'f', -1, 64)
	} else {
		leValue = strconv.FormatFloat(*p.Le, 'f', -1, 64)
	}
	attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "le"}, Value: leValue})

	/*if p.Ce != 0 {
		var ceValue string
		if p.Ce == DefaultValue {
			ceValue = fmt.Sprintf("%d", int(p.Ce))
		} else {
			ceValue = strconv.FormatFloat(p.Ce, 'f', -1, 64)
		}
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "ce"}, Value: ceValue})
	}

	if p.Le != 0 {
		var leValue string
		if p.Le == DefaultValue {
			leValue = fmt.Sprintf("%d", int(p.Le))
		} else {
			leValue = strconv.FormatFloat(p.Le, 'f', -1, 64)
		}
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "le"}, Value: leValue})
	}*/

	// Create the element with our custom attributes
	start.Name = xml.Name{Local: "point"}
	start.Attr = attrs
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	// Close the element
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}
