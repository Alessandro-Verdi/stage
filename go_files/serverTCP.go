package main

import (
    "fmt"
    "net/http"
    "io"
)

func setupHandler(www string) http.Handler {
	mux := http.NewServeMux()

	if len(www) > 0 {
		mux.Handle("/", http.FileServer(http.Dir(www)))
	}

	mux.HandleFunc("/img", func(w http.ResponseWriter, r *http.Request) {
		// Small 40x40 png
		w.Write([]byte{
			0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
			0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x28,
			0x01, 0x03, 0x00, 0x00, 0x00, 0xb6, 0x30, 0x2a, 0x2e, 0x00, 0x00, 0x00,
			0x03, 0x50, 0x4c, 0x54, 0x45, 0x5a, 0xc3, 0x5a, 0xad, 0x38, 0xaa, 0xdb,
			0x00, 0x00, 0x00, 0x0b, 0x49, 0x44, 0x41, 0x54, 0x78, 0x01, 0x63, 0x18,
			0x61, 0x00, 0x00, 0x00, 0xf0, 0x00, 0x01, 0xe2, 0xb8, 0x75, 0x22, 0x00,
			0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
		})
	})

	mux.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
        	io.WriteString(w, "<html><body>")
		for i := 0; i < 5000000; i++ {
			fmt.Fprintf(w, `<p> %d </p>`, i)
		}
		io.WriteString(w, "</body></html>")
	})

	return mux
} 

func main() {
    // Register the index handler
    mux := http.NewServeMux()
    handler := setupHandler("/tmp/quic-data")
    mux.Handle("/", handler)
    
    // Create an HTTP server with the ServeMux
    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    // Start the server
    fmt.Println("Server listening on port 8080")
    if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil {
        fmt.Println("Error starting server:", err)
    }
}

