// server.go
package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"
)

var (
    clients    = make(map[chan string]struct{})
    clientsMux sync.Mutex
)

func handleClient(w http.ResponseWriter, r *http.Request) {
    // Set response header to indicate SSE (Server-Sent Events)
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    // Create a new channel for this client
    clientChannel := make(chan string)

    // Register client channel
    clientsMux.Lock()
    clients[clientChannel] = struct{}{}
    clientsMux.Unlock()

    // Send initial message
    fmt.Fprintf(w, "data: %s\n\n", "Welcome to the server!")

    // Continuously send packets to the client
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            fmt.Fprintf(w, "data: %s\n", "Packet from server")
            w.(http.Flusher).Flush()
        case <-r.Context().Done():
            // Remove client channel when client disconnects
            clientsMux.Lock()
            delete(clients, clientChannel)
            clientsMux.Unlock()
            return
        }
    }
}

func main() {
    http.HandleFunc("/", handleClient)

    server := &http.Server{
        Addr: ":8080",
    }

    log.Println("Server started on :8080")
    log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}
