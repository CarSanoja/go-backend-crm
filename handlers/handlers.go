package handlers

import (
	"encoding/csv"
	"html/template"
	"net/http"
	"os"

	"go-backend-crm/config"
	"go-backend-crm/errors"
	"go-backend-crm/models"
	"go-backend-crm/validation"

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

func HandleGetCustomers(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("handlers/views/index.html")
	if err != nil {
		errors.HandleError(w, err)
		return
	}
	t.Execute(w, customers)
}

func HandleGetCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, customer := range customers {
		if customer.ID == params["id"] {
			t, err := template.ParseFiles("handlers/views/view.html")
			if err != nil {
				errors.HandleError(w, err)
				return
			}
			t.Execute(w, customer)
			return
		}
	}
	errors.HandleError(w, errors.NewCustomError("Customer not found", http.StatusNotFound))
}

func HandleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("handlers/views/add.html")
		if err != nil {
			errors.HandleError(w, err)
			return
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		var customer models.Customer
		customer.ID = r.FormValue("id")
		customer.Name = r.FormValue("name")
		customer.Email = r.FormValue("email")
		customer.Phone = r.FormValue("phone")
		customer.Address = r.FormValue("address")

		if err := validation.ValidateCustomer(customer); err != nil {
			errors.HandleError(w, err)
			return
		}

		customers = append(customers, customer)
		if err := SaveCustomers(); err != nil {
			errors.HandleError(w, err)
			return
		}

		http.Redirect(w, r, "/customers", http.StatusSeeOther)
	}
}

func HandleUpdateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Method == "GET" {
		for _, customer := range customers {
			if customer.ID == params["id"] {
				t, err := template.ParseFiles("handlers/views/update.html")
				if err != nil {
					errors.HandleError(w, err)
					return
				}
				t.Execute(w, customer)
				return
			}
		}
		errors.HandleError(w, errors.NewCustomError("Customer not found", http.StatusNotFound))
	} else if r.Method == "POST" {
		var updatedCustomer models.Customer
		updatedCustomer.ID = r.FormValue("id")
		updatedCustomer.Name = r.FormValue("name")
		updatedCustomer.Email = r.FormValue("email")
		updatedCustomer.Phone = r.FormValue("phone")
		updatedCustomer.Address = r.FormValue("address")

		for i, customer := range customers {
			if customer.ID == params["id"] {
				customers[i] = updatedCustomer
				if err := SaveCustomers(); err != nil {
					errors.HandleError(w, err)
					return
				}
				http.Redirect(w, r, "/customers", http.StatusSeeOther)
				return
			}
		}
		errors.HandleError(w, errors.NewCustomError("Customer not found", http.StatusNotFound))
	}
}

func HandleDeleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Method == "GET" {
		for _, customer := range customers {
			if customer.ID == params["id"] {
				t, err := template.ParseFiles("handlers/views/delete.html")
				if err != nil {
					errors.HandleError(w, err)
					return
				}
				t.Execute(w, customer)
				return
			}
		}
		errors.HandleError(w, errors.NewCustomError("Customer not found", http.StatusNotFound))
	} else if r.Method == "POST" {
		for i, customer := range customers {
			if customer.ID == params["id"] {
				customers = append(customers[:i], customers[i+1:]...)
				if err := SaveCustomers(); err != nil {
					errors.HandleError(w, err)
					return
				}
				http.Redirect(w, r, "/customers", http.StatusSeeOther)
				return
			}
		}
		errors.HandleError(w, errors.NewCustomError("Customer not found", http.StatusNotFound))
	}
}
