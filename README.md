# MyBiblio - Online Bookstore API

## Description

**MyBiblio** is an advanced RESTful API for managing an online bookstore. It allows users to manage books, authors, customers, and orders while generating periodic sales reports to analyze performance.

---

## Collaborators

- **Kably Malak**
- Project Advisor: **Abdelghafour Mourchid**

---

## Features

### Current Features
- **Books Management**
  - Add, update, delete, and search for books.
- **Authors Management**
  - Add, update, delete, and list authors.
- **Customers Management**
  - Add, update, delete, and list customers.
- **Orders Management**
  - Create, update, delete, and view orders.
- **Sales Reports**
  - Automatically generate daily sales reports, including:
    - Total revenue.
    - Total number of orders.
    - Total books sold.
    - Top-selling books.

### Planned Milestones

Our project follows a structured development roadmap. You can track our progress through the **Milestones** section on GitHub.

#### **Current Milestones:**
- **Milestone 1:** Implement a Database (PostgreSQL with GORM).
- **Milestone 2:** Authentication & Authorization (JWT-based security).
- **Milestone 3:** Input Validation & Error Handling.
- **Milestone 4:** Caching & Comprehensive Testing.
- **Milestone 5:** Containerization, CI/CD Pipeline, Monitoring & Logging.
- **Milestone 6:** Database Migrations using Flyway.

---

## Technical Requirements

- **Go 1.20+**
- **Postman** or **Swagger UI** for API testing.
- **Git** for version control.

---

## Installation Instructions

### Prerequisites
1. Ensure you have **Go** installed ([Install Go](https://go.dev/)).
2. Clone the repository to your local machine:
   ```bash
   git clone https://github.com/yourusername/MyBiblio.git
   cd MyBiblio
   ```

### Steps to Run
1. Run the application:
   ```bash
   go run main.go
   ```
2. Access the API at `http://localhost:8080`.

---

## Usage Examples

### Authentication
First, you need to register and login to get your JWT token:

#### Register a New User
```bash
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{
    "name": "Test Admin",
    "email": "admin@mybiblio.com",
    "password": "adminpass123",
    "role": "admin"
}'
```

#### Login
```bash
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{
    "email": "admin@mybiblio.com",
    "password": "adminpass123"
}'
```

The login response will include a JWT token. Use this token in the Authorization header for subsequent requests:
```bash
export TOKEN="your_jwt_token_here"
```

### Basic Endpoints
#### Books
- **Create a Book (Requires Admin or Manager Role)**
  ```bash
  curl -X POST http://localhost:8080/books \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
      "title": "The Go Programming Language",
      "author": { "id": 1, "first_name": "Alan", "last_name": "Donovan" },
      "genres": ["Programming"],
      "price": 45.99,
      "stock": 100
    }'
  ```

- **Get All Books (Public Access)**
  ```bash
  curl http://localhost:8080/books
  ```

### Role-Based Examples

#### Manager Operations
```bash
# Register as Manager
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{
    "name": "Test Manager",
    "email": "manager@mybiblio.com",
    "password": "managerpass123",
    "role": "manager"
}'

# Create New Book (Allowed for Manager)
curl -X POST http://localhost:8080/books \
  -H "Authorization: Bearer $MANAGER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "New Book",
    "author": {"id": 1},
    "price": 29.99,
    "stock": 50
  }'
```

#### Employee Operations
```bash
# Register as Employee
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{
    "name": "Test Employee",
    "email": "employee@mybiblio.com",
    "password": "employeepass123",
    "role": "employee"
}'

# List Customers (Allowed for Employee)
curl -X GET http://localhost:8080/customers \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN"
```

#### Regular User Operations
```bash
# Register as User
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{
    "name": "Test User",
    "email": "user@mybiblio.com",
    "password": "userpass123",
    "role": "user"
}'

# Create Order (Allowed for User)
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customerID": 1,
    "items": [
      {
        "bookID": 1,
        "quantity": 2
      }
    ]
  }'
```

---

## Role-Based Access Control (RBAC) Testing Guide

### 1. Database Setup
Ensure PostgreSQL is running and your `.env` file contains:
```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=mybiblio
DB_PORT=5432
DB_SSL=disable
JWT_SECRET=your_secure_secret_key
```

### 2. Available Roles
The system supports four roles with different permissions:
- **Admin**: Full access to all features
- **Manager**: Can manage books, authors, and view reports
- **Employee**: Can manage customers and orders
- **User**: Can browse books and manage their own orders

### 3. Testing with Postman

#### Step 1: User Registration and Authentication

1. **Register Users** (Create a request for each role)
   - Method: `POST`
   - URL: `http://localhost:8080/register`
   - Headers: `Content-Type: application/json`
   - Body (Admin):
     ```json
     {
         "name": "Admin User",
         "email": "admin@mybiblio.com",
         "password": "adminpass123",
         "role": "admin"
     }
     ```
   - Body (Manager):
     ```json
     {
         "name": "Manager User",
         "email": "manager@mybiblio.com",
         "password": "managerpass123",
         "role": "manager"
     }
     ```
   - Similar for "employee" and "user" roles

2. **Login** (Test each user)
   - Method: `POST`
   - URL: `http://localhost:8080/login`
   - Headers: `Content-Type: application/json`
   - Body:
     ```json
     {
         "email": "admin@mybiblio.com",
         "password": "adminpass123"
     }
     ```
   Save the returned token for each user.

#### Step 2: Testing Role Permissions

Create these requests in Postman and test with different user tokens:

1. **Admin Tests** (Should all succeed)
   ```http
   # Create Book
   POST http://localhost:8080/books
   Authorization: Bearer <admin_token>
   Content-Type: application/json
   
   {
       "title": "Test Book",
       "author": {
           "id": 1,
           "first_name": "Test Author",
           "last_name": "Test"
       },
       "genres": "Fiction",
       "price": 29.99,
       "stock": 100
   }

   # Delete User
   DELETE http://localhost:8080/users/2
   Authorization: Bearer <admin_token>

   # View Reports
   GET http://localhost:8080/sales-reports
   Authorization: Bearer <admin_token>
   ```

2. **Manager Tests**
   ```http
   # Create Book (Should succeed)
   POST http://localhost:8080/books
   Authorization: Bearer <manager_token>
   
   # Delete User (Should fail)
   DELETE http://localhost:8080/users/2
   Authorization: Bearer <manager_token>
   ```

3. **Employee Tests**
   ```http
   # View Customers (Should succeed)
   GET http://localhost:8080/customers
   Authorization: Bearer <employee_token>
   
   # Create Book (Should fail)
   POST http://localhost:8080/books
   Authorization: Bearer <employee_token>
   ```

4. **User Tests**
   ```http
   # View Books (Should succeed)
   GET http://localhost:8080/books
   Authorization: Bearer <user_token>
   
   # Create Order (Should succeed)
   POST http://localhost:8080/orders
   Authorization: Bearer <user_token>
   Content-Type: application/json
   
   {
       "customerID": 1,
       "items": [
           {
               "bookID": 1,
               "quantity": 2
           }
       ]
   }
   ```

### Expected Test Results

| Operation           | Admin | Manager | Employee | User |
|--------------------|-------|---------|----------|------|
| View Books         | ✅    | ✅      | ✅       | ✅   |
| Create/Edit Books  | ✅    | ✅      | ❌       | ❌   |
| Manage Customers   | ✅    | ✅      | ✅       | ❌   |
| View Reports       | ✅    | ✅      | ❌       | ❌   |
| Create Orders      | ✅    | ✅      | ✅       | ✅   |
| Delete Users       | ✅    | ❌      | ❌       | ❌   |

### Common HTTP Status Codes
- 200: Success
- 201: Created
- 401: Unauthorized (no/invalid token)
- 403: Forbidden (insufficient permissions)
- 404: Not Found
- 409: Conflict (e.g., duplicate email)
- 500: Internal Server Error

---

## API Documentation

The full API documentation is available in the `swagger.json` file. Use tools like [Swagger UI](https://swagger.io/tools/swagger-ui/) to visualize and interact with the API.

### Key Endpoints
#### Books
- `GET /books`
- `POST /books`
- `GET /books/{id}`
- `PUT /books/{id}`
- `DELETE /books/{id}`

#### Authors
- `GET /authors`
- `POST /authors`
- `GET /authors/{id}`
- `PUT /authors/{id}`
- `DELETE /authors/{id}`

#### Sales Reports
- `GET /sales-reports`

---

## Contributing

We welcome contributions to improve MyBiblio. To contribute:
1. Fork the repository.
2. Create a new branch (`feature/your-feature-name`).
3. Commit your changes.
4. Push your branch and create a Pull Request.

Refer to [GitHub's guide on contributing](https://docs.github.com/en/get-started/quickstart/contributing-to-projects) for best practices.

---

## License

This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/). Feel free to use and modify it.

---

## Acknowledgments

Thanks to all collaborators and contributors for their hard work and dedication to this project.
