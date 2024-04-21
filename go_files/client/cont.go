// client.go
package main

import (
    "bufio"
    "crypto/tls"
    "fmt"
    "log"
    "net/http"
    "strings"
)

func main() {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}

    // Make GET request to server
    resp, err := client.Get("https://10.0.0.1:8080")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    // Read response body
    scanner := bufio.NewScanner(resp.Body)
    
    i := 0
    for scanner.Scan() {
        fmt.Printf("Packet %d received: %s\n", i, scanner.Text())

        // Check if stop command received
        if strings.Contains(scanner.Text(), "stop") {
            break
        }
     	i = i + 1
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}
