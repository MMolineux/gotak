package cot

import "encoding/xml"

// Shape represents the shape element as defined in shape.xsd
type Shape struct {
	XMLName  xml.Name  `xml:"shape" json:"shape,omitempty"`
	Ellipse  *Ellipse  `xml:"ellipse,omitempty" json:"ellipse,omitempty"`
	Polyline *Polyline `xml:"polyline,omitempty" json:"polyline,omitempty"`
	Link     *Link     `xml:"link,omitempty" json:"link,omitempty"`
}

// AddEllipse adds an ellipse element to the shape
func (s *Shape) AddEllipse(major, minor, angle float64) *Ellipse {
	s.Ellipse = &Ellipse{
		Major: major,
		Minor: minor,
		Angle: angle,
	}
	return s.Ellipse
}

// SetEllipse sets the ellipse element of the shape
func (s *Shape) SetEllipse(ellipse *Ellipse) *Shape {
	s.Ellipse = ellipse
	return s
}

// AddPolyline adds a polyline element to the shape
func (s *Shape) AddPolyline(points string) *Polyline {
	s.Polyline = &Polyline{
		Points: points,
	}
	return s.Polyline
}

// SetPolyline sets the polyline element of the shape
func (s *Shape) SetPolyline(polyline *Polyline) *Shape {
	s.Polyline = polyline
	return s
}

// SetLink sets the link element of the shape
func (s *Shape) SetLink(link *Link) *Shape {
	s.Link = link
	return s
}

// Ellipse represents the ellipse element within a shape
type Ellipse struct {
	XMLName xml.Name `xml:"ellipse" json:"-"`
	Major   float64  `xml:"major,attr" json:"major"`
	Minor   float64  `xml:"minor,attr" json:"minor"`
	Angle   float64  `xml:"angle,attr" json:"angle"`
}

// SetMajor sets the major axis value of the ellipse
func (e *Ellipse) SetMajor(major float64) *Ellipse {
	e.Major = major
	return e
}

// SetMinor sets the minor axis value of the ellipse
func (e *Ellipse) SetMinor(minor float64) *Ellipse {
	e.Minor = minor
	return e
}

// SetAngle sets the angle value of the ellipse
func (e *Ellipse) SetAngle(angle float64) *Ellipse {
	e.Angle = angle
	return e
}

// Polyline represents the polyline element within a shape
type Polyline struct {
	XMLName xml.Name `xml:"polyline" json:"-"`
	Points  string   `xml:"points,attr" json:"points"`
}

// SetPoints sets the points value of the polyline
func (p *Polyline) SetPoints(points string) *Polyline {
	p.Points = points
	return p
}
