
# 🅿️ **PenbunAPI v1.7 Aplha**

PenbunAPI is a RESTful API designed to manage the distribution and supply of books and stationery. It provides robust features for inventory management, order processing, and user authentication using JWT.

## 📘 **Development Standards and Tasks**

To ensure consistent, maintainable development, PenbunAPI is now governed by:

- [`STANDARD.md`](./STANDARD.md): Defines conventions, query patterns, response formats, transaction management, naming, validation rules, and controller templates.
- [`TODO.md`](./TODO.md): Task tracker that manages priorities and work items in checklist format.

## 📦 **New in v1.7.5**
- ✅ เพิ่ม **Product Type API** พร้อมฟังก์ชัน 8 แบบ (Select All, Page, By ID, By Name, Insert, Update, Delete, Remove)
- ✅ รองรับการค้นหาแบบ LIKE โดยใช้ `select/name/:name`
- ✅ ทุกฟังก์ชันใช้ `executeTransaction()` และ `models.ApiResponse`
- ✅ เพิ่ม `unit_type_id`, `unit_type_name`, `description`, `update_by`, `update_date`, `id_status` และ `is_delete` ตามมาตรฐาน Master Table

## 🚀 **Features**

- **Authentication**: รองรับ JWT-based Authentication สำหรับ API ที่ต้องการความปลอดภัย
- **Publisher Management**: เพิ่มฟังก์ชันครบถ้วนสำหรับจัดการข้อมูล Publisher
- **Publisher Type Management**: เพิ่มฟังก์ชันครบถ้วนสำหรับจัดการข้อมูล Publisher Type
- **Customer Management**: เพิ่มฟังก์ชันครบถ้วนสำหรับจัดการข้อมูล Customer
- **Customer Type Management**: เพิ่มฟังก์ชันครบถ้วนสำหรับจัดการข้อมูล Customer Type
- **Book Management**: เพิ่มฟังก์ชันครบถ้วนสำหรับจัดการข้อมูล Customer
- **Book Type Management**: เพิ่มฟังก์ชันครบถ้วนสำหรับจัดการข้อมูล Customer Type
- **Vendor Type Management**: เพิ่มฟังก์ชันครบถ้วนสำหรับจัดการข้อมูล Vendor Type
- **Unit Type Management**: เพิ่มฟังก์ชันครบถ้วนสำหรับจัดการข้อมูล Unit Type
- **Paging**: รองรับการดึงข้อมูลแบบแบ่งหน้า (Pagination)
- **Logging**: จัดการบันทึกข้อมูล Log สำหรับ Audit
- **Versioned**: API (v1, v2)
- **Graceful Shutdown**

## ⚙️ **Fundamental Functions**

> ฟังก์ชันพื้นฐานที่ PenbunAPI ทุก Master Data จะต้องมี ครบ 7 Function โดยโครงสร้างจะทำงานและมีลักษณะเหมือนกันทั้งหมด เพื่อให้การพัฒนาง่ายต่อการดูแลและขยายในอนาคต

| #  | Function         | Description                                                   |
|----|-----------------|---------------------------------------------------------------|
| 1  | Select All       | ดึงข้อมูลทั้งหมด โดย where `is_delete = 0`                  |
| 2  | Select By Paging | รองรับ Query Parameter `?page=<number>&limit=<number>` เพื่อแบ่งหน้า |
| 3  | Select By ID     | ดึงข้อมูลตาม Primary Key เช่น `customer_code` หรือ `publisher_code` หรือ `type_id` |
| 4  | Select By NAME   | ดึงข้อมูลตาม โดยใช้ชื่อ เช่น Select By Name (LIKE `%name%`) |
| 5  | Insert           | เพิ่มข้อมูลใหม่ โดย Insert เฉพาะ field ที่จำเป็น         |
| 6  | Update By ID     | แก้ไขข้อมูลตาม ID โดยไม่แก้ไขค่า Auto Generate เช่น Code ต่าง ๆ |
| 7  | Delete By ID     | Soft Delete โดย Update `is_delete = 1` เท่านั้น             |
| 8  | Remove By ID     | Hard Delete การลบข้อมูลออกจาก Database จริง ๆ              |

- Every Insert/Update/Delete uses centralized `executeTransaction()` from `utils/transaction.go`
- All responses wrapped with `models.ApiResponse`

---

> หมายเหตุเพิ่มเติม:
- ทุกฟังก์ชันที่เกี่ยวข้องกับ Insert / Update / Delete จะมีการใช้ Transaction (Rollback / Panic) ป้องกันข้อมูลไม่ให้เสียหายหากเกิด Error
- ทุก Select จะต้องเช็ค `is_delete = 0` เสมอ
- ฟังก์ชัน Select By Paging จะใช้ Query Parameters:
```
?page=1&limit=20
```
ตัวอย่าง Route:
```
/api/v1/protected/publishers/select/page
/api/v1/protected/customertype/select/page
```
- ฟังก์ชัน Select By NAME ยังไม่มีการ Implement

## ↩️ **Previous Version** (1.7.2)
- **Authentication**: Secure login with JWT-based authentication.
- **Inventory Management**: CRUD operations for books, book types, and references.
- **Order and Delivery**: Manage orders and deliveries.
- **Logging**: Transaction logging for audit purposes.
- **Versioned API**: Support for multiple API versions (e.g., v1, v2).
- **Graceful Shutdown**: Handles safe server shutdown for cleanup and database disconnections.
- **Publisher Management**: เพิ่มฟังก์ชันครบถ้วนสำหรับจัดการข้อมูล Publisher
- **Customer Management**: เพิ่มฟังก์ชันครบถ้วนสำหรับจัดการข้อมูล Customer
- **Discount Management**: เพิ่มฟังก์ชันครบถ้วนสำหรับจัดการข้อมูล Discount
- **Select By Name**: เพิ่มฟังก์ชันสำหรับทุกโมดูล แต่ยังไม่มีการ implement!!
- ** ใช้ `executeTransaction` จาก `utils/transaction.go`
- ** ใช้ `models.ApiResponse` เป็นมาตรฐานการตอบกลับ  
- ** ปรับรูปแบบ Response ให้เป็น `models.ApiResponse` แบบมี key ทุกจุด
- ** เพิ่ม Book API (`tb_book`) พร้อม 8 ฟังก์ชัน
- ** รองรับ LIKE Search ใน `Publisher`, `Book`, `Customer`, `Type` ทุกประเภท
- ** Every Insert/Update/Delete uses centralized `executeTransaction()` from `utils/transaction.go`
- ** All responses wrapped with `models.ApiResponse`
- ** `Publisher` API now returns `publisher_type_name` via JOIN for better usability
- ** Updated all models to use `*string` for nullable fields
- ** Bug fixes and code cleanups for consistency
- ** JOIN type_name in Publisher for display name better than type id 
- ** เพิ่ม **Vendor Type API** พร้อมฟังก์ชัน 8 แบบ (Select All, Page, By ID, By Name, Insert, Update, Delete, Remove)
- ** ทุกฟังก์ชันของ Vendor Type ใช้ `executeTransaction()` และ `models.ApiResponse`
- ** ปรับปรุงการ Query Paging ให้รองรับ SQL Server (OFFSET + FETCH NEXT)
- ** ปรับปรุง Naming Function, Error Handling, และโครงสร้าง Response ให้สอดคล้องกับมาตรฐานล่าสุด

---

## 🧩 **Project Structure**

```
PenbunAPI/
├── main.go
├── config/
│   ├── database.go       # Database connection setup
│   ├── blacklist.go      # Token blacklist
│   ├── env.go            # Environment variable management
│   └── logger.go         # Log configuration
│
├── controllers/
│   ├── auth.go           # Authentication endpoints
│   ├── books.go          # Book management endpoints
│   ├── publishers.go     # Publisher management endpoints
│   ├── publisherType.go  # Publisher Type management endpoints
│   ├── references.go     # Reference management endpoints
│   ├── customer.go       # Customer management endpoints
│   ├── customerType.go   # Customer Type management endpoints
│   ├── book.go           # Book management endpoints
│   ├── bookType.go       # Book Type management endpoints
│   ├── discountType.go   # Discount Type management endpoints
│   ├── discount.go       # Discount management endpoints
│   ├── unitType.go      # Product Type management endpoints
│   ├── productType.go    # Product Type management endpoints
│   ├── vendorType.go     # Vendor Type management endpoints
│   └── vendor.go         # Vendor management endpoints
│
├── models/
│   ├── user.go           # User-related structs and logic
│   ├── book.go           # Book-related structs and logic
│   ├── bookType.go       # Book Type-related structs and logic
│   ├── publisher.go      # Publisher-related structs and logic
│   ├── publisherType.go  # Publisher Type-related structs and logic
│   ├── references.go     # Reference-related structs and logic
│   ├── book.go           # Book management structs and logic
│   ├── bookType.go       # Book Type management structs and logic
│   ├── discountType.go   # Discount Type management structs and logic
│   ├── discount.go       # Discount management structs and logic
│   ├── unitType.go       # Product Type management endpoints
│   ├── productType.go    # Product Type management endpoints
│   ├── vendorType.go     # Vendor Type management structs and logic
│   └── vendor.go         # Vendor management structs and logic
│
├── routes/
│   ├── public.go         # Public API version routes
│   ├── v1.go             # API version 1 routes and grouping
│   └── v2.go             # API version 2 routes (placeholder)
│
├── middleware/
│   └── jwt.go            # JWT middleware for secure endpoints
│
├── logs/
│   └── transaction.log   # Log file for transactions
│
├── .env                  # Environment variables
│
└── go.mod                # Go module file
```

---

## 🪛 **API Documentation**

API Endpoints
-----------------------

# PenbunAPI v1.7.3

## 📘 BookType API
### Base Path: `/api/v1/protected/booktype`

| Method   | Endpoint                      | Description                                  | Required Headers                 | Body Example |
|----------|-------------------------------|----------------------------------------------|----------------------------------|--------------|
| `POST`   | `/insert`                     | เพิ่มข้อมูลประเภทหนังสือใหม่               | `Authorization: Bearer <token>` | `{ "type_name": "Science", "description": "วิทยาศาสตร์", "update_by": "admin" }` |
| `GET`    | `/select/all`                 | ดึงข้อมูลทั้งหมด (is_delete = 0)           | `Authorization: Bearer <token>` | — |
| `GET`    | `/select/page?page=1&limit=10`| ดึงข้อมูลแบบ Paging                         | `Authorization: Bearer <token>` | — |
| `GET`    | `/select/:id`                 | ดึงข้อมูลตาม ID                             | `Authorization: Bearer <token>` | — |
| `GET`    | `/select/:name`               | ดึงข้อมูลตาม Name (LIKE)                    | `Authorization: Bearer <token>` | — |
| `PUT`    | `/update/:id`                 | อัปเดตข้อมูลประเภทหนังสือ                  | `Authorization: Bearer <token>` | `{ "type_name": "Health Sci", "description": "หมวดสุขภาพ", "update_by": "admin" }` |
| `PUT`    | `/delete/:id`                 | Soft Delete (`is_delete = 1`)                | `Authorization: Bearer <token>` | — |
| `DELETE` | `/remove/:id`                 | ลบข้อมูลจริงจากฐานข้อมูล                   | `Authorization: Bearer <token>` | — |

---

## 📗 Book API
### Base Path: `/api/v1/protected/book`

| Method   | Endpoint                      | Description                                  | Required Headers                 | Body Example |
|----------|-------------------------------|----------------------------------------------|----------------------------------|--------------|
| `POST`   | `/insert`                     | เพิ่มหนังสือใหม่                            | `Authorization: Bearer <token>` | `{ "book_name": "AI for Everyone", "book_type_id": "BKTYP0001", "publisher_code": "PUB001", "book_price": 299.0, "book_discount": 50.0, "update_by": "admin" }` |
| `GET`    | `/select/all`                 | ดึงหนังสือทั้งหมด                          | `Authorization: Bearer <token>` | — |
| `GET`    | `/select/page?page=1&limit=10`| ดึงหนังสือแบบแบ่งหน้า                      | `Authorization: Bearer <token>` | — |
| `GET`    | `/select/:id`                 | ดึงหนังสือตาม book_code                     | `Authorization: Bearer <token>` | — |
| `GET`    | `/select/:name`               | ดึงหนังสือโดยใช้ LIKE จาก book_name        | `Authorization: Bearer <token>` | — |
| `PUT`    | `/update/:id`                 | อัปเดตข้อมูลหนังสือ                         | `Authorization: Bearer <token>` | `{ "book_name": "AI Revised", "book_discount": 80.0, "update_by": "editor" }` |
| `PUT`    | `/delete/:id`                 | Soft Delete (`is_delete = 1`)                | `Authorization: Bearer <token>` | — |
| `DELETE` | `/remove/:id`                 | ลบข้อมูลจริงจากฐานข้อมูล                   | `Authorization: Bearer <token>` | — |

---

### 👨‍👩‍👧‍👧 Customer API 
### Base Path: (`/api/v1/protected/customer`)

| Method   | Endpoint                     | Description                                 | Required Headers                  | Body Example |
|----------|--------------------------------|---------------------------------------------|----------------------------------|--------------|
| POST     | `/insert`                     | เพิ่ม Customer ใหม่                        | `Authorization: Bearer <Token>`  | `{ "customer_name": "Siam Bookstore", "biz_id": "BIZ001", "customer_type_id": "CUTMT0001", "first_name": "Somchai", "last_name": "Jaidee", "phone1": "0999999999", "update_by": "admin" }` |
| GET      | `/select/all`                 | ดึงข้อมูล Customer ทั้งหมด               | `Authorization: Bearer <Token>`  | N/A          |
| GET      | `/select/page`                | ดึงข้อมูล Customer แบบ Paging             | `Authorization: Bearer <Token>`  | Query: `?page=1&limit=20` |
| GET      | `/select/:id`                 | ดึงข้อมูล Customer ตาม customer_code     | `Authorization: Bearer <Token>`  | N/A          |
| PUT      | `/update/:id`                 | อัพเดท Customer ตาม customer_code        | `Authorization: Bearer <Token>`  | `{ "customer_name": "Siam Bookstore Updated", "first_name": "Somchai", "last_name": "Jaidee", "update_by": "admin" }` |
| PUT      | `/delete/:id`                 | Soft Delete เปลี่ยน `is_delete = 1`      | `Authorization: Bearer <Token>`  | N/A          |
| DELETE   | `/remove/:id`                 | Hard Delete ลบข้อมูลจริงออกจาก Database  | `Authorization: Bearer <Token>`  | N/A          |

### 🕺 Customer Type API 
### Base Path: (`/api/v1/protected/customertype`)

| Method   | Endpoint                  | Description                                | Required Headers           | Body Example |
|----------|---------------------------|--------------------------------------------|----------------------------|-------------------------------------------------------------------------------------------------------|
| POST   | `/insert`                 | เพิ่ม Customer Type                | `Authorization: Bearer <Token>` | { "type_name": "Wholesale", "description": "Sell for dealer", "update_by": "admin" } |
| GET    | `/select/all`             | ดึงข้อมูลทั้งหมด                  | `Authorization: Bearer <Token>` | - |
| GET    | `/select/page`            | ดึงแบบ Paging                     | `Authorization: Bearer <Token>` | - (Parameter ?page=1&limit=20) |
| GET    | `/select/:id`             | ดึงจาก customer_type_id           | `Authorization: Bearer <Token>` | - |
| PUT    | `/update/:id`             | แก้ไข Customer Type                | `Authorization: Bearer <Token>` | { "type_name": "Retail", "description": "Normal retail customer", "update_by": "admin" } |
| PUT    | `/delete/:id`             | Soft Delete (is_delete = 1)        | `Authorization: Bearer <Token>` | - |
| DELETE | `/remove/:id`             | ลบข้อมูลจริง                      | `Authorization: Bearer <Token>` | - |

### 🔖 Publisher API
### Base Path: (`/api/v1/protected/publishers`)

| Method   | Endpoint                  | Description                                | Required Headers           | Body Example                                                                                           |
|----------|---------------------------|--------------------------------------------|----------------------------|-------------------------------------------------------------------------------------------------------|
| `POST`   | `/publishers/insert`      | Insert a new Publisher                    | `Authorization: Bearer <Token>` | `{ "publisher_name": "Publisher Name", "publisher_type_id": "PUBT001", "contact_name1": "John Doe", ... }` |
| `GET`    | `/publishers/select/all`  | Select all Publishers                     | `Authorization: Bearer <Token>` | N/A                                                                                                   |
| `GET`    | `/publishers/select/page` | Select Publishers with Paging             | `Authorization: Bearer <Token>` | Query: `page` (int), `limit` (int)                                                                    |
| `GET`    | `/publishers/select/:id`  | Select a Publisher by ID                  | `Authorization: Bearer <Token>` | N/A                                                                                                   |
| `PUT`    | `/publishers/update/:id`  | Update a Publisher by ID                  | `Authorization: Bearer <Token>` | `{ "publisher_name": "Updated Name", "contact_name1": "Jane Doe", ... }`                             |
| `PUT`    | `/publishers/delete/:id`  | Soft delete a Publisher (`is_delete = 1`) | `Authorization: Bearer <Token>` | N/A                                                                                                   |
| `DELETE` | `/publishers/remove/:id`  | Hard delete a Publisher                   | `Authorization: Bearer <Token>` | N/A                                                                                                   |

### 📙 Publisher Type API
### Base Path: (`/api/v1/protected/publishertype`)

| Method   | Endpoint                      | Description                                | Required Headers           | Body Example                                                                                           |
|----------|-------------------------------|--------------------------------------------|----------------------------|-------------------------------------------------------------------------------------------------------|
| `POST`   | `/publishertype/insert`       | Insert a new Publisher Type               | `Authorization: Bearer <Token>` | `{ "type_name": "Bookstore", "description": "Retail bookstore type", ... }`                          |
| `GET`    | `/publishertype/select/all`   | Select all Publisher Types                | `Authorization: Bearer <Token>` | N/A                                                                                                   |
| `GET`    | `/publishertype/select/page`  | Select Publisher Types with Paging        | `Authorization: Bearer <Token>` | Query: `page` (int), `limit` (int)                                                                    |
| `GET`    | `/publishertype/select/:id`   | Select a Publisher Type by ID             | `Authorization: Bearer <Token>` | N/A                                                                                                   |
| `PUT`    | `/publishertype/update/:id`   | Update a Publisher Type by ID             | `Authorization: Bearer <Token>` | `{ "type_name": "Wholesale", "description": "Wholesale distributor type", ... }`                     |
| `PUT`    | `/publishertype/delete/:id`   | Soft delete a Publisher Type (`is_delete = 1`) | `Authorization: Bearer <Token>` | N/A                                                                                                   |
| `DELETE` | `/publishertype/remove/:id`   | Hard delete a Publisher Type              | `Authorization: Bearer <Token>` | N/A                                                                                                   |

### 💸 Discount Type API  
### Base Path: (`/api/v1/protected/discounttype`)

| Method   | Endpoint                          | Description                                   | Required Headers                | Body Example |
|----------|-----------------------------------|-----------------------------------------------|----------------------------------|--------------|
| `POST`   | `/discounttype/insert`            | Insert a new Discount Type                   | `Authorization: Bearer <Token>` | `{ "type_name": "Summer Sale", "discount_unit_type": "percent", "update_by": "admin" }` |
| `GET`    | `/discounttype/select/all`        | Select all Discount Types                    | `Authorization: Bearer <Token>` | — |
| `GET`    | `/discounttype/select/page`       | Select Discount Types with Paging            | `Authorization: Bearer <Token>` | Query: `page=1&limit=20` |
| `GET`    | `/discounttype/select/:id`        | Select a Discount Type by ID                 | `Authorization: Bearer <Token>` | — |
| `PUT`    | `/discounttype/update/:id`        | Update a Discount Type by ID                 | `Authorization: Bearer <Token>` | `{ "type_name": "Holiday Promo", "discount_unit_type": "fixed", "update_by": "admin" }` |
| `PUT`    | `/discounttype/delete/:id`        | Soft delete a Discount Type (`is_delete = 1`) | `Authorization: Bearer <Token>` | — |
| `DELETE` | `/discounttype/remove/:id`        | Hard delete a Discount Type                  | `Authorization: Bearer <Token>` | — |
---

## 💽 **Libraries and Frameworks**

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

## 💾 **Installation and Setup**

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

## ©️ **License**

This project is licensed under the PENBUN License. See the LICENSE file for details.
