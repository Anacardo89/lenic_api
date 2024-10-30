#!/bin/bash

CERT_DIR="./cmd/ssl"
mkdir -p "$CERT_DIR"

CERT_FILE="$CERT_DIR/cert.pem"
SIGNED_CERT_FILE="$CERT_DIR/certificate.pem"
KEY_FILE="$CERT_DIR/key.pem"

echo "Generating self-signed certificate..."
openssl genrsa -out "$KEY_FILE"
openssl req -new -key "$KEY_FILE" -subj "/C=PT/L=Porto/O=CINEL" -out "$CERT_FILE"
openssl req -x509 -days 365 -key "$KEY_FILE" -in "$CERT_FILE" -out "$SIGNED_CERT_FILE"

if [ -f "$CERT_FILE" ] && [ -f "$KEY_FILE" ] && [-f "$SIGNED_CERT_FILE"]; then
    echo "SSL certificate and key have been successfully created."
    echo "Certificate: $CERT_FILE"
    echo "Certificate: $SIGNED_CERT_FILE"
    echo "Private Key: $KEY_FILE"
else
    echo "Error: Failed to create SSL certificate and/or key."
    exit 1
fi