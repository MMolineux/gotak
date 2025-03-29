package cot

import (
	"encoding/xml"
	"testing"
)

func TestContact_SetMethods(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Contact
		validate func(*Contact) bool
	}{
		{
			name: "SetCallsign",
			setup: func() *Contact {
				return new(Contact).SetCallsign("TEST1")
			},
			validate: func(c *Contact) bool {
				return c.Callsign == "TEST1"
			},
		},
		{
			name: "SetEmailAddress",
			setup: func() *Contact {
				return new(Contact).SetEmailAddress("test@example.com")
			},
			validate: func(c *Contact) bool {
				return c.EmailAddress == "test@example.com"
			},
		},
		{
			name: "SetEndpoint",
			setup: func() *Contact {
				return new(Contact).SetEndpoint("192.168.1.1:8080")
			},
			validate: func(c *Contact) bool {
				return c.Endpoint == "192.168.1.1:8080"
			},
		},
		{
			name: "SetPhone",
			setup: func() *Contact {
				return new(Contact).SetPhone("123-456-7890")
			},
			validate: func(c *Contact) bool {
				return c.Phone == "123-456-7890"
			},
		},
		{
			name: "SetXmppUsername",
			setup: func() *Contact {
				return new(Contact).SetXmppUsername("user123")
			},
			validate: func(c *Contact) bool {
				return c.XmppUsername == "user123"
			},
		},
		{
			name: "ChainedSetters",
			setup: func() *Contact {
				return new(Contact).
					SetCallsign("TEST1").
					SetEmailAddress("test@example.com").
					SetEndpoint("192.168.1.1:8080").
					SetPhone("123-456-7890").
					SetXmppUsername("user123")
			},
			validate: func(c *Contact) bool {
				return c.Callsign == "TEST1" &&
					c.EmailAddress == "test@example.com" &&
					c.Endpoint == "192.168.1.1:8080" &&
					c.Phone == "123-456-7890" &&
					c.XmppUsername == "user123"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contact := tt.setup()
			if !tt.validate(contact) {
				t.Errorf("Contact.%s() did not set field correctly", tt.name)
			}
		})
	}
}

func TestContact_XMLMarshal(t *testing.T) {
	contact := new(Contact).
		SetCallsign("TEST1").
		SetEmailAddress("test@example.com").
		SetEndpoint("192.168.1.1:8080").
		SetPhone("123-456-7890").
		SetXmppUsername("user123")

	data, err := xml.Marshal(contact)
	if err != nil {
		t.Fatalf("Failed to marshal Contact: %v", err)
	}

	expected := `<contact callsign="TEST1" emailAddress="test@example.com" endpoint="192.168.1.1:8080" phone="123-456-7890" xmppUsername="user123"></contact>`
	if string(data) != expected {
		t.Errorf("Contact XML did not match expected output\nGot: %s\nExpected: %s", string(data), expected)
	}
}

func TestTakv_SetMethods(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Takv
		validate func(*Takv) bool
	}{
		{
			name: "SetPlatform",
			setup: func() *Takv {
				return new(Takv).SetPlatform("ANDROID")
			},
			validate: func(takv *Takv) bool {
				return takv.Platform == "ANDROID"
			},
		},
		{
			name: "SetVersion",
			setup: func() *Takv {
				return new(Takv).SetVersion("4.1.0.0")
			},
			validate: func(takv *Takv) bool {
				return takv.Version == "4.1.0.0"
			},
		},
		{
			name: "SetOS",
			setup: func() *Takv {
				return new(Takv).SetOS("29")
			},
			validate: func(takv *Takv) bool {
				return takv.OS == "29"
			},
		},
		{
			name: "SetDevice",
			setup: func() *Takv {
				return new(Takv).SetDevice("Pixel 4")
			},
			validate: func(takv *Takv) bool {
				return takv.Device == "Pixel 4"
			},
		},
		{
			name: "ChainedSetters",
			setup: func() *Takv {
				return new(Takv).
					SetPlatform("ANDROID").
					SetVersion("4.1.0.0").
					SetOS("29").
					SetDevice("Pixel 4")
			},
			validate: func(takv *Takv) bool {
				return takv.Platform == "ANDROID" &&
					takv.Version == "4.1.0.0" &&
					takv.OS == "29" &&
					takv.Device == "Pixel 4"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			takv := tt.setup()
			if !tt.validate(takv) {
				t.Errorf("Takv.%s() did not set field correctly", tt.name)
			}
		})
	}
}

func TestTakv_XMLMarshal(t *testing.T) {
	takv := new(Takv).
		SetPlatform("ANDROID").
		SetVersion("4.1.0.0").
		SetOS("29").
		SetDevice("Pixel 4")

	data, err := xml.Marshal(takv)
	if err != nil {
		t.Fatalf("Failed to marshal Takv: %v", err)
	}

	expected := `<takv platform="ANDROID" version="4.1.0.0" os="29" device="Pixel 4"></takv>`
	if string(data) != expected {
		t.Errorf("Takv XML did not match expected output\nGot: %s\nExpected: %s", string(data), expected)
	}
}
