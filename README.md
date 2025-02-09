# Online Bookstore API

## Overview

The **Online Bookstore API** is a backend system designed for managing books, authors, customers, and orders. It also includes an advanced feature for generating periodic sales reports, providing insights into total revenue, order counts, and top-selling books. This project serves as an educational example of building RESTful APIs, in-memory data storage, concurrency handling, and automated tasks using Go.

## Features

- Manage **Books**: Add, update, delete, and search for books.
- Manage **Authors**: Add, update, delete, and list authors.
- Manage **Customers**: Add, update, delete, and list customers.
- Manage **Orders**: Create, update, delete, and view orders.
- Generate **Sales Reports**: Automatically generate daily reports, including:
  - Total revenue.
  - Total number of orders.
  - Total books sold.
  - Top-selling books.

## Project Structure

```
final-project/
│
├── database/                  # JSON files for data storage
│   ├── authors.json           # Stores author data
│   ├── books.json             # Stores book data
│   ├── customers.json         # Stores customer data
│   ├── orders.json            # Stores order data
│   ├── sales_reports.json     # Stores generated sales reports
│
├── errorhandling/             # Custom error definitions
│   ├── errorHandling.go       # Error handling logic for the API
│
├── handlers/                  # Handlers for API endpoints
│   ├── AuthorsHandlers.go     # Handlers for author-related operations
│   ├── BooksHandlers.go       # Handlers for book-related operations
│   ├── CustomersHandlers.go   # Handlers for customer-related operations
│   ├── OrdersHandlers.go      # Handlers for order-related operations
│
├── inmemorystores/            # In-memory data storage and operations
│   ├── inMemoryAuthorStore.go # Storage and operations for authors
│   ├── inMemoryBookStore.go   # Storage and operations for books
│   ├── inMemoryCustomerStore.go # Storage and operations for customers
│   ├── inMemoryOrderStore.go  # Storage and operations for orders
│   ├── inMemoryReportStore.go # Storage and operations for sales reports
│
├── interfaces/                # Interface definitions
│   ├── interfaces.go          # Interfaces for data stores
│
├── models/                    # Data models for the project
│   ├── models.go              # Defines Book, Author, Customer, Order, and Report models
│
├── swagger.json               # API documentation in OpenAPI format
├── main.go                    # Entry point of the application
├── go.mod                     # Dependency file for Go modules
├── go.sum                     # Checksum file for dependencies
├── README.md                  # Documentation (this file)
```

### Folder Explanations

1. **`database/`**:

   - Houses `.json` files that serve as persistent storage for various entities:
     - `authors.json`: Stores data about authors.
     - `books.json`: Contains information about books.
     - `customers.json`: Maintains customer records.
     - `orders.json`: Tracks order details.
     - `sales_reports.json`: Saves generated sales reports for analytics.

2. **`errorhandling/`**:

   - Centralizes error handling with custom error definitions like `ErrBookNotFound` and `ErrOrderNotFound`, ensuring consistent and clear error messages throughout the API.

3. **`handlers/`**:

   - Contains the core logic for handling API requests, with each file focused on a specific entity:
     - `AuthorsHandlers.go`: Manages author-related operations.
     - `BooksHandlers.go`: Handles book-related CRUD operations.
     - `CustomersHandlers.go`: Oversees customer-related actions.
     - `OrdersHandlers.go`: Implements order management logic.

4. **`inmemorystores/`**:

   - Implements in-memory data storage with methods for CRUD operations. Each file encapsulates logic for a specific entity:
     - `inMemoryAuthorStore.go`: Manages author data.
     - `inMemoryBookStore.go`: Handles book-related storage and queries.
     - `inMemoryCustomerStore.go`: Stores customer details.
     - `inMemoryOrderStore.go`: Tracks orders and supports report generation.
     - `inMemoryReportStore.go`: Stores and manages generated sales reports.

5. **`interfaces/`**:

   - Defines abstraction layers for data stores, promoting modularity and simplifying testing by allowing interchangeable implementations.

6. **`models/`**:

   - Defines data models used throughout the application:
     - `Book`, `Author`, `Order`, `Customer`: Represent core entities.
     - `SalesReport`: Structure for periodic reports.
     - `SearchCriteria`: Used for filtering books during searches.

7. **`swagger.json`**:

   - Contains the OpenAPI specification for the API. This file facilitates documentation and testing with tools like Swagger UI, enabling developers to visualize and interact with the endpoints effortlessly.

8. **`main.go`**:
   - Serves as the application's entry point, handling:
     - Server setup and initialization.
     - Configuration of in-memory stores.
     - Registration of API endpoints.
     - Background task execution, such as periodic sales report generation, ensuring smooth application functionality.

## How to Run the Project

### Prerequisites

- Install [Go](https://go.dev/) on your system.
- Clone the repository to your local machine.

### Steps

1. **Clone the repository**:

   ```bash
   git clone https://github.com/malakkbl/GoProgramming.git
   cd final-project
   ```

2. **Run the application**:

   ```bash
   go run main.go
   ```

3. **Access the API**:
   - The server will start on `http://localhost:8080`.
   - Use tools like [Postman](https://www.postman.com/) or [cURL](https://curl.se/) to interact with the API.

### Endpoints

#### Books

- `GET /books`: Retrieve all books or search using query parameters.
- `POST /books`: Add a new book.
- `GET /books/{id}`: Retrieve a book by its ID.
- `PUT /books/{id}`: Update book details.
- `DELETE /books/{id}`: Delete a book.

#### Authors

- `GET /authors`: Retrieve all authors.
- `POST /authors`: Add a new author.
- `GET /authors/{id}`: Retrieve an author by their ID.
- `PUT /authors/{id}`: Update author details.
- `DELETE /authors/{id}`: Delete an author.

#### Customers

- `GET /customers`: Retrieve all customers.
- `POST /customers`: Add a new customer.
- `GET /customers/{id}`: Retrieve a customer by their ID.
- `PUT /customers/{id}`: Update customer details.
- `DELETE /customers/{id}`: Delete a customer.

#### Orders

- `GET /orders`: Retrieve all orders.
- `POST /orders`: Create a new order.
- `GET /orders/{id}`: Retrieve an order by its ID.
- `PUT /orders/{id}`: Update order details.
- `DELETE /orders/{id}`: Delete an order.

#### Sales Reports

- `GET /sales-reports`: Retrieve all generated sales reports.

### Sales Report Generation

- Sales reports are automatically generated every 24 hours in the background.
- Reports include:
  - Total revenue.
  - Total number of orders.
  - Total books sold.
  - Top-selling books.
- Reports are saved to `database/sales_reports.json` and can be accessed via the `/sales-reports` endpoint.

## API Documentation

The `swagger.json` file provides OpenAPI documentation for the API. You can use tools like [Swagger UI](https://swagger.io/tools/swagger-ui/) or [Postman](https://www.postman.com/) to visualize and test the endpoints.

## Testing :
These manual tests are to ensure all endpoints are functional and adhere to the project’s structure and logic:

### **Books Endpoints**

1. **Create a Book**  
   **POST** `/books`  
   **Body**:

   ```json
   {
     "id": 1,
     "title": "The Go Programming Language",
     "author": {
       "id": 1,
       "first_name": "Alan",
       "last_name": "Donovan",
       "bio": "Author and programmer."
     },
     "genres": ["Programming", "Technology"],
     "published_at": "2015-10-26T00:00:00Z",
     "price": 45.99,
     "stock": 100
   }
   ```

   **Expected Response**:

   ```json
   {
     "id": 1,
     "title": "The Go Programming Language",
     "author": {
       "id": 1,
       "first_name": "Alan",
       "last_name": "Donovan",
       "bio": "Author and programmer."
     },
     "genres": ["Programming", "Technology"],
     "published_at": "2015-10-26T00:00:00Z",
     "price": 45.99,
     "stock": 100
   }
   ```

2. **Get a Book by ID**  
   **GET** `/books/1`  
   **Expected Response**:

   ```json
   {
     "id": 1,
     "title": "The Go Programming Language",
     "author": {
       "id": 1,
       "first_name": "Alan",
       "last_name": "Donovan",
       "bio": "Author and programmer."
     },
     "genres": ["Programming", "Technology"],
     "published_at": "2015-10-26T00:00:00Z",
     "price": 45.99,
     "stock": 100
   }
   ```

3. **Update a Book**  
   **PUT** `/books/1`  
   **Body**:

   ```json
   {
     "id": 1,
     "title": "Updated Book Title",
     "author": {
       "id": 1,
       "first_name": "Alan",
       "last_name": "Donovan",
       "bio": "Author and programmer."
     },
     "genres": ["Programming"],
     "published_at": "2015-10-26T00:00:00Z",
     "price": 50.0,
     "stock": 90
   }
   ```

   **Expected Response**:

   ```json
   {
     "id": 1,
     "title": "Updated Book Title",
     "author": {
       "id": 1,
       "first_name": "Alan",
       "last_name": "Donovan",
       "bio": "Author and programmer."
     },
     "genres": ["Programming"],
     "published_at": "2015-10-26T00:00:00Z",
     "price": 50.0,
     "stock": 90
   }
   ```

4. **Delete a Book**  
   **DELETE** `/books/1`  
   **Expected Response**:

   ```
   HTTP 204 No Content
   ```

5. **Search Books**  
   **GET** `/books?title=Go`  
   **Expected Response**:
   ```json
   [
     {
       "id": 1,
       "title": "The Go Programming Language",
       "author": {
         "id": 1,
         "first_name": "Alan",
         "last_name": "Donovan",
         "bio": "Author and programmer."
       },
       "genres": ["Programming", "Technology"],
       "published_at": "2015-10-26T00:00:00Z",
       "price": 45.99,
       "stock": 100
     }
   ]
   ```

---

### **Authors Endpoints**

1. **Create an Author**  
   **POST** `/authors`  
   **Body**:

   ```json
   {
     "id": 1,
     "first_name": "Alan",
     "last_name": "Donovan",
     "bio": "Author and programmer."
   }
   ```

   **Expected Response**:

   ```json
   {
     "id": 1,
     "first_name": "Alan",
     "last_name": "Donovan",
     "bio": "Author and programmer."
   }
   ```

2. **Get an Author by ID**  
   **GET** `/authors/1`  
   **Expected Response**:

   ```json
   {
     "id": 1,
     "first_name": "Alan",
     "last_name": "Donovan",
     "bio": "Author and programmer."
   }
   ```

3. **Update an Author**  
   **PUT** `/authors/1`  
   **Body**:

   ```json
   {
     "id": 1,
     "first_name": "Updated First Name",
     "last_name": "Updated Last Name",
     "bio": "Updated bio."
   }
   ```

   **Expected Response**:

   ```json
   {
     "id": 1,
     "first_name": "Updated First Name",
     "last_name": "Updated Last Name",
     "bio": "Updated bio."
   }
   ```

4. **Delete an Author**  
   **DELETE** `/authors/1`  
   **Expected Response**:

   ```
   HTTP 204 No Content
   ```

5. **List All Authors**  
   **GET** `/authors`  
   **Expected Response**:
   ```json
   [
     {
       "id": 1,
       "first_name": "Alan",
       "last_name": "Donovan",
       "bio": "Author and programmer."
     }
   ]
   ```

---

### **Customers Endpoints**

1. **Create a Customer**  
   **POST** `/customers`  
   **Body**:

   ```json
   {
     "id": 1,
     "name": "John Doe",
     "email": "john.doe@example.com",
     "address": {
       "street": "123 Main St",
       "city": "Springfield",
       "state": "IL",
       "postal_code": "62701",
       "country": "USA"
     },
     "created_at": "2025-01-12T00:00:00Z"
   }
   ```

   **Expected Response**:

   ```json
   {
     "id": 1,
     "name": "John Doe",
     "email": "john.doe@example.com",
     "address": {
       "street": "123 Main St",
       "city": "Springfield",
       "state": "IL",
       "postal_code": "62701",
       "country": "USA"
     },
     "created_at": "2025-01-12T00:00:00Z"
   }
   ```

2. **Get a Customer by ID**  
   **GET** `/customers/1`  
   **Expected Response**:
   ```json
   {
     "id": 1,
     "name": "John Doe",
     "email": "john.doe@example.com",
     "address": {
       "street": "123 Main St",
       "city": "Springfield",
       "state": "IL",
       "postal_code": "62701",
       "country": "USA"
     },
     "created_at": "2025-01-12T00:00:00Z"
   }
   ```

---

### **Orders Endpoints**

1. **Create an Order**  
   **POST** `/orders`  
   **Body**:

   ```json
   {
     "id": 1,
     "customer": {
       "id": 1,
       "name": "John Doe",
       "email": "john.doe@example.com",
       "address": {
         "street": "123 Main St",
         "city": "Springfield",
         "state": "IL",
         "postal_code": "62701",
         "country": "USA"
       },
       "created_at": "2025-01-12T00:00:00Z"
     },
     "items": [
       {
         "book": {
           "id": 1,
           "title": "The Go Programming Language",
           "author": {
             "id": 1,
             "first_name": "Alan",
             "last_name": "Donovan",
             "bio": "Author and programmer."
           },
           "genres": ["Programming", "Technology"],
           "published_at": "2015-10-26T00:00:00Z",
           "price": 45.99,
           "stock": 100
         },
         "quantity": 1
       }
     ],
     "total_price": 45.99,
     "created_at": "2025-01-12T00:00:00Z",
     "status": "Processing"
   }
   ```

   **Expected Response**:

   ```json
   {
     "id": 1,
     "customer": {
       "id": 1,
       "name": "John Doe",
       "email": "john.doe@example.com",
       "address": {
         "street": "123 Main St",
         "city": "Springfield",
         "state": "IL",
         "postal_code": "62701",
         "country": "USA"
       },
       "created_at": "2025-01-12T00:00:00Z"
     },
     "items": [
       {
         "book": {
           "id": 1,
           "title": "The Go Programming Language",
           "author": {
             "id": 1,
             "first_name": "Alan",
             "last_name": "Donovan",
             "bio": "Author and programmer."
           },
           "genres": ["Programming", "Technology"],
           "published_at": "2015-10-26T00:00:00Z",
           "price": 45.99,
           "stock": 100
         },
         "quantity": 1
       }
     ],
     "total_price": 45.99,
     "created_at": "2025-01-12T00:00:00Z",
     "status": "Processing"
   }
   ```

2. **Get All Orders**  
   **GET** `/orders`  
   **Expected Response**:
   ```json
   [
     {
       "id": 1,
       "customer": {
         "id": 1,
         "name": "John Doe",
         "email": "john.doe@example.com",
         "address": {
           "street": "123 Main St",
           "city": "Springfield",
           "state": "IL",
           "postal_code": "62701",
           "country": "USA"
         },
         "created_at": "2025-01-12T00:00:00Z"
       },
       "items": [
         {
           "book": {
             "id": 1,
             "title": "The Go Programming Language",
             "author": {
               "id": 1,
               "first_name": "Alan",
               "last_name": "Donovan",
               "bio": "Author and programmer."
             },
             "genres": ["Programming", "Technology"],
             "published_at": "2015-10-26T00:00:00Z",
             "price": 45.99,
             "stock": 100
           },
           "quantity": 1
         }
       ],
       "total_price": 45.99,
       "created_at": "2025-01-12T00:00:00Z",
       "status": "Processing"
     }
   ]
   ```

---

### **Sales Reports Endpoints**

1. **Get Sales Reports**  
   **GET** `/sales-reports`  
   **Expected Response**:
   ```json
   [
     {
       "timestamp": "2025-01-12T00:00:00Z",
       "total_revenue": 100.0,
       "total_orders": 2,
       "top_selling_books": [
         {
           "book": {
             "id": 1,
             "title": "The Go Programming Language",
             "author": {
               "id": 1,
               "first_name": "Alan",
               "last_name": "Donovan",
               "bio": "Author and programmer."
             },
             "genres": ["Programming", "Technology"],
             "published_at": "2015-10-26T00:00:00Z",
             "price": 45.99,
             "stock": 100
           },
           "quantity": 10
         }
       ]
     }
   ]
   ```


