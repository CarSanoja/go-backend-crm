package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"go-backend-crm/handlers"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Error al leer el archivo de configuración: %s", err))
	}
}

func TestHandleGet(t *testing.T) {
	req, err := http.NewRequest("GET", "/get?param1=value1&param2=value2", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandleGet)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":200,"message":"GET exitoso","data":{"param1":["value1"],"param2":["value2"]}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandlePostJSON(t *testing.T) {
	payload := []byte(`{"key1":"value1","key2":"value2"}`)
	req, err := http.NewRequest("POST", "/post", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandlePost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":200,"message":"POST exitoso","data":{"key1":"value1","key2":"value2"}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandlePostForm(t *testing.T) {
	form := url.Values{}
	form.Add("field1", "value1")
	form.Add("field2", "value2")

	req, err := http.NewRequest("POST", "/post", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandlePost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":200,"message":"POST exitoso","data":{"field1":["value1"],"field2":["value2"]}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandleUpload(t *testing.T) {
	viper.Set("upload_dir", "./uploads_test")
	if _, err := os.Stat(viper.GetString("upload_dir")); os.IsNotExist(err) {
		os.MkdirAll(viper.GetString("upload_dir"), os.ModePerm)
	}
	defer os.RemoveAll(viper.GetString("upload_dir"))

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
	handler := http.HandlerFunc(handlers.HandleUpload)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := fmt.Sprintf(`{"status":200,"message":"Archivo subido exitosamente","data":"%s"}`, fileName)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	uploadedFile, err := ioutil.ReadFile(viper.GetString("upload_dir") + "/" + fileName)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(uploadedFile, fileContent) {
		t.Errorf("uploaded file content does not match: got %v want %v", string(uploadedFile), string(fileContent))
	}
}

func TestHandleSecure(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(viper.GetString("jwt_secret")))

	req, err := http.NewRequest("GET", "/secure", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr := httptest.NewRecorder()
	handler := handlers.JwtAuth(http.HandlerFunc(handlers.HandleSecure))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":200,"message":"Acceso seguro autorizado"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
