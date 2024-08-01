package main

import (
	"go-backend-crm/config"
	"go-backend-crm/handlers"
	"log"
	"net/http"
)

func main() {
	config.LoadConfig()
	err := handlers.LoadCustomers()
	if err != nil {
		log.Fatalf("Failed to load customers: %v", err)
	}

	r := handlers.NewRouter()
	log.Printf("Server running at port %s", config.GetConfig().Port)
	log.Fatal(http.ListenAndServe(":"+config.GetConfig().Port, r))
}
