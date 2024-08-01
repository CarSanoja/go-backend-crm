#!/bin/bash

MODULE_NAME="go-backend-crm"

echo "Initializing Go module..."
go mod init $MODULE_NAME

echo "Installing necessary dependencies..."

# Install UUID package
go get github.com/google/uuid

# Install other dependencies
go get github.com/spf13/viper
go get github.com/gorilla/mux
go get github.com/rs/cors
go get github.com/stretchr/testify

echo "Installation complete."
