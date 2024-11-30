
# PenbunAPI

PenbunAPI is a RESTful API designed to manage the distribution and supply of books and stationery. It provides robust features for inventory management, order processing, and user authentication using JWT.

## Features

- **Authentication**: Secure login with JWT-based authentication.
- **Inventory Management**: CRUD operations for books, book types, and references.
- **Order and Delivery**: Manage orders and deliveries.
- **Logging**: Transaction logging for audit purposes.
- **Versioned API**: Support for multiple API versions (e.g., v1, v2).
- **Graceful Shutdown**: Handles safe server shutdown for cleanup and database disconnections.

## Project Structure

```
PenbunAPI/
├── main.go
├── config/
│   ├── database.go       # Database connection setup
│   ├── blacklist.go      # Token blacklist
│   ├── env.go            # Environment variable management
│   └── logger.go         # Log configuration
├── controllers/
│   ├── auth.go           # Authentication endpoints
│   ├── books.go          # Book management endpoints
│   └── references.go     # Reference management endpoints
├── models/
│   ├── user.go           # User-related structs and logic
│   ├── book.go           # Book-related structs and logic
│   ├── bookType.go       # Book Type management endpoints
│   └── reference.go      # Reference management endpoints
├── routes/
│   ├── public.go         # Public API version routes for testing or pinging
│   ├── v1.go             # API version 1 routes and separate Public/Protected 
│   └── v2.go             # API version 2 routes and separate Public/Protected 
├── middleware/
│   └── jwt.go            # JWT middleware for secure endpoints
├── logs/
│   └── transaction.log   # Log file for transactions
├── .env                  # Environment variables
└── go.mod                # Go module file
```

## Libraries and Frameworks

### **Backend Framework**
- [Fiber](https://gofiber.io/) - High-performance web framework for Go.

### **Authentication**
- [JWT (golang-jwt)](https://github.com/golang-jwt/jwt) - JWT implementation in Go for secure authentication.

### **Database**
- [MSSQL (go-mssqldb)](https://github.com/denisenkom/go-mssqldb) - Microsoft SQL Server driver for Go.

### **Hashing**
- [Bcrypt (golang.org/x/crypto/bcrypt)](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Secure password hashing.

### **Environment Variables**
- [Godotenv](https://github.com/joho/godotenv) - Load environment variables from `.env` file.

### **Logging**
- [Logrus](https://github.com/sirupsen/logrus) - Structured logging for Go.

## Installation and Setup

### Prerequisites
- Go (1.19 or higher)
- Microsoft SQL Server
- Git (optional, for cloning the repository)

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/PenbunAPI.git
   cd PenbunAPI
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Configure the `.env` file:
   ```
   DB_HOST=your_db_host
   DB_PORT=1433
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=your_db_name
   JWT_SECRET=your_jwt_secret
   LOG_FILE=logs/transaction.log
   ```

4. Run the server:
   ```bash
   go run main.go
   ```

5. **Optional**: Create a new user using `bcrypt` for password hashing.

   - Install `htpasswd`:
     ```bash
     sudo apt update
     sudo apt install apache2-utils -y
     ```

   - Generate `bcrypt` hash:
     ```bash
     htpasswd -nbBC 10 username password
     ```

     Sample output:
     ```
     username:$2y$10$KfQ8mU5VvJ5QGk7/LN9OeOujOPEwLjD3Oo4yEWDwEpr6/LkfuPWoK
     ```

   - Insert into Database:
     ```sql
     DELETE FROM tb_users;
     DBCC CHECKIDENT ('tb_users', RESEED, 0);

     INSERT INTO tb_users (user_name, user_password)
     VALUES ('username', '$2y$10$KfQ8mU5VvJ5QGk7/LN9OeOujOPEwLjD3Oo4yEWDwEpr6/LkfuPWoK');
     ```

## Endpoints Overview

### Version 1 (`/api/v1`)
- **Authentication**:
  - `POST /public/login` - User login and JWT generation.
  - `POST /public/logout` - Logout and blacklist token.
- **Books**:
  - `POST /protected/books` - Create a new book.
  - `PUT /protected/books/:id` - Update book details.
  - `DELETE /protected/books/:id` - Delete a book.
- **References**:
  - `GET /protected/reference` - Retrieve reference data.
- **Ping**:
  - `GET /public/ping` - Check server health.

### Version 2 (`/api/v2`)
- **Books**:
  - `GET /public/books` - Retrieve all books.
  - `GET /protected/books/:id` - Retrieve book by ID.

## Graceful Shutdown

- Fiber supports Graceful Shutdown for safe cleanup and resource disconnection.
- Handled in `main.go` using Unix Signals (`SIGINT`, `SIGTERM`).

## License

This project is licensed under the MIT License. See the LICENSE file for details.
