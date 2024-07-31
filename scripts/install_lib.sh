#!/bin/bash

MODULE_NAME=$1

if [ -z "$MODULE_NAME" ]; then
  echo "Por favor, proporciona el nombre del módulo como parámetro."
  exit 1
fi

echo "Inicializando módulo de Go..."
go mod init $MODULE_NAME

echo "Instalando dependencias necesarias..."

# Instalar gorilla/mux para enrutamiento de HTTP
go get github.com/gorilla/mux@latest

# Instalar rs/cors para manejo de CORS
go get github.com/rs/cors@latest

# Instalar spf13/viper para manejo de configuración
go get github.com/spf13/viper@latest

# Instalar dgrijalva/jwt-go para manejo de JWT
go get github.com/dgrijalva/jwt-go@latest

# Instalar stretchr/testify para pruebas
go get github.com/stretchr/testify@latest

echo "Instalación completa."
