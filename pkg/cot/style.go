package cot

import "encoding/xml"

// Style represents a KML style element
type Style struct {
	XMLName   xml.Name   `xml:"Style" json:"-"`
	LineStyle *LineStyle `xml:"LineStyle,omitempty" json:"line_style,omitempty"`
	PolyStyle *PolyStyle `xml:"PolyStyle,omitempty" json:"poly_style,omitempty"`
}

// SetLineStyle sets the LineStyle element of the style
func (s *Style) SetLineStyle(lineStyle *LineStyle) *Style {
	s.LineStyle = lineStyle
	return s
}

// SetPolyStyle sets the PolyStyle element of the style
func (s *Style) SetPolyStyle(polyStyle *PolyStyle) *Style {
	s.PolyStyle = polyStyle
	return s
}

// LineStyle represents the KML LineStyle element
type LineStyle struct {
	XMLName xml.Name `xml:"LineStyle" json:"-"`
	Color   string   `xml:"color,omitempty" json:"color,omitempty"`
	Width   float64  `xml:"width,omitempty" json:"width,omitempty"`
}

// SetColor sets the color of the line style
func (ls *LineStyle) SetColor(color string) *LineStyle {
	ls.Color = color
	return ls
}

// SetWidth sets the width of the line style
func (ls *LineStyle) SetWidth(width float64) *LineStyle {
	ls.Width = width
	return ls
}

// PolyStyle represents the KML PolyStyle element
type PolyStyle struct {
	XMLName xml.Name `xml:"PolyStyle" json:"-"`
	Color   string   `xml:"color,omitempty" json:"color,omitempty"`
}

// SetColor sets the color of the poly style
func (ps *PolyStyle) SetColor(color string) *PolyStyle {
	ps.Color = color
	return ps
}

// Color represents a color element
type Color struct {
	XMLName xml.Name `xml:"color" json:"-"`
	Value   int64    `xml:"argb,attr" json:"value"`
}

// SetValue sets the argb value of the color
func (c *Color) SetValue(value int64) *Color {
	c.Value = value
	return c
}

// StrokeColor represents the strokeColor element
type StrokeColor struct {
	XMLName xml.Name `xml:"strokeColor" json:"-"`
	Value   int64    `xml:"value,attr" json:"value"`
}

// SetValue sets the value of the stroke color
func (sc *StrokeColor) SetValue(value int64) *StrokeColor {
	sc.Value = value
	return sc
}

// StrokeWeight represents the strokeWeight element
type StrokeWeight struct {
	XMLName xml.Name `xml:"strokeWeight" json:"-"`
	Value   float64  `xml:"value,attr" json:"value"`
}

// SetValue sets the value of the stroke weight
func (sw *StrokeWeight) SetValue(value float64) *StrokeWeight {
	sw.Value = value
	return sw
}

// FillColor represents the fillColor element
type FillColor struct {
	XMLName xml.Name `xml:"fillColor" json:"-"`
	Value   int64    `xml:"value,attr" json:"value"`
}

// SetValue sets the value of the fill color
func (fc *FillColor) SetValue(value int64) *FillColor {
	fc.Value = value
	return fc
}

// LabelsOn represents the labels_on element
type LabelsOn struct {
	XMLName xml.Name `xml:"labels_on" json:"-"`
	Value   bool     `xml:"value,attr" json:"value"`
}

// SetValue sets the value of the labels_on property
func (lo *LabelsOn) SetValue(value bool) *LabelsOn {
	lo.Value = value
	return lo
}
