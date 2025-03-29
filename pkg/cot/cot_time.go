package cot

import (
	"encoding/xml"
	"fmt"
	"time"
)

// Time format constants for CoT
const (
	// CotTimeFormat is the standard CoT time format with fractional seconds
	CotTimeFormat = "2006-01-02T15:04:05.000Z"
	// CotTimeFormatNoFraction is the CoT time format without fractional seconds
	CotTimeFormatNoFraction = "2006-01-02T15:04:05Z"
)

// CotTime is a custom time type that marshals/unmarshals in CoT format
type CotTime time.Time

// NewCotTime creates a new CotTime from a time.Time
func NewCotTime(t time.Time) CotTime {
	return CotTime(t)
}

// Now returns the current time as CotTime
func Now() CotTime {
	return CotTime(time.Now().UTC())
}

// NowWithLocation returns the current time as CotTime in the specified location
func NowWithLocation(loc *time.Location) CotTime {
	return CotTime(time.Now().In(loc))
}

// MarshalXML implements xml.Marshaler for CotTime
func (t CotTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	timeStr := FormatCotTime(time.Time(t))
	return e.EncodeElement(timeStr, start)
}

// UnmarshalXML implements xml.Unmarshaler for CotTime
func (t *CotTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var timeStr string
	if err := d.DecodeElement(&timeStr, &start); err != nil {
		return err
	}

	parsedTime, err := time.Parse(CotTimeFormat, timeStr)
	if err != nil {
		// Try without fractional seconds
		parsedTime, err = time.Parse(CotTimeFormatNoFraction, timeStr)
		if err != nil {
			return err
		}
	}

	*t = CotTime(parsedTime)
	return nil
}

// MarshalXMLAttr implements xml.MarshalerAttr for CotTime
func (t CotTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: FormatCotTime(time.Time(t)),
	}, nil
}

// UnmarshalXMLAttr implements xml.UnmarshalerAttr for CotTime
func (t *CotTime) UnmarshalXMLAttr(attr xml.Attr) error {
	parsedTime, err := time.Parse(CotTimeFormat, attr.Value)
	if err != nil {
		// Try without fractional seconds
		parsedTime, err = time.Parse(CotTimeFormatNoFraction, attr.Value)
		if err != nil {
			return err
		}
	}

	*t = CotTime(parsedTime)
	return nil
}

// MarshalJSON implements json.Marshaler for CotTime
func (t CotTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", FormatCotTime(time.Time(t)))), nil
}

// UnmarshalJSON implements json.Unmarshaler for CotTime
func (t *CotTime) UnmarshalJSON(data []byte) error {
	// Remove quotes
	timeStr := string(data)
	if len(timeStr) >= 2 && timeStr[0] == '"' && timeStr[len(timeStr)-1] == '"' {
		timeStr = timeStr[1 : len(timeStr)-1]
	}

	parsedTime, err := time.Parse(CotTimeFormat, timeStr)
	if err != nil {
		// Try without fractional seconds
		parsedTime, err = time.Parse(CotTimeFormatNoFraction, timeStr)
		if err != nil {
			return err
		}
	}

	*t = CotTime(parsedTime)
	return nil
}

// Time returns the time.Time representation of CotTime
func (t CotTime) Time() time.Time {
	return time.Time(t)
}

// String returns the string representation of CotTime
func (t CotTime) String() string {
	return FormatCotTime(time.Time(t))
}

// Add returns the time t+d
func (t CotTime) Add(d time.Duration) CotTime {
	return CotTime(time.Time(t).Add(d))
}

// FormatCotTime formats time in CoT-compatible format (ISO 8601 with optional fractional seconds)
// Examples: 2002-10-05T18:00:23Z, 2002-10-05T18:00:23.12Z, 2002-10-05T18:00:23.123456Z
func FormatCotTime(t time.Time) string {
	return t.Format(CotTimeFormat)
}

// FormatCotTimeWithPrecision formats time with specified decimal precision for seconds
// precision: number of decimal places for seconds (0-6)
func FormatCotTimeWithPrecision(t time.Time, precision int) string {
	if precision <= 0 {
		return t.Format(CotTimeFormatNoFraction)
	}

	if precision > 6 {
		precision = 6 // Cap at microsecond precision
	}

	format := fmt.Sprintf("2006-01-02T15:04:05.%sZ", "999999"[:precision])
	return t.Format(format)
}
