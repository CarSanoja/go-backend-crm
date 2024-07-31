package main

import "net/http"

type CustomError struct {
	Message string
	Code    int
}

func (e *CustomError) Error() string {
	return e.Message
}

func customError(message string, code int) error {
	return &CustomError{Message: message, Code: code}
}

func handleError(w http.ResponseWriter, err error) {
	if customErr, ok := err.(*CustomError); ok {
		http.Error(w, customErr.Message, customErr.Code)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
