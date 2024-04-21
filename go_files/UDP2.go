package main

import (
    "fmt"
    "net"
    "strings"
)

func handlePacket(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) {
    // Interpret the packet based on its content
    msg := strings.TrimSpace(string(buf))

    switch msg {
    case "GET /img":
        // Respond with image data (mocked)
        imgData := []byte{
            0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
            0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x28,
            0x01, 0x03, 0x00, 0x00, 0x00, 0xb6, 0x30, 0x2a, 0x2e, 0x00, 0x00, 0x00,
            0x03, 0x50, 0x4c, 0x54, 0x45, 0x5a, 0xc3, 0x5a, 0xad, 0x38, 0xaa, 0xdb,
            0x00, 0x00, 0x00, 0x0b, 0x49, 0x44, 0x41, 0x54, 0x78, 0x01, 0x63, 0x18,
            0x61, 0x00, 0x00, 0x00, 0xf0, 0x00, 0x01, 0xe2, 0xb8, 0x75, 0x22, 0x00,
            0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
        }
        _, _ = conn.WriteToUDP(imgData, addr)

    case "GET /page":
        // Respond with a long HTML page (mocked)
        htmlData := "<html><body>"
        for i := 0; i < 100; i++ {
            htmlData += fmt.Sprintf(`<p>                     %d </p>`, i)
        }
        //htmlData += "</body></html>"
	for i:= 0; i < 150000; i++ {
        _, _ = conn.WriteToUDP([]byte(htmlData), addr)
	}


    default:
        // Respond with a generic message
        _, _ = conn.WriteToUDP([]byte("Hello, UDP client!"), addr)
    }
}

func main() {
    // Resolve UDP address
    udpAddr, err := net.ResolveUDPAddr("udp", ":8080")
    if err != nil {
        fmt.Println("Error resolving UDP address:", err)
        return
    }

    // Create UDP listener
    udpConn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        fmt.Println("Error creating UDP listener:", err)
        return
    }
    defer udpConn.Close()

    fmt.Println("UDP server listening on port 8080")

    // Buffer for incoming data
    buf := make([]byte, 1024)

    // Infinite loop to handle incoming UDP packets
    for {
        // Read data from UDP connection
        n, addr, err := udpConn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println("Error reading from UDP:", err)
            continue
        }

        // Handle the received packet
        go handlePacket(udpConn, addr, buf[:n])
    }
}

