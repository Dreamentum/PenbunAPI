
# PenbunAPI v1.4.3

PenbunAPI is a RESTful API designed to manage the distribution and supply of books and stationery. It provides robust features for inventory management, order processing, and user authentication using JWT.

## **Features**

- **Authentication**: รองรับ JWT-based Authentication สำหรับ API ที่ต้องการความปลอดภัย
- **Publisher Management**: เพิ่มฟังก์ชันครบถ้วนสำหรับจัดการข้อมูล Publisher
- **Publisher Type Management**: เพิ่มฟังก์ชันครบถ้วนสำหรับจัดการข้อมูล Publisher Type
- **Paging**: รองรับการดึงข้อมูลแบบแบ่งหน้า (Pagination) สำหรับ Publisher และ Publisher Type
- **Logging**: จัดการบันทึกข้อมูล Log สำหรับ Audit

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
- Added Paging support for retrieving Publisher data.
  - Route: /api/v1/protected/publishers/select/page
  - Query Parameters: ?page=<page_number>&limit=<items_per_page>

## What's New in v1.4.3

### **Publisher API**
1. เพิ่มฟังก์ชันใหม่:
   - Insert Publisher
   - Select All Publishers
   - Select Publisher By ID
   - Select Publishers with Paging
   - Update Publisher By ID
   - Soft Delete Publisher (is_delete)
   - Hard Delete Publisher

2. Routing ใหม่สำหรับ Publisher:
   - `/api/v1/protected/publishers/insert`
   - `/api/v1/protected/publishers/select/all`
   - `/api/v1/protected/publishers/select/page`
   - `/api/v1/protected/publishers/select/:id`
   - `/api/v1/protected/publishers/update/:id`
   - `/api/v1/protected/publishers/delete/:id`
   - `/api/v1/protected/publishers/remove/:id`

---

### **Publisher Type API**
1. เพิ่มฟังก์ชันใหม่:
   - Insert Publisher Type
   - Select All Publisher Types
   - Select Publisher Type By ID
   - Select Publisher Types with Paging
   - Update Publisher Type By ID
   - Soft Delete Publisher Type (is_delete)
   - Hard Delete Publisher Type

2. Routing ใหม่สำหรับ Publisher Type:
   - `/api/v1/protected/publishertype/insert`
   - `/api/v1/protected/publishertype/select/all`
   - `/api/v1/protected/publishertype/select/page`
   - `/api/v1/protected/publishertype/select/:id`
   - `/api/v1/protected/publishertype/update/:id`
   - `/api/v1/protected/publishertype/delete/:id`
   - `/api/v1/protected/publishertype/remove/:id`

---

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
│   ├── publishers.go     # Publisher management endpoints
│   ├── publisherType.go  # Publisher Type management endpoints
│   └── references.go     # Reference management endpoints
├── models/
│   ├── user.go           # User-related structs and logic
│   ├── book.go           # Book-related structs and logic
│   ├── bookType.go       # Book Type-related structs and logic
│   ├── publisher.go      # Publisher-related structs and logic
│   ├── publisherType.go  # Publisher Type-related structs and logic
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

---

## **API Documentation**

API Endpoints
-----------------------

### Base Path: `/api/v1/protected/publishers`

### **Publisher API**

| Method   | Endpoint                  | Description                                | Required Headers           | Body Example                                                                                           |
|----------|---------------------------|--------------------------------------------|----------------------------|-------------------------------------------------------------------------------------------------------|
| `POST`   | `/publishers/insert`      | Insert a new Publisher                    | `Authorization: Bearer <Token>` | `{ "publisher_name": "Publisher Name", "publisher_type_id": "PUBT001", "contact_name1": "John Doe", ... }` |
| `GET`    | `/publishers/select/all`  | Select all Publishers                     | `Authorization: Bearer <Token>` | N/A                                                                                                   |
| `GET`    | `/publishers/select/page` | Select Publishers with Paging             | `Authorization: Bearer <Token>` | Query: `page` (int), `limit` (int)                                                                    |
| `GET`    | `/publishers/select/:id`  | Select a Publisher by ID                  | `Authorization: Bearer <Token>` | N/A                                                                                                   |
| `PUT`    | `/publishers/update/:id`  | Update a Publisher by ID                  | `Authorization: Bearer <Token>` | `{ "publisher_name": "Updated Name", "contact_name1": "Jane Doe", ... }`                             |
| `PUT`    | `/publishers/delete/:id`  | Soft delete a Publisher (`is_delete = 1`) | `Authorization: Bearer <Token>` | N/A                                                                                                   |
| `DELETE` | `/publishers/remove/:id`  | Hard delete a Publisher                   | `Authorization: Bearer <Token>` | N/A                                                                                                   |

### Base Path: `/api/v1/protected/publishertype`

### **Publisher Type API**

| Method   | Endpoint                      | Description                                | Required Headers           | Body Example                                                                                           |
|----------|-------------------------------|--------------------------------------------|----------------------------|-------------------------------------------------------------------------------------------------------|
| `POST`   | `/publishertype/insert`       | Insert a new Publisher Type               | `Authorization: Bearer <Token>` | `{ "type_name": "Bookstore", "description": "Retail bookstore type", ... }`                          |
| `GET`    | `/publishertype/select/all`   | Select all Publisher Types                | `Authorization: Bearer <Token>` | N/A                                                                                                   |
| `GET`    | `/publishertype/select/page`  | Select Publisher Types with Paging        | `Authorization: Bearer <Token>` | Query: `page` (int), `limit` (int)                                                                    |
| `GET`    | `/publishertype/select/:id`   | Select a Publisher Type by ID             | `Authorization: Bearer <Token>` | N/A                                                                                                   |
| `PUT`    | `/publishertype/update/:id`   | Update a Publisher Type by ID             | `Authorization: Bearer <Token>` | `{ "type_name": "Wholesale", "description": "Wholesale distributor type", ... }`                     |
| `PUT`    | `/publishertype/delete/:id`   | Soft delete a Publisher Type (`is_delete = 1`) | `Authorization: Bearer <Token>` | N/A                                                                                                   |
| `DELETE` | `/publishertype/remove/:id`   | Hard delete a Publisher Type              | `Authorization: Bearer <Token>` | N/A                                                                                                   |

---

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

5. **Optional**: Create a new user using `bcrypt` for password hashing.
   Install `htpasswd`:
   ```bash
   sudo apt update
   sudo apt install apache2-utils -y
   ```
   Generate `bcrypt` hash:
   ```bash
   htpasswd -nbBC 10 username password
   ```
   Sample output:
   ```
   username:$2y$10$KfQ8mU5VvJ5QGk7/LN9OeOujOPEwLjD3Oo4yEWDwEpr6/LkfuPWoK
   ```
   Insert into Database:
   ```sql
   DELETE FROM tb_users;
   DBCC CHECKIDENT ('tb_users', RESEED, 0);
   INSERT INTO tb_users (user_name, user_password)
   VALUES ('username', '$2y$10$KfQ8mU5VvJ5QGk7/LN9OeOujOPEwLjD3Oo4yEWDwEpr6/LkfuPWoK');
   ```


## License

This project is licensed under the MIT License. See the LICENSE file for details.
