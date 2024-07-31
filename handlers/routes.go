package handlers

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.Use(loggingMiddleware)
	r.Use(corsMiddleware)

	// Rutas públicas
	r.HandleFunc("/get", handleGet).Methods("GET")
	r.HandleFunc("/post", handlePost).Methods("POST")
	r.HandleFunc("/put", handlePut).Methods("PUT")
	r.HandleFunc("/delete", handleDelete).Methods("DELETE")
	r.HandleFunc("/upload", handleUpload).Methods("POST")

	// Rutas seguras
	secure := r.PathPrefix("/secure").Subrouter()
	secure.Use(jwtAuth)
	secure.HandleFunc("", handleSecure).Methods("GET")

	return r
}
