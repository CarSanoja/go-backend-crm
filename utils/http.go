package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Configuración de GODEBUG
func configureGODEBUG() {
	os.Setenv("GODEBUG", "http2debug=1,gctrace=1")
}

// Estructura de la respuesta
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Estructura de la configuración del servidor
type Config struct {
	Port      string
	UploadDir string
}

// Middleware para autenticación básica
func basicAuth(next http.HandlerFunc, username, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			http.Error(w, "No autorizado", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Manejadores de las diferentes rutas y métodos HTTP
func handleRequests(config *Config) {
	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/post", handlePost)
	http.HandleFunc("/put", handlePut)
	http.HandleFunc("/delete", handleDelete)
	http.HandleFunc("/upload", handleUpload(config))
}

// Función para manejar peticiones GET
func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()
	response := Response{Status: http.StatusOK, Message: "GET exitoso", Data: queryParams}
	writeJSONResponse(w, response)
}

// Función para manejar peticiones POST
func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		response := Response{Status: http.StatusOK, Message: "POST exitoso", Data: data}
		writeJSONResponse(w, response)
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Datos de formulario inválidos", http.StatusBadRequest)
			return
		}
		response := Response{Status: http.StatusOK, Message: "POST exitoso", Data: r.PostForm}
		writeJSONResponse(w, response)
	} else {
		http.Error(w, "Tipo de contenido no soportado", http.StatusUnsupportedMediaType)
	}
}

// Función para manejar peticiones PUT
func handlePut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		response := Response{Status: http.StatusOK, Message: "PUT exitoso", Data: data}
		writeJSONResponse(w, response)
	} else {
		http.Error(w, "Tipo de contenido no soportado", http.StatusUnsupportedMediaType)
	}
}

// Función para manejar peticiones DELETE
func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()
	response := Response{Status: http.StatusOK, Message: "DELETE exitoso", Data: queryParams}
	writeJSONResponse(w, response)
}

// Función para manejar carga de archivos
func handleUpload(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, "Error al analizar datos del formulario", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error al obtener el archivo", http.StatusBadRequest)
			return
		}
		defer file.Close()

		f, err := os.OpenFile(config.UploadDir+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, "Error al guardar el archivo", http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		response := Response{Status: http.StatusOK, Message: "Archivo subido exitosamente", Data: handler.Filename}
		writeJSONResponse(w, response)
	}
}

// Función para escribir respuestas JSON
func writeJSONResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)
}

// Función principal
func main() {
	config := &Config{
		Port:      "8080",
		UploadDir: "./uploads",
	}
	configureGODEBUG()

	// Creación del directorio de subida si no existe
	if _, err := os.Stat(config.UploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(config.UploadDir, os.ModePerm)
		if err != nil {
			log.Fatalf("No se pudo crear el directorio de subida: %v", err)
		}
	}

	handleRequests(config)
	log.Printf("Servidor corriendo en el puerto %s\n", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

// Funciones de prueba

// TestHandleGet verifica que el manejador GET funcione correctamente.
func TestHandleGet(t *testing.T) {
	req, err := http.NewRequest("GET", "/get?param1=value1&param2=value2", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleGet)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":200,"message":"GET exitoso","data":{"param1":["value1"],"param2":["value2"]}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestHandlePostJSON verifica que el manejador POST funcione correctamente para datos JSON.
func TestHandlePostJSON(t *testing.T) {
	var jsonStr = []byte(`{"name":"test","age":30}`)
	req, err := http.NewRequest("POST", "/post", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlePost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":200,"message":"POST exitoso","data":{"age":30,"name":"test"}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestHandlePostForm verifica que el manejador POST funcione correctamente para datos de formulario.
func TestHandlePostForm(t *testing.T) {
	req, err := http.NewRequest("POST", "/post", strings.NewReader("name=test&age=30"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlePost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":200,"message":"POST exitoso","data":{"age":["30"],"name":["test"]}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestHandleUpload verifica que el manejador de carga de archivos funcione correctamente.
func TestHandleUpload(t *testing.T) {
	config := &Config{UploadDir: "./uploads_test"}
	if _, err := os.Stat(config.UploadDir); os.IsNotExist(err) {
		os.MkdirAll(config.UploadDir, os.ModePerm)
	}
	defer os.RemoveAll(config.UploadDir)

	fileContent := []byte("Este es el contenido del archivo de prueba.")
	fileName := "testfile.txt"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		t.Fatal(err)
	}
	part.Write(fileContent)
	writer.Close()

	req, err := http.NewRequest("POST", "/upload", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := handleUpload(config)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := fmt.Sprintf(`{"status":200,"message":"Archivo subido exitosamente","data":"%s"}`, fileName)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	uploadedFile, err := ioutil.ReadFile(config.UploadDir + "/" + fileName)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(uploadedFile, fileContent) {
		t.Errorf("uploaded file content does not match: got %v want %v", string(uploadedFile), string(fileContent))
	}
}
