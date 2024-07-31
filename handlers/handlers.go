package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"go-backend-crm/config"
	"go-backend-crm/errors"
)

// Estructura de la respuesta
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	response := Response{Status: http.StatusOK, Message: "GET exitoso", Data: queryParams}
	writeJSONResponse(w, response)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			errors.HandleError(w, errors.NewCustomError("JSON inválido", http.StatusBadRequest))
			return
		}
		response := Response{Status: http.StatusOK, Message: "POST exitoso", Data: data}
		writeJSONResponse(w, response)
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		err := r.ParseForm()
		if err != nil {
			errors.HandleError(w, errors.NewCustomError("Datos de formulario inválidos", http.StatusBadRequest))
			return
		}
		response := Response{Status: http.StatusOK, Message: "POST exitoso", Data: r.PostForm}
		writeJSONResponse(w, response)
	} else {
		errors.HandleError(w, errors.NewCustomError("Tipo de contenido no soportado", http.StatusUnsupportedMediaType))
	}
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			errors.HandleError(w, errors.NewCustomError("JSON inválido", http.StatusBadRequest))
			return
		}
		response := Response{Status: http.StatusOK, Message: "PUT exitoso", Data: data}
		writeJSONResponse(w, response)
	} else {
		errors.HandleError(w, errors.NewCustomError("Tipo de contenido no soportado", http.StatusUnsupportedMediaType))
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	response := Response{Status: http.StatusOK, Message: "DELETE exitoso", Data: queryParams}
	writeJSONResponse(w, response)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		errors.HandleError(w, errors.NewCustomError("Error al analizar datos del formulario", http.StatusBadRequest))
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		errors.HandleError(w, errors.NewCustomError("Error al obtener el archivo", http.StatusBadRequest))
		return
	}
	defer file.Close()

	f, err := os.OpenFile(config.GetConfig().UploadDir+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		errors.HandleError(w, errors.NewCustomError("Error al guardar el archivo", http.StatusInternalServerError))
		return
	}
	defer f.Close()
	io.Copy(f, file)

	response := Response{Status: http.StatusOK, Message: "Archivo subido exitosamente", Data: handler.Filename}
	writeJSONResponse(w, response)
}

func handleSecure(w http.ResponseWriter, r *http.Request) {
	response := Response{Status: http.StatusOK, Message: "Acceso seguro autorizado"}
	writeJSONResponse(w, response)
}

func writeJSONResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)
}
