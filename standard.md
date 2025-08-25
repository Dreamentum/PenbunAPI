# 📏 PenbunAPI Development Standard

เอกสารนี้กำหนดแนวทางและมาตรฐานในการพัฒนา PenbunAPI อย่างเป็นระบบ เพื่อความสอดคล้อง รวดเร็ว และง่ายต่อการดูแลในระยะยาว

---

## 1. 🧱 Database Table Structure

- ทุกตารางต้องมี:
  - `autoID` (INT, Primary Key)
  - `prefix` (NVARCHAR(5)) สำหรับใช้สร้างรหัสแบบ dynamic
  - `[table]_id` หรือ `[table]_code` (NVARCHAR(50)) ใช้เป็นรหัสหลัก
  - `update_by`, `update_date`, `is_delete` (BIT DEFAULT 0)
- ใช้ Trigger:
  - `TRIG_AUTO_UPDATE_DATE_[TABLE_NAME]`
  - `TRIG_GENERATE_[TABLE]_ID` (ใช้ `prefix` เป็น dynamic ID)

---

## 2. 🧠 API Pattern (8 ฟังก์ชันหลัก)

1. Select All  
2. Select By Paging  
3. Select By ID  
4. Select By Name (LIKE `%name%`)  
5. Insert  
6. Update By ID  
7. Delete By ID (Soft Delete)  
8. Remove By ID (Hard Delete)  

---

## 3. 📦 API Design Guideline

- ใช้ `executeTransaction()` จาก `utils/transaction.go`
- ทุก Response ใช้ `models.ApiResponse` เท่านั้น
- ทุก API Route อยู่ภายใต้ `/api/v1/protected/[module]`

---

## 4. 🔐 Authentication

- ใช้ JWT สำหรับ API ทั้งหมดที่ขึ้นต้นด้วย `/protected`
- Middleware ตรวจสอบ Token จาก `middleware/jwt.go`

---

## 5. 📄 Naming Convention

- ชื่อไฟล์: `camelCase`
- ชื่อตัวแปร/ฟังก์ชัน: `hungarian + camelCase`
- ชื่อ column/table: `snake_case`
- Prefix Code เช่น `B`, `C`, `P`, `V`

---

## 6. 🔄 Transaction Handling

- ทุก Insert, Update, Delete ต้องใช้ `executeTransaction()` เสมอ
- Rollback ถ้าเกิด panic/error เพื่อป้องกันข้อมูลเสียหาย

---

## 7. 🔍 LIKE Search (Select By Name)

```sql
WHERE type_name LIKE '%' + @type_name + '%' AND is_delete = 0
```
- ใช้ parameter `c.Params("name")`

---

## 8. 🧪 Data Validation

- Validate field ที่จำเป็นก่อน Insert/Update
- ตรวจสอบ foreign key ว่ามีอยู่จริง เช่น `publisher_type_id`, `customer_type_id`

---

## 9. 🌏 TimeZone & Logging

- ทุกตารางใช้เวลา `SE Asia Standard Time`
- บันทึก log ที่ `logs/transaction.log`

---

## 🧩 PenbunAPI Controller Template (8 Standard Functions)

Template นี้ใช้เป็นโครงสร้างพื้นฐานของ Controller สำหรับ Entity ทุกตัวในระบบ PenbunAPI

---

### 🔷 1. Select All

```go
func SelectAll<Entity>(c *fiber.Ctx) error {
    query := `SELECT ... FROM tb_<entity> WHERE is_delete = 0`
    rows, err := config.DB.Query(query)
    ...
}
```

---

### 🔷 2. Select By Paging

```go
func SelectPage<Entity>(c *fiber.Ctx) error {
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    offset := (page - 1) * limit

    query := `SELECT ... FROM tb_<entity> WHERE is_delete = 0 ORDER BY update_date DESC OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY`
    ...
}
```

---

### 🔷 3. Select By ID

```go
func Select<Entity>ByID(c *fiber.Ctx) error {
    id := c.Params("id")
    query := `SELECT ... FROM tb_<entity> WHERE <entity>_id = @ID AND is_delete = 0`
    ...
}
```

---

### 🔷 4. Select By Name

```go
func Select<Entity>ByName(c *fiber.Ctx) error {
    name := c.Params("name")
    query := `SELECT ... FROM tb_<entity> WHERE type_name LIKE '%' + @Name + '%' AND is_delete = 0`
    ...
}
```

---

### 🔷 5. Insert

```go
func Insert<Entity>(c *fiber.Ctx) error {
    var item models.<Entity>
    if err := c.BodyParser(&item); err != nil {
        ...
    }

    query := `INSERT INTO tb_<entity> (...) VALUES (...)`
    ...
}
```

---

### 🔷 6. Update By ID

```go
func Update<Entity>ByID(c *fiber.Ctx) error {
    id := c.Params("id")
    var item models.<Entity>
    if err := c.BodyParser(&item); err != nil {
        ...
    }

    query := `UPDATE tb_<entity> SET ... WHERE <entity>_id = @ID AND is_delete = 0`
    ...
}
```

---

### 🔷 7. Delete By ID (Soft Delete)

```go
func Delete<Entity>ByID(c *fiber.Ctx) error {
    id := c.Params("id")
    username := c.Query("user") // รับชื่อผู้ใช้จาก query string เช่น ?user=ROOT
    if username == "" {
        username = "UNKNOWN"
    }

    query := `UPDATE tb_<entity> SET is_delete = 1, update_by = @UpdateBy WHERE <entity>_id = @ID`
```

### 🔷 8. Remove By ID (Hard Delete)

```go
func Remove<Entity>ByID(c *fiber.Ctx) error {
    id := c.Params("id")
    query := `DELETE FROM tb_<entity> WHERE <entity>_id = @ID`
    ...
}
```

---

🛠️ เปลี่ยน `<Entity>` เป็นชื่อ struct และ `<entity>` เป็นชื่อ table เช่น `VendorType`, `vendor_type`

> 🔖 เวอร์ชันล่าสุดของมาตรฐานนี้: v1.9
