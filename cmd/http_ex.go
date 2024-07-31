package main

import (
	"log"
	"net/http"

	"go-backend-crm/config"
	"go-backend-crm/handlers"
)

func main() {
	config.LoadConfig()

	r := handlers.NewRouter()

	addr := ":" + config.GetConfig().Port
	log.Printf("Servidor corriendo en el puerto %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
