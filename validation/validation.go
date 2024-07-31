package validation

import (
	"errors"
	"regexp"
)

// Valida entradas
func ValidateInput(data map[string]interface{}) error {
	for key, value := range data {
		if key == "email" {
			if !isValidEmail(value.(string)) {
				return errors.New("Correo electrónico inválido")
			}
		}
	}
	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}
