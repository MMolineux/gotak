package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/angry-kivi/gotak/pkg/cot"
	"github.com/angry-kivi/gotak/pkg/parser"
	"github.com/angry-kivi/gotak/pkg/tak"
	"github.com/sirupsen/logrus"
)

func spotEvent(uid string) *cot.Event {
	spotEvent := cot.NewEvent("a-f-S-C", "spot-1"+uid)
	spotEvent.Point.Lat = 38.85606343062312
	spotEvent.Point.Lon = -77.0563755018233
	spotEvent.SetType("b-m-p-s-m")

	// Set the Europe/Stockholm timezone
	loc, _ := time.LoadLocation("Europe/Stockholm")
	now := cot.CotTime(time.Now().In(loc))
	staleTime := now.Add(10 * time.Minute)

	point := cot.NewPoint(38.85606343062312, -77.0563755018233)

	detail := cot.Detail{
		Status:  &cot.Status{Readiness: true},
		Archive: &cot.Archive{},
		Contact: &cot.Contact{Callsign: "R1"},
	}
	detail.AddLink(
		&cot.Link{
			UID:        "spot-1.Link",
			Type:       "a-f-G-U-C",
			Parent:     "spot-1",
			Relation:   "p-p",
			Production: &now,
		},
	)
	return &cot.Event{
		Version: "2.0",
		UID:     fmt.Sprintf("spot-1-%s", uid),
		Type:    "b-m-p-s-m",
		Time:    now,
		Start:   now,
		Stale:   staleTime,
		How:     "h-g-i-g-o",
		Point:   point,
		Detail:  detail,
	}
}

func createRect(uid string) *cot.Event {

	// Set the Europe/Stockholm timezone
	loc, _ := time.LoadLocation("Europe/Stockholm")
	now := cot.CotTime(time.Now().In(loc))
	staleTime := now.Add(10 * time.Minute)

	point := cot.NewPoint(38.83771744357615, 77.06824708113128)

	detail := cot.Detail{
		Archive: &cot.Archive{},
		Contact: &cot.Contact{Callsign: "Rectangle"},
	}
	detail.AddPointLink("38.83884480020009,-77.06896916307281")
	detail.AddPointLink("38.83878017039543,-77.06737849735573")
	detail.AddPointLink("38.8365895820601,-77.0675250569016")
	detail.AddPointLink("38.83665521881576,-77.06911560643489")
	detail.SetLabelsOn(true)
	detail.SetTog(false)
	detail.SetPrecisionLocation("???")
	detail.SetStrokeColor(-1)
	detail.SetStrokeWeight(3.0)
	detail.SetFillColor(-1761607681)

	return &cot.Event{
		Version: "2.0",
		UID:     fmt.Sprintf("rect-1-%s", uid),
		Type:    "u-d-r",
		Time:    now,
		Start:   now,
		Stale:   staleTime,
		How:     "h-e",
		Point:   point,
		Detail:  detail,
	}
}

func main() {
	// Define command-line flags
	serverAddr := flag.String("server", "takserver", "TAK server address")
	serverPort := flag.Int("port", tak.DefaultTCPPort, "TAK server port")
	connectionType := flag.String("connection", "tcp", "Connection type (tcp, tls, udp, multicast)")
	multicastAddr := flag.String("multicast-addr", tak.DefaultMulticastAddr, "Multicast group address")
	multicastPort := flag.Int("multicast-port", tak.DefaultMulticastPort, "Multicast port")
	clientID := flag.String("id", "gotak-client", "Client identifier")
	logLevel := flag.String("log-level", "info", "Log level (trace, debug, info, warn, error, fatal, panic)")

	// TLS certificate flags
	certFile := flag.String("cert", "", "Client certificate file (.pem or .p12)")
	keyFile := flag.String("key", "", "Client private key file (.key) - not needed for .p12")
	p12Password := flag.String("password", "", "Password for .p12 certificate file")
	caFile := flag.String("ca", "", "CA certificate file (.pem)")
	skipVerify := flag.Bool("skip-verify", false, "Skip TLS certificate verification")

	flag.Parse()

	// Configure logger
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set log level from command line flag
	level, err := logrus.ParseLevel(strings.ToLower(*logLevel))
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	log.SetLevel(level)
	log.Debugf("Log level set to: %s", level)

	// Determine connection type from the flag
	var connType tak.ConnectionType
	switch strings.ToLower(*connectionType) {
	case "tcp":
		connType = tak.ConnectionTypeTCP
	case "tls":
		connType = tak.ConnectionTypeTLS
	case "udp":
		connType = tak.ConnectionTypeUDP
	case "multicast":
		connType = tak.ConnectionTypeMulticast
	default:
		log.Fatalf("Invalid connection type: %s. Valid options are: tcp, tls, udp, multicast", *connectionType)
	}

	// Configure client
	config := tak.ClientConfig{
		Address:        *serverAddr,
		Port:           *serverPort,
		ClientID:       *clientID,
		ConnectionType: connType,
		DialTimeout:    10 * time.Second,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   10 * time.Second,

		// TLS certificate configuration
		CertFile:      *certFile,
		KeyFile:       *keyFile,
		P12Password:   *p12Password,
		CAFile:        *caFile,
		SkipTLSVerify: *skipVerify,

		// Multicast configuration
		MulticastAddr: *multicastAddr,
		MulticastPort: *multicastPort,

		// Add logger to client config
		Logger: log,
	}

	client, err := tak.NewClient(config)
	if err != nil {
		log.WithError(err).Fatal(fmt.Sprintf("Failed to create %s client", connType))
	}

	switch connType {
	case tak.ConnectionTypeTLS:
		log.Info("Using TLS connection")
	case tak.ConnectionTypeUDP:
		log.Info("Using UDP connection")
	case tak.ConnectionTypeMulticast:
		log.WithFields(logrus.Fields{
			"address": *multicastAddr,
			"port":    *multicastPort,
		}).Info("Using multicast UDP connection")
	default: // TCP
		log.Info("Using TCP connection")
	}

	// Connect to server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.WithFields(logrus.Fields{
		"server": *serverAddr,
		"port":   *serverPort,
	}).Info("Connecting to TAK server...")

	if err := client.Connect(ctx); err != nil {
		log.WithError(err).Fatal("Failed to connect")
	}
	defer func() {
		if err := client.Disconnect(); err != nil {
			log.WithError(err).Error("Error during disconnect")
		}
	}()

	log.Info("Connected successfully")

	// Create and send a simple ping event
	// pingEvent := cot.NewPingEvent(*clientID)

	posEvent := cot.NewEvent("a-f-S-C", "pos-1-"+*connectionType)
	posEvent.Point.Lat = 41.04911339629699
	posEvent.Point.Lon = -102.97315828649731
	posEvent.Detail = cot.Detail{
		Contact: &cot.Contact{
			Callsign:     "GOTAK",
			Phone:        "555-555-5555",
			EmailAddress: "test@gotak.com",
		},
	}

	xmlParser := parser.NewXMLParser()
	posData, err := xmlParser.SerializePrettyCoT(posEvent)
	if err != nil {
		log.WithError(err).Fatal("Failed to serialize CoT event")
	}

	// Print the cot event
	log.Info("Sending CoT event:")
	fmt.Println(string(posData))

	log.Debug("Sending CoT event...")
	if err := client.Send(posData); err != nil {
		log.WithError(err).Fatal("Failed to send data")
	}
	log.Info("CoT event sent successfully")
	// Wait for a response
	if connType == tak.ConnectionTypeUDP {
		log.Info("Waiting for response (UDP mode may not receive a response)...")
	} else {
		log.Info("Waiting for response...")
	}

	// Set a shorter timeout for UDP
	responseTimeout := 30 * time.Second
	if connType == tak.ConnectionTypeUDP {
		responseTimeout = 5 * time.Second
	}

	// If using multicast, update the receive timeout
	if connType == tak.ConnectionTypeMulticast {
		responseTimeout = 60 * time.Second
		log.Info("Multicast mode: listening for 60 seconds...")
	}

	// Create context with timeout for receive operation
	receiveCtx, receiveCancel := context.WithTimeout(context.Background(), responseTimeout)
	defer receiveCancel()

	// Use a channel to handle the receive operation
	respChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		resp, err := client.Receive()
		if err != nil {
			errChan <- err
			return
		}
		respChan <- resp
	}()

	// Wait for response or timeout
	select {
	case response := <-respChan:
		// Try to parse the response
		responseEvent, err := xmlParser.ParseCoT(response)
		if err != nil {
			log.WithField("raw_data", string(response)).Warn("Received unparseable data")
		} else {
			log.WithFields(logrus.Fields{
				"type": responseEvent.Type,
				"uid":  responseEvent.UID,
			}).Info("Received CoT event")

			if log.IsLevelEnabled(logrus.DebugLevel) {
				fmt.Printf("Full CoT event details: %+v\n", responseEvent)
			}
		}
	case err := <-errChan:
		log.WithError(err).Error("Error receiving data")
	case <-receiveCtx.Done():
		if connType == tak.ConnectionTypeUDP {
			log.Info("No response received within timeout (expected for UDP)")
		} else {
			log.Fatal("No response received within timeout")
		}
	}

	log.Info("Demonstration complete")
}
