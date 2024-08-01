package handlers

import (
	"encoding/csv"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"go-backend-crm/config"
	"go-backend-crm/errors"
	"go-backend-crm/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var customers []models.Customer

func LoadCustomers() error {
	file, err := os.Open(config.GetConfig().CSVFile)
	if err != nil {
		return err
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	customers = []models.Customer{}
	for _, record := range records[1:] {
		customer := models.Customer{
			ID:      record[0],
			Name:    record[1],
			Email:   record[2],
			Phone:   record[3],
			Address: record[4],
		}
		customers = append(customers, customer)
	}
	return nil
}

func SaveCustomers() error {
	file, err := os.Create(config.GetConfig().CSVFile)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	w.Write([]string{"ID", "Name", "Email", "Phone", "Address"})
	for _, customer := range customers {
		w.Write([]string{customer.ID, customer.Name, customer.Email, customer.Phone, customer.Address})
	}
	return nil
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func respondWithHTML(w http.ResponseWriter, templateName string, data interface{}) {
	t, err := template.ParseFiles("handlers/views/" + templateName)
	if err != nil {
		errors.HandleError(w, err)
		return
	}
	t.Execute(w, data)
}

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	acceptHeader := r.Header.Get("Accept")
	if acceptHeader == "application/json" {
		respondWithJSON(w, http.StatusOK, customers)
	} else {
		respondWithHTML(w, "index.html", customers)
	}
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, customer := range customers {
		if customer.ID == params["id"] {
			acceptHeader := r.Header.Get("Accept")
			if acceptHeader == "application/json" {
				respondWithJSON(w, http.StatusOK, customer)
			} else {
				respondWithHTML(w, "view.html", customer)
			}
			return
		}
	}
	errors.HandleError(w, errors.NewCustomError("Customer not found", http.StatusNotFound))
}

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errors.HandleError(w, err)
			return
		}
		if err := json.Unmarshal(body, &customer); err != nil {
			errors.HandleError(w, err)
			return
		}
	} else {
		customer.Name = r.FormValue("name")
		customer.Email = r.FormValue("email")
		customer.Phone = r.FormValue("phone")
		customer.Address = r.FormValue("address")
	}

	customer.ID = uuid.New().String()
	customers = append(customers, customer)
	if err := SaveCustomers(); err != nil {
		errors.HandleError(w, err)
		return
	}

	acceptHeader := r.Header.Get("Accept")
	if acceptHeader == "application/json" {
		respondWithJSON(w, http.StatusCreated, customer)
	} else {
		http.Redirect(w, r, "/customers", http.StatusSeeOther)
	}
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedCustomer models.Customer
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errors.HandleError(w, err)
			return
		}
		if err := json.Unmarshal(body, &updatedCustomer); err != nil {
			errors.HandleError(w, err)
			return
		}
	} else {
		updatedCustomer.ID = params["id"]
		updatedCustomer.Name = r.FormValue("name")
		updatedCustomer.Email = r.FormValue("email")
		updatedCustomer.Phone = r.FormValue("phone")
		updatedCustomer.Address = r.FormValue("address")
	}

	for i, customer := range customers {
		if customer.ID == params["id"] {
			customers[i] = updatedCustomer
			if err := SaveCustomers(); err != nil {
				errors.HandleError(w, err)
				return
			}

			acceptHeader := r.Header.Get("Accept")
			if acceptHeader == "application/json" {
				respondWithJSON(w, http.StatusOK, updatedCustomer)
			} else {
				http.Redirect(w, r, "/customers", http.StatusSeeOther)
			}
			return
		}
	}
	errors.HandleError(w, errors.NewCustomError("Customer not found", http.StatusNotFound))
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, customer := range customers {
		if customer.ID == params["id"] {
			customers = append(customers[:i], customers[i+1:]...)
			if err := SaveCustomers(); err != nil {
				errors.HandleError(w, err)
				return
			}

			acceptHeader := r.Header.Get("Accept")
			if acceptHeader == "application/json" {
				respondWithJSON(w, http.StatusOK, map[string]string{"message": "Customer deleted"})
			} else {
				http.Redirect(w, r, "/customers", http.StatusSeeOther)
			}
			return
		}
	}
	errors.HandleError(w, errors.NewCustomError("Customer not found", http.StatusNotFound))
}
