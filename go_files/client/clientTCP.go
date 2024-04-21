package main

import (
    "log"
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
    
    startTime := time.Now()

    resp, err := httpClient.Get(addr)
    if err != nil {
        log.Fatal(err)
    }

    
   
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

