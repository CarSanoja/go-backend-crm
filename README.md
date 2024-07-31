# Go Backend CRM

This project is a template for a Go backend that includes authentication and authorization with JWT, logging and monitoring, flexible configuration, custom error handling, CORS handling, and input validation. It also includes unit and integration tests.
Project Structure

go

your-project/
├── cmd/
│   └── main.go
├── config/
│   └── config.go
├── errors/
│   └── errors.go
├── handlers/
│   ├── handlers.go
│   ├── middleware.go
│   └── routes.go
├── scripts/
│   └── install_lib.sh
├── validation/
│   └── validation.go
├── tests/
│   └── main_test.go
├── config.yaml
└── go.mod

Directory Description

    cmd/: Contains the main entry point main.go for starting the application.
    config/: Handles server configuration in config.go.
    errors/: Defines custom errors in errors.go.
    handlers/: Contains route handlers (handlers.go), middlewares (middleware.go), and route configuration (routes.go).
    scripts/: Contains the install_lib.sh script for initializing the Go module and installing dependencies.
    validation/: Handles input validation in validation.go.
    tests/: Contains unit and integration tests in main_test.go.
    config.yaml: Server configuration file.
    go.mod: Go module configuration file.

Prerequisites

    Go 1.16 or higher installed on your machine.

Installation
Step 1: Clone the Repository

Clone this repository to your local machine:

bash

git clone https://github.com/your-username/your-repository.git
cd your-repository

Step 2: Run the Installation Script

The install_lib.sh script initializes the Go module and installs the necessary dependencies.

bash

./scripts/install_lib.sh your-module-name

Replace your-module-name with the desired name for your module.
Step 3: Configure the config.yaml File

Edit the config.yaml file with the necessary values:

yaml

port: "8080"
upload_dir: "./uploads"
jwt_secret: "your_jwt_secret"

Running the Server

To run the server, use the following command:

bash

go run cmd/main.go

Testing
Running Tests

To run the unit and integration tests, use the following command:

bash

go test -v ./tests

Features
Authentication and Authorization

    JWT Authentication: Protects secure routes with JWT.

Logging and Monitoring

    Logging Middleware: Logs all HTTP requests.

Flexible Configuration

    Configuration with Viper: Manages configurations through a config.yaml file.

Error Handling

    Custom Errors: Defines and handles custom errors.

CORS

    CORS Handling: Allows requests from any origin.

Input Validation

    Validates emails and other input data.

Routes and Endpoints

    GET /get: Handles GET requests.
    POST /post: Handles POST requests (JSON and form-urlencoded).
    PUT /put: Handles PUT requests (JSON).
    DELETE /delete: Handles DELETE requests.
    POST /upload: Handles file uploads.
    GET /secure: Secure route protected with JWT.

Usage Examples
GET Request

bash

curl "http://localhost:8080/get?param1=value1&param2=value2"

POST Request (JSON)

bash

curl -X POST "http://localhost:8080/post" -H "Content-Type: application/json" -d '{"key1":"value1","key2":"value2"}'

POST Request (Form-urlencoded)

bash

curl -X POST "http://localhost:8080/post" -H "Content-Type: application/x-www-form-urlencoded" -d "field1=value1&field2=value2"

PUT Request

bash

curl -X PUT "http://localhost:8080/put" -H "Content-Type: application/json" -d '{"key1":"value1","key2":"value2"}'

DELETE Request

bash

curl -X DELETE "http://localhost:8080/delete?param1=value1"

File Upload Request

bash

curl -X POST "http://localhost:8080/upload" -F "file=@path/to/your/file"

Secure Route Request (JWT)

bash

TOKEN=$(curl -X POST "http://localhost:8080/login" -d '{"username":"your-username","password":"your-password"}' | jq -r .token)
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/secure"

Contributing

Contributions are welcome. Please open an issue or a pull request to contribute.