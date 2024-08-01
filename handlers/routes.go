package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Use(corsMiddleware)

	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/customers", GetCustomers).Methods("GET")
	r.HandleFunc("/customers/{id}", GetCustomer).Methods("GET")
	r.HandleFunc("/customers/add", CreateCustomer).Methods("POST")
	r.HandleFunc("/customers/update/{id}", UpdateCustomer).Methods("PUT")
	r.HandleFunc("/customers/delete/{id}", DeleteCustomer).Methods("DELETE")

	return r
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	respondWithHTML(w, "home.html", nil)
}
