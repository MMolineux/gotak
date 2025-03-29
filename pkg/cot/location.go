package cot

import "encoding/xml"

// PrecisionLocation represents the precisionlocation element as defined in precisionlocation.xsd
type PrecisionLocation struct {
	XMLName           xml.Name `xml:"precisionlocation" json:"-"`
	GeoPointSrc       string   `xml:"geopointsrc,attr,omitempty" json:"geopointsrc,omitempty"`
	AltSrc            string   `xml:"altsrc,attr" json:"altsrc"`
	PreciseImageFile  string   `xml:"PRECISE_IMAGE_FILE,attr,omitempty" json:"precise_image_file,omitempty"`
	PreciseImageFileX float64  `xml:"PRECISE_IMAGE_FILE_X,attr,omitempty" json:"precise_image_file_x,omitempty"`
	PreciseImageFileY float64  `xml:"PRECISE_IMAGE_FILE_Y,attr,omitempty" json:"precise_image_file_y,omitempty"`
}

// SetGeoPointSrc sets the geopointsrc attribute of the precision location
func (p *PrecisionLocation) SetGeoPointSrc(src string) *PrecisionLocation {
	p.GeoPointSrc = src
	return p
}

// SetAltSrc sets the altsrc attribute of the precision location
func (p *PrecisionLocation) SetAltSrc(src string) *PrecisionLocation {
	p.AltSrc = src
	return p
}

// SetPreciseImageFile sets the PRECISE_IMAGE_FILE attribute of the precision location
func (p *PrecisionLocation) SetPreciseImageFile(file string) *PrecisionLocation {
	p.PreciseImageFile = file
	return p
}

// SetPreciseImageFileX sets the PRECISE_IMAGE_FILE_X attribute of the precision location
func (p *PrecisionLocation) SetPreciseImageFileX(x float64) *PrecisionLocation {
	p.PreciseImageFileX = x
	return p
}

// SetPreciseImageFileY sets the PRECISE_IMAGE_FILE_Y attribute of the precision location
func (p *PrecisionLocation) SetPreciseImageFileY(y float64) *PrecisionLocation {
	p.PreciseImageFileY = y
	return p
}
