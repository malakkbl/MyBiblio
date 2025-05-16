# MyBiblio - Online Bookstore API

## Description

**MyBiblio** is an educational project implementing a RESTful API for online bookstore management.

### Technical Overview
The project is built with Go and uses:
- PostgreSQL database with GORM for data persistence
- JWT-based authentication for user security
- Role-based access control for different user types
- Request validation and error handling
- RESTful endpoints for resource management

### Learning Objectives
This project focuses on understanding:
- Backend API development principles
- Database design and ORM usage
- Authentication and authorization concepts
- Input validation and error handling patterns
- API documentation practices

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
    - Total books sold.    - Top-selling books.

### Development Progress

#### Current Implementation Status
| Feature          | Implementation | Input Validation | Documentation |
|------------------|---------------------|------------------|---------------|
| Authentication   | ✅                  | ✅               | ✅            |
| Books           | ✅                  | ✅               | ✅            |
| Authors         | ✅                  | ✅               | ✅            |
| Customers       | ✅                  | ✅               | ✅            |
| Orders          | ✅                  | ✅               | ✅            |
| Sales Reports   | ✅                  | ✅               | ✅            |

#### Key Features Implemented
- **Authentication**: Basic JWT-based user authentication
- **Input Validation**: Request validation with helpful error messages
- **Error Handling**: Standardized error responses
- **API Documentation**: Basic Swagger documentation available

### Project Milestones & Status

The project follows a structured development roadmap. Below is the current status of each milestone:

#### **Completed Milestones:**
✅ **Milestone 1: Database Implementation**
- Implemented PostgreSQL with GORM
- Set up database models and relationships
- Established connection handling and configuration

✅ **Milestone 2: Authentication & Authorization**
- Implemented JWT-based security
- Role-based access control (RBAC)
- Secure password handling
- Permission-based endpoints

✅ **Milestone 3: Input Validation & Error Handling**
- Comprehensive request validation
- Custom email format validation
- Password strength requirements
- Standardized error responses
- Detailed validation feedback

#### **Pending Milestones:**
🔲 **Milestone 4:** Caching & Comprehensive Testing
- Redis caching implementation
- Unit tests coverage
- Integration tests
- End-to-end testing
- Performance testing

🔲 **Milestone 5:** Containerization & CI/CD
- Docker containerization
- CI/CD pipeline setup
- Monitoring implementation
- Logging system

🔲 **Milestone 6:** Database Migrations
- Flyway integration
- Version-controlled schemas
- Automated migrations

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

## Role-Based Access Control (RBAC) Testing Guide

MyBiblio implements a secure JWT-based authentication system with role-based access control (RBAC). The system includes comprehensive input validation and error handling.

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

#### Authentication & Error Handling

1. **Register User**
   ```http
   POST /api/auth/register
   ```
   **Required Fields:**
   - `name`: String (2-100 characters)
   - `email`: String (Valid email format)
     - Must be 3-64 characters before @
     - Domain must be 2-255 characters
     - No consecutive dots
     - Only allowed special characters (!#$%&'*+-/=?^_`{|}~.)
   - `password`: String with requirements:
     - Minimum 8 characters     - At least one uppercase letter
     - At least one lowercase letter
     - At least one number
     - At least one special character
   - `role`: String (one of: admin, manager, employee, user)

2. **Login**
   ```http
   POST /api/auth/login
   ```
   **Required Fields:**
   - `email`: String (same validation as register)
   - `password`: String
   
   **Returns:**
   ```json
   {
     "token": "JWT_TOKEN",
     "user": {
       "id": 1,
       "name": "User Name",
       "email": "user@example.com",
       "role": "user"
     },
     "permissions": ["read:books", "write:orders"]
   }
   ```

### Error Handling

The API implements comprehensive error handling with detailed feedback:

```json
{
  "code": "ERROR_CODE",
  "message": "User-friendly error message",
  "details": [
    {
      "field": "field_name",
      "tag": "validation_tag",
      "value": "invalid_value",
      "message": "Detailed error message"
    }
  ],
  "debug": "Additional debug information (development only)"
}
```

#### Common Error Codes

1. **Authentication Errors**
   - `INVALID_CREDENTIALS`: Invalid email or password
   - `INVALID_TOKEN`: Invalid or malformed token
   - `EXPIRED_TOKEN`: Token has expired
   - `MISSING_TOKEN`: No token provided
   - `WEAK_PASSWORD`: Password requirements not met
   - `INVALID_ROLE`: Role must be one of: admin, manager, employee, user
   - `INVALID_EMAIL`: Email format validation failed

2. **Validation Errors**
   - `VALIDATION_ERROR`: Input validation failed
   - `INVALID_INPUT`: Invalid request format or data
   - `DUPLICATE_ENTRY`: Resource already exists

3. **Authorization Errors**
   - `UNAUTHORIZED`: Authentication required
   - `FORBIDDEN`: Insufficient permissions

4. **Database Errors**
   - `DATABASE_ERROR`: Database operation failed
   - `NOT_FOUND`: Requested resource not found

#### Error Response Examples

1. **Invalid Credentials**
   ```json
   {
     "code": "INVALID_CREDENTIALS",
     "message": "Invalid email or password"
   }
   ```

2. **Email Validation Error**
   ```json
   {
     "code": "VALIDATION_ERROR",
     "message": "Validation failed",
     "details": [
       {
         "field": "email",
         "tag": "custom_email",
         "value": "invalid@email",
         "message": "email must be a valid email address between 3-64 characters before @ and 2-255 characters after @, containing only allowed special characters"
       }
     ]
   }
   ```

3. **Password Validation Error**
   ```json
   {
     "code": "VALIDATION_ERROR",
     "message": "Validation failed",
     "details": [
       {
         "field": "password",
         "tag": "passwd",
         "message": "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character"
       }
     ]
   }
   ```

3. **Authorization Error**
   ```json
   {
     "code": "FORBIDDEN",
     "message": "Access denied. Required roles: admin, manager"
   }
   ```

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

## License

MIT License

Copyright (c) 2025 Kably Malak

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish,distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.


---

Thank you for checking out MyBiblio! I hope you find it useful for managing your online bookstore. If you have any questions or feedback, feel free to reach out.