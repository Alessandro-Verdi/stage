#!/bin/bash

# Generate a private key
openssl genrsa -out server.key 2048

# Generate a certificate signing request (CSR)
openssl req -new -key server.key -out server.csr -subj "/CN=localhost"

# Generate a self-signed SSL certificate valid for 365 days
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt

# Remove the CSR file (not needed anymore)
rm server.csr

echo "SSL certificate (server.crt) and private key (server.key) generated successfully."

