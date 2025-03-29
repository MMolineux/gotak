package cot

import "encoding/xml"

// Contact represents a contact element as defined in contact.xsd
type Contact struct {
	XMLName      xml.Name `xml:"contact" json:"-"`
	Callsign     string   `xml:"callsign,attr" json:"callsign"`
	EmailAddress string   `xml:"emailAddress,attr,omitempty" json:"email_address,omitempty"`
	Endpoint     string   `xml:"endpoint,attr,omitempty" json:"endpoint,omitempty"`
	Phone        string   `xml:"phone,attr,omitempty" json:"phone,omitempty"`
	XmppUsername string   `xml:"xmppUsername,attr,omitempty" json:"xmpp_username,omitempty"`
}

// SetCallsign sets the callsign attribute of the contact
func (c *Contact) SetCallsign(callsign string) *Contact {
	c.Callsign = callsign
	return c
}

// SetEmailAddress sets the emailAddress attribute of the contact
func (c *Contact) SetEmailAddress(email string) *Contact {
	c.EmailAddress = email
	return c
}

// SetEndpoint sets the endpoint attribute of the contact
func (c *Contact) SetEndpoint(endpoint string) *Contact {
	c.Endpoint = endpoint
	return c
}

// SetPhone sets the phone attribute of the contact
func (c *Contact) SetPhone(phone string) *Contact {
	c.Phone = phone
	return c
}

// SetXmppUsername sets the xmppUsername attribute of the contact
func (c *Contact) SetXmppUsername(username string) *Contact {
	c.XmppUsername = username
	return c
}

// Takv represents a takv element with platform info
type Takv struct {
	XMLName  xml.Name `xml:"takv" json:"-"`
	Platform string   `xml:"platform,attr" json:"platform"`
	Version  string   `xml:"version,attr" json:"version"`
	OS       string   `xml:"os,attr,omitempty" json:"os,omitempty"`
	Device   string   `xml:"device,attr,omitempty" json:"device,omitempty"`
}

// SetPlatform sets the platform attribute of the takv
func (t *Takv) SetPlatform(platform string) *Takv {
	t.Platform = platform
	return t
}

// SetVersion sets the version attribute of the takv
func (t *Takv) SetVersion(version string) *Takv {
	t.Version = version
	return t
}

// SetOS sets the os attribute of the takv
func (t *Takv) SetOS(os string) *Takv {
	t.OS = os
	return t
}

// SetDevice sets the device attribute of the takv
func (t *Takv) SetDevice(device string) *Takv {
	t.Device = device
	return t
}
