package main

import (
    "fmt"
    "net"
    "net/http"
)

func main() {
    // Example server code
    // Create a listener
    listener, err := net.Listen("tcp", "localhost:8080")
    if err != nil {
        panic(err)
    }
    defer listener.Close()

    // HTTP server handler for "/page"
    http.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
        // Set Content-Type header to indicate HTML content
        w.Header().Set("Content-Type", "text/html")

        // Write HTML response
        fmt.Fprintf(w, "<html><body><h1>This is a simple HTML page served over HTTP</h1></body></html>")
    })

    // Start the HTTP server
    go func() {
        fmt.Println("HTTP server is listening on port 8081...")
        if err := http.ListenAndServe(":8081", nil); err != nil {
            panic(err)
        }
    }()

    for {
        // Accept incoming connection
        conn, err := listener.Accept()
        if err != nil {
            panic(err)
        }
        defer conn.Close()

        // Reduce the maximum segment size to 1000 bytes
        tcpConn := conn.(*net.TCPConn)
        err = tcpConn.SetWriteBuffer(1000)
        if err != nil {
            panic(err)
        }

        // Handle connection
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    // Handle connection logic here
}

