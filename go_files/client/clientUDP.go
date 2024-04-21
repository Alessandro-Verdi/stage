package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("udp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Send request
	_, err = conn.Write([]byte("/page"))
	if err != nil {
		fmt.Println("Error sending:", err)
		return
	}

	// Receive response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error receiving:", err)
		return
	}

	// Print response
	fmt.Println(string(buffer[:n]))
}

