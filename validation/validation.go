package validation

import (
	"errors"
	"go-backend-crm/models"
	"regexp"
)

func ValidateCustomer(customer models.Customer) error {
	if customer.ID == "" || customer.Name == "" || customer.Email == "" || customer.Phone == "" || customer.Address == "" {
		return errors.New("all fields are required")
	}

	if !isValidEmail(customer.Email) {
		return errors.New("invalid email format")
	}

	if !isValidPhone(customer.Phone) {
		return errors.New("invalid phone number format")
	}

	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func isValidPhone(phone string) bool {
	re := regexp.MustCompile(`^\d{10}$`)
	return re.MatchString(phone)
}
