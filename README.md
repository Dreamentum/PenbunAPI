
# PenbunAPI v1.4.2

PenbunAPI is a RESTful API designed to manage the distribution and supply of books and stationery. It provides robust features for inventory management, order processing, and user authentication using JWT.

## Previous Version
- **Authentication**: Secure login with JWT-based authentication.
- **Inventory Management**: CRUD operations for books, book types, and references.
- **Order and Delivery**: Manage orders and deliveries.
- **Logging**: Transaction logging for audit purposes.
- **Versioned API**: Support for multiple API versions (e.g., v1, v2).
- **Graceful Shutdown**: Handles safe server shutdown for cleanup and database disconnections.
- **Added comprehensive management for `tb_publisher` including:
  - Insert Publisher
  - Select All Publishers
  - Select Publisher by ID
  - Update Publisher by ID
  - Delete Publisher (soft delete by updating `is_delete = 1`)
  - Remove Publisher (hard delete from the database)
- **Introduced transaction handling (`Rollback`, `Panic`) for critical operations to ensure data consistency.
- **Enhanced API structure for improved modularity and separation of concerns.

## What's New in v1.4.2

- Added Paging support for retrieving Publisher data.
  - Route: /api/v1/protected/publishers/select/page
  - Query Parameters: ?page=<page_number>&limit=<items_per_page>

## Updated Project Structure

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
│   └── publishers.go     # Publisher management endpoints
│   └── references.go     # Reference management endpoints
├── models/
│   ├── user.go           # User-related structs and logic
│   ├── book.go           # Book-related structs and logic
│   ├── bookType.go       # Book Type management endpoints
│   └── publisher.go      # Publisher-related structs and logic
│   └── reference.go      # Reference-related structs and logic
├── routes/
│   ├── public.go         # Public API version routes
│   ├── v1.go             # API version 1 routes and grouping
│   └── v2.go             # API version 2 routes (placeholder)
├── middleware/
│   └── jwt.go            # JWT middleware for secure endpoints
├── logs/
│   └── transaction.log   # Log file for transactions
├── .env                  # Environment variables
└── go.mod                # Go module file
```

## Publisher API Endpoints

### Base Path: `/api/v1/protected/publishers`

Publisher API Endpoints
-----------------------

### Base Path: `/api/v1/protected/publishers`

| Method | Endpoint | Description | Required Headers | Body Example |
| --- | --- | --- | --- | --- |
| `GET` | `/select/all` | Retrieve all Publishers | `Authorization: Bearer <Token>` | N/A |
| `GET` | `/select/page` | Retrieve Publishers with Paging | `Authorization: Bearer <Token>` | Query Parameters: `?page=<page_number>&limit=<items_per_page>` |
| `GET` | `/select/:id` | Retrieve a Publisher by ID | `Authorization: Bearer <Token>` | N/A |
| `POST` | `/insert` | Add a new Publisher | `Authorization: Bearer <Token>` | `{ "publisher_type_id": "PUBT001", "publisher_name": "Publisher Name", "contact_name1": "John Doe", ... }` |
| `PUT` | `/update/:id` | Update a Publisher by ID | `Authorization: Bearer <Token>` | `{ "publisher_name": "Updated Name", "contact_name1": "Jane Doe", ... }` |
| `PUT` | `/delete/:id` | Soft delete a Publisher (`is_delete = 1`) | `Authorization: Bearer <Token>` | N/A |
| `DELETE` | `/remove/:id` | Hard delete a Publisher | `Authorization: Bearer <Token>` | N/A |

## All Example in Postman

| Method | Endpoint | Description | Header | Body (Example) |
| --- | --- | --- | --- | --- |
| `GET` | `/api/v1/protected/publishers/select/all` | ดึงข้อมูล Publisher ทั้งหมด | `Authorization: Bearer <Token>` | N/A |
| `GET` | `/api/v1/protected/publishers/select/page` | ดึงข้อมูล Publisher แบบ Paging | `Authorization: Bearer <Token>` | Query Parameters: `?page=1&limit=10` |
| `GET` | `/api/v1/protected/publishers/select/:id` | ดึงข้อมูล Publisher ตาม ID | `Authorization: Bearer <Token>` | N/A |
| `POST` | `/api/v1/protected/publishers/insert` | เพิ่ม Publisher ใหม่ | `Authorization: Bearer <Token>` | `{ "publisher_type_id": "PUBT001", "publisher_name": "Publisher Name", "contact_name1": "John Doe", ... }` |
| `PUT` | `/api/v1/protected/publishers/update/:id` | อัปเดตข้อมูล Publisher ตาม ID | `Authorization: Bearer <Token>` | `{ "publisher_name": "Updated Name", "contact_name1": "Jane Doe", ... }` |
| `PUT` | `/api/v1/protected/publishers/delete/:id` | เปลี่ยน is_delete = 1 (Soft Delete) | `Authorization: Bearer <Token>` | N/A |
| `DELETE` | `/api/v1/protected/publishers/remove/:id` | ลบข้อมูล Publisher ออกจากฐานข้อมูลจริง (Hard Delete) | `Authorization: Bearer <Token>` | N/A |

## Libraries and Frameworks

### Backend Framework
- [Fiber](https://gofiber.io/) - High-performance web framework for Go.

### Authentication
- [JWT (golang-jwt)](https://github.com/golang-jwt/jwt) - JWT implementation in Go for secure authentication.

### Database
- [MSSQL (go-mssqldb)](https://github.com/denisenkom/go-mssqldb) - Microsoft SQL Server driver for Go.

### Hashing
- [Bcrypt (golang.org/x/crypto/bcrypt)](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Secure password hashing.

### Environment Variables
- [Godotenv](https://github.com/joho/godotenv) - Load environment variables from `.env` file.

### Logging
- Built-in `log` package in Go for lightweight logging.

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

## License

This project is licensed under the MIT License. See the LICENSE file for details.
