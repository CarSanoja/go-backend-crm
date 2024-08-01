#!/bin/bash

MODULE_NAME="go-backend-crm"

echo "Initializing Go module..."
go mod init $MODULE_NAME

echo "Installing necessary dependencies..."

# Install gorilla/mux for HTTP routing
go get github.com/gorilla/mux@latest

# Install rs/cors for CORS handling
go get github.com/rs/cors@latest

# Install spf13/viper for configuration management
go get github.com/spf13/viper@latest

# Install dgrijalva/jwt-go for JWT handling
go get github.com/dgrijalva/jwt-go@latest

# Install stretchr/testify for testing
go get github.com/stretchr/testify@latest

# Install excelize for Excel file manipulation
go get github.com/xuri/excelize/v2

echo "Installation complete."
