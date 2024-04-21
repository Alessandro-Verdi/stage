package main

import (
    "log"
    "context"
    "time"
    "bytes"
    "io"
    "os"
    "net"
    "net/http"
    "crypto/tls"

)

func main() {

        httpClient := &http.Client{
                Transport: &http.Transport{
                        TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Skip SSL certificate verification (not recommended for production)
                        DialContext: (&net.Dialer{
                                Timeout:   30 * time.Second,
                                KeepAlive: 30 * time.Second,
                        }).DialContext,
                },
        }

    // Make a GET request to the server
    //addr := 'https://10.0.0.1:8080/'
    addr := os.Args[1]
    log.Printf("GET %s", addr)

    ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
    if err != nil {
        log.Fatal(err)
    }

    startTime := time.Now()

    resp, err := httpClient.Do(req)
    if err != nil {
        // Check for context cancellation error
        if err == context.DeadlineExceeded {
            log.Fatalf("request cancelled: %v", err)
        }
        log.Fatal(err)
    }
    defer resp.Body.Close()

    log.Printf("Got response for %s: %#v", addr, resp)

    body := &bytes.Buffer{}
    _, err = io.Copy(body, resp.Body)

    if err != nil {
        log.Fatal(err)
    }

    //finalTime := time.Now()
  //  diff := finalTime.Sub(startTime)
    duration := time.Since(startTime)
 //   throughput := float64(body.Len()) / duration.Seconds()

    log.Printf("Response Body (%d bytes):\n%s", body.Len(), body.Len())
    log.Printf("Throughput: %.2f bytes/second and duration: %d\n", duration.Seconds(), duration.Milliseconds())
}
