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

- ทุกโมดูลต้องรองรับ `LIKE` search โดยใช้:
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

> 🛠 เวอร์ชันล่าสุดของมาตรฐานนี้: v1.7.4
