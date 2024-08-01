package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Use(corsMiddleware)

	r.HandleFunc("/", handleIndex).Methods("GET")
	r.HandleFunc("/customers", HandleGetCustomers).Methods("GET")
	r.HandleFunc("/customers/view/{id}", HandleGetCustomer).Methods("GET")
	r.HandleFunc("/customers/add", HandleCreateCustomer).Methods("GET", "POST")
	r.HandleFunc("/customers/update/{id}", HandleUpdateCustomer).Methods("GET", "POST")
	r.HandleFunc("/customers/delete/{id}", HandleDeleteCustomer).Methods("GET", "POST")

	return r
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/customers", http.StatusSeeOther)
}
