
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

### Basic Endpoints
#### Books
- **Create a Book**
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{
      "title": "The Go Programming Language",
      "author": { "id": 1, "first_name": "Alan", "last_name": "Donovan" },
      "genres": ["Programming"],
      "price": 45.99,
      "stock": 100
  }' http://localhost:8080/books
  ```

- **Get All Books**
  ```bash
  curl http://localhost:8080/books
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

---

## License

This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/). Feel free to use and modify it.

---

## Acknowledgments

Thanks to all collaborators and contributors for their hard work and dedication to this project.
