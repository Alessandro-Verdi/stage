
package main

import (

	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// Create UDP connection
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		log.Fatal("Error resolving UDP address:", err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal("Error establishing UDP connection:", err)
	}
	defer conn.Close()

	// Make a UDP request to the server
	addr := os.Args[1]
	log.Printf("Sending UDP request to %s", addr)

	startTime := time.Now()

	_, err = conn.Write([]byte("GET /page"))
	if err != nil {
		log.Fatal("Error sending UDP request:", err)
	}

	// Read responses from UDP connection
	var responses [][]byte
	timeout := time.After(5 * time.Second) // Change timeout duration as needed
	expectedResponses := 76000
	r := 0               // Number of expected responses
	for {
		select {
		case <-timeout:
			log.Println("Timeout reached. Exiting...")
			break
		default:
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			log.Println("received")
			if err != nil {
				if err == io.EOF {
					break // No more data available
				}
				log.Fatal("Error reading UDP response:", err)
			}
			r = r + 1
			responses = append(responses, buf[:n])
			}
			if r == expectedResponses {
				log.Println("stop")
				break
			
		}
	}

	// Calculate duration
	duration := time.Since(startTime)

	// Log received responses
	log.Printf("Received %d UDP responses from %s", len(responses), addr)
	//for i, response := range responses {
		//log.Printf("Response %d: %s", i+1, response)
	//}
	log.Printf("Round-trip time: %s", duration)
}

