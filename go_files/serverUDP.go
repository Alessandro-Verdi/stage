package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Create a UDP connection on port 8080
	conn, err := net.ListenPacket("udp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Buffer to store incoming data
	buffer := make([]byte, 1024)

	for {
		// Read data from the connection
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		// Handle request
		go handleRequest(conn, addr, buffer[:n])
	}
}

func handleRequest(conn net.PacketConn, addr net.Addr, data []byte) {
	// Check request path
	if string(data) == "/page" {
		// Prepare response
		response := "<html><body>"
		for i := 0; i < 5000000; i++ {
			response += fmt.Sprintf(`<p>%d</p>`, i)
		}
		response += "</body></html>"

		// Send response to client
		_, err := conn.WriteTo([]byte(response), addr)
		if err != nil {
			fmt.Println("Error writing:", err)
		}
	}
}

