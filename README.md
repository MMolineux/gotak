# GoTAK - Go Implementation of TAK Protocol

GoTAK is a Go library for working with the TAK (Tactical Assault Kit) protocol and CoT (Cursor on Target) data format. It provides a set of tools for connecting to TAK servers, as well as parsing and generating CoT messages.

## Features

- TAK Protocol Support: Connect with ATAK, WinTAK, iTAK, and TAK Server
- Data Handling: Manage TAK, Cursor on Target (CoT), and non-CoT data
- Data Parsing and Serialization: Parse and serialize TAK and CoT data
- Network Communication: Send and receive TAK and CoT data over TCP, TLS, and UDP
- TLS Security: Support for certificate-based authentication with .pem and .key files

## Getting Started

### Installation

```bash
go get github.com/angry-kivi/gotak
```

### Basic Usage

Connect to a TAK server and send a CoT message:

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/angry-kivi/gotak/pkg/cot"
	"github.com/angry-kivi/gotak/pkg/tak"
	"github.com/angry-kivi/gotak/pkg/parser"
)

func main() {
	// Configure TAK client
	config := tak.ClientConfig{
		Address:        "takserver.example.com",
		Port:           8087,
		ClientID:       "gotak-client",
		ConnectionType: tak.ConnectionTypeTCP, // Or ConnectionTypeUDP
		DialTimeout:    10 * time.Second,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	// Create and connect client
	client, err := tak.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Connect to server
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect()

	// Create a CoT event
	event := cot.NewEvent("a-f-G-U-C", "GOTAK-DEMO")
	event.Point.Lat = "38.8977"
	event.Point.Lon = "-77.0365"

	// Serialize the event to XML
	xmlParser := parser.NewXMLParser()
	data, err := xmlParser.SerializeCoT(event)
	if err != nil {
		log.Fatalf("Failed to serialize CoT event: %v", err)
	}

	// Send the event
	if err := client.Send(data); err != nil {
		log.Fatalf("Failed to send data: %v", err)
	}
}
```

### Using TLS with Certificates

For secure TLS connections with certificate authentication:

```go
config := tak.ClientConfig{
    Address:        "takserver.example.com",
    Port:           8089,  // TLS port
    ClientID:       "gotak-client",
    ConnectionType: tak.ConnectionTypeTLS,
    CertFile:       "/path/to/client.pem",
    KeyFile:        "/path/to/client.key",
    CAFile:         "/path/to/ca.pem",  // Optional
    SkipTLSVerify:  false,  // Set to true to skip server certificate verification
    ReadTimeout:    30 * time.Second,
    WriteTimeout:   10 * time.Second,
}

client, err := tak.NewTLSClient(config)
```

### Using UDP Mode

For connectionless UDP communication:

```go
config := tak.ClientConfig{
    Address:        "takserver.example.com",
    Port:           8087,
    ClientID:       "gotak-client",
    ConnectionType: tak.ConnectionTypeUDP,
    ReadTimeout:    5 * time.Second,  // Shorter timeout for UDP
    WriteTimeout:   5 * time.Second,
}

client, err := tak.NewUDPClient(config)
```

## Command Line Usage

The GoTAK command-line client supports various connection options:

```bash
# Connect via TCP (default)
./gotak --server takserver.example.com --port 8087

# Connect via TLS with certificates
./gotak --server takserver.example.com --port 8089 --tls --cert client.pem --key client.key --ca ca.pem

# Connect via UDP
./gotak --server takserver.example.com --port 8087 --udp
```

## Project Structure

- `pkg/tak` - Core TAK protocol implementation
- `pkg/cot` - CoT (Cursor on Target) data types and utilities
- `pkg/parser` - XML and other format parsers
- `pkg/util` - Utility functions and helpers
- `cmd/gotak` - Command-line client example

## Documentation

GoTAK includes comprehensive documentation for TAK Protocol and Cursor on Target (CoT):

- **CoT XML Schemas**: Full XSD documentation is available in the `/doc/html/` directory.
- **TAK Protocol**: Implementation details and protocol specifications are available in the source code.
- **API Documentation**: Run `godoc -http=:6060` and visit `http://localhost:6060/pkg/github.com/angry-kivi/gotak/` for API docs.

## Examples

### Receiving CoT Messages

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/angry-kivi/gotak/pkg/tak"
    "github.com/angry-kivi/gotak/pkg/parser"
)

func main() {
    config := tak.ClientConfig{
        Address:        "takserver.example.com",
        Port:           8087,
        ClientID:       "gotak-receiver",
        ConnectionType: tak.ConnectionTypeTCP,
        ReadTimeout:    60 * time.Second,
    }

    client, err := tak.NewClient(config)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    ctx := context.Background()
    if err := client.Connect(ctx); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer client.Disconnect()

    xmlParser := parser.NewXMLParser()
    
    // Receive and process CoT messages
    for {
        data, err := client.Receive()
        if err != nil {
            log.Printf("Error receiving data: %v", err)
            continue
        }

        event, err := xmlParser.ParseCoT(data)
        if err != nil {
            log.Printf("Error parsing CoT: %v", err)
            continue
        }

        fmt.Printf("Received event type: %s from: %s at location: %s, %s\n", 
            event.Type, event.UID, event.Point.Lat, event.Point.Lon)
    }
}
```

### Creating a Custom Marker with Details

```go
package main

import (
    "github.com/angry-kivi/gotak/pkg/cot"
    "github.com/angry-kivi/gotak/pkg/util"
)

func createCustomMarker() *cot.Event {
    // Create base event
    event := cot.NewEvent("a-f-G-U-C", "CUSTOM-MARKER")
    
    // Set coordinates
    event.Point.Lat = "38.8977"
    event.Point.Lon = "-77.0365"
    event.Point.Hae = "100"
    
    // Use the color converter to create color values
    colorConverter := util.NewColorConverter()
    
    // Convert from named color to CoT format (signed integer as string)
    redColor, _ := colorConverter.ConvertToCoTColor("red")
    // redColor will be "-16777216" for red
    
    // Add detail elements
    event.Detail = &cot.Detail{
        Contact: &cot.Contact{
            Callsign: "Observer1",
            Endpoint: "192.168.1.1:8087",
        },
        Remarks: &cot.Remarks{
            Value: "Important observation point",
            Source: "GoTAK",
        },
        Color: &cot.Color{
            Value: redColor, // Uses signed integer format expected by CoT
        },
    }
    
    return event
}
```

### Working with Colors

```go
package main

import (
    "fmt"
    "github.com/angry-kivi/gotak/pkg/util"
)

func main() {
    colorConverter := util.NewColorConverter()
    
    // Convert between color formats
    colorInt, _ := colorConverter.NameToInt("blue")
    fmt.Printf("Blue as unsigned integer: 0x%X\n", colorInt)
    
    // Get signed integer (CoT format)
    signedInt, _ := colorConverter.GetSignedInt("blue")
    fmt.Printf("Blue as signed integer for CoT: %d\n", signedInt)
    
    // Parse a CoT color (signed integer) back to a color
    cotColor := "-1761607681"  // A semi-transparent gray color in CoT format
    uint32Color, _ := colorConverter.ParseCoTColor(cotColor)
    r, g, b, a := colorConverter.IntToRGBA(uint32Color)
    fmt.Printf("CoT color %s is RGBA(%d, %d, %d, %d)\n", cotColor, r, g, b, a)
    
    // Parse various color formats
    colors := []string{
        "red",
        "#00FF00",
        "rgb(0, 0, 255)",
        "rgba(255, 255, 0, 128)",
    }
    
    for _, color := range colors {
        cotColor, _ := colorConverter.ConvertToCoTColor(color)
        fmt.Printf("Color %s converted to CoT format: %s\n", color, cotColor)
    }
}
```

## Compatibility

GoTAK has been tested with:

- ATAK (Android Team Awareness Kit)
- WinTAK (Windows Team Awareness Kit)
- iTAK (iOS Team Awareness Kit)
- TAK Server
- FreeTAKServer

## Development

### Prerequisites

- Go 1.16 or higher
- Basic understanding of TAK Protocol and CoT message format

### Building from Source

```bash
# Clone the repository
git clone https://github.com/angry-kivi/gotak.git
cd gotak

# Build the library
go build ./...

# Run tests
go test ./...
```

### Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
