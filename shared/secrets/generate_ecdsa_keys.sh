#!/bin/bash

# Define output directory
KEY_DIR="./secrets"
PRIVATE_KEY="${KEY_DIR}/private_key.pem"
PUBLIC_KEY="${KEY_DIR}/public_key.pem"

# Create directory if it doesn't exist
mkdir -p "$KEY_DIR"

# Generate ECDSA private key using NIST P-256 curve
openssl ecparam -name prime256v1 -genkey -noout -out "$PRIVATE_KEY"

# Generate public key from private key
openssl ec -in "$PRIVATE_KEY" -pubout -out "$PUBLIC_KEY"

# Set file permissions (read/write for owner only)
chmod 600 "$PRIVATE_KEY"
chmod 644 "$PUBLIC_KEY"

echo "✅ ECDSA key pair generated:"
echo "   Private Key: $PRIVATE_KEY"
echo "   Public Key : $PUBLIC_KEY"
