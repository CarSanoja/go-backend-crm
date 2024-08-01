# CRM Backend API

This project represents the backend of a customer relationship management (CRM) web application. It allows users to manage customer data via HTTP requests, supporting various CRUD operations. The server can respond with both HTML and JSON based on the request headers.

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installing Dependencies](#installing-dependencies)
  - [Configuration](#configuration)
- [Running the Server](#running-the-server)
- [API Endpoints](#api-endpoints)
  - [Get All Customers](#1-get-all-customers)
  - [Get a Single Customer](#2-get-a-single-customer)
  - [Add a New Customer](#3-add-a-new-customer)
  - [Update a Customer](#4-update-a-customer)
  - [Delete a Customer](#5-delete-a-customer)
- [HTML Templates](#html-templates)
- [License](#license)

## Getting Started

To get started with this project, you need to install the necessary dependencies and set up the configuration.

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.16 or higher)

### Installing Dependencies

A script is provided to install all necessary dependencies. Run the following command:

```
./scripts/install_lib.sh
```

# Configuration

The application configuration is managed using a config.yaml file. Ensure you have the following structure in your config.yaml:

# yaml

port: "8080"
csv_file: "./data/customers.csv"

# Running the Server

To run the server, execute the following command:

```
go run cmd/main.go
```

The server will start and listen on the port specified in the config.yaml file (default is 8080).
API Endpoints

The API provides the following endpoints for managing customer data:

1. Get All Customers
```
    Endpoint: /customers
    Method: GET
    Description: Retrieve a list of all customers.
    Response: A JSON array of customer objects.
```
Example Request:


```
GET http://localhost:8080/customers
```

Example Response:


```
[
    {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "name": "John Doe",
        "email": "john.doe@example.com",
        "phone": "123-456-7890",
        "address": "123 Main St"
    },
    ...
]
```

2. Get a Single Customer
```
    Endpoint: /customers/{id}
    Method: GET
    Description: Retrieve a single customer by their ID.
    Response: A JSON object representing the customer.
```
Example Request:


```
GET http://localhost:8080/customers/123e4567-e89b-12d3-a456-426614174000
```

Example Response:

```
{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "123-456-7890",
    "address": "123 Main St"
}
```

3. Add a New Customer
```
    Endpoint: /customers/add
    Method: POST
    Description: Add a new customer. The customer ID will be automatically generated.
    Request Body (JSON):
```

```
{
    "name": "Jane Doe",
    "email": "jane.doe@example.com",
    "phone": "987-654-3210",
    "address": "456 Elm St"
}
```

    Response: A JSON object representing the newly created customer.

Example Request:



POST http://localhost:8080/customers/add
Content-Type: application/json

```
{
    "name": "Jane Doe",
    "email": "jane.doe@example.com",
    "phone": "987-654-3210",
    "address": "456 Elm St"
}
```

Example Response:


```
{
    "id": "223e4567-e89b-12d3-a456-426614174001",
    "name": "Jane Doe",
    "email": "jane.doe@example.com",
    "phone": "987-654-3210",
    "address": "456 Elm St"
}
```

4. Update a Customer
```
    Endpoint: /customers/update/{id}
    Method: PUT
    Description: Update an existing customer's information by their ID.
    Request Body (JSON):
```

```
{
    "name": "Jane Doe Updated",
    "email": "jane.doe.updated@example.com",
    "phone": "987-654-3210",
    "address": "456 Elm St Updated"
}
```

    Response: A JSON object representing the updated customer.

Example Request:


```
PUT http://localhost:8080/customers/update/223e4567-e89b-12d3-a456-426614174001
Content-Type: application/json

{
    "name": "Jane Doe Updated",
    "email": "jane.doe.updated@example.com",
    "phone": "987-654-3210",
    "address": "456 Elm St Updated"
}
```

Example Response:


```
{
    "id": "223e4567-e89b-12d3-a456-426614174001",
    "name": "Jane Doe Updated",
    "email": "jane.doe.updated@example.com",
    "phone": "987-654-3210",
    "address": "456 Elm St Updated"
}
```

5. Delete a Customer

```
    Endpoint: /customers/delete/{id}
    Method: DELETE
    Description: Delete a customer by their ID.
    Response: A JSON object indicating success or failure.
```

Example Request:

```
DELETE http://localhost:8080/customers/delete/223e4567-e89b-12d3-a456-426614174001
```

Example Response:

```
{
    "message": "Customer deleted"
}
```

HTML Templates

The project includes HTML templates for rendering the UI:

    home.html: Home page with information about the API and endpoints.
    index.html: List of all customers.
    view.html: View details of a single customer.
    add.html: Form to add a new customer.
    update.html: Form to update an existing customer.
    delete.html: Confirm deletion of a customer.