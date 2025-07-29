# ✅ Todo.md – PenbunAPI Development Tracker

> ใช้เพื่อติดตามสถานะการพัฒนา API และ Module ต่าง ๆ ในระบบ PenbunAPI  
> อัปเดตตาม version ล่าสุด (v1.7.4)

---

## 🧩 Master Data Module

- ~~[x] เพิ่ม Publisher API (8 ฟังก์ชัน)~~
- ~~[x] เพิ่ม Publisher Type API (8 ฟังก์ชัน)~~
- [x] เพิ่ม Customer API (8 ฟังก์ชัน)
- [x] เพิ่ม Customer Type API (8 ฟังก์ชัน)
- [x] เพิ่ม Vendor Type API (8 ฟังก์ชัน)
- [x] เพิ่ม Vendor API (8 ฟังก์ชัน)
- [x] เพิ่ม Book API (8 ฟังก์ชัน)
- [x] เพิ่ม Book Type API (8 ฟังก์ชัน)
- [x] เพิ่ม Discount Type API (8 ฟังก์ชัน)
- [x] เพิ่ม Discount API (8 ฟังก์ชัน)
- [x] เพิ่ม Unit Type API (8 ฟังก์ชัน) 
- [ ] เพิ่ม Product Type API (8 ฟังก์ชัน) ← 🛠️ กำลังทำ
- [ ] เพิ่ม Product Format Type API (8 ฟังก์ชัน)
- [ ] เพิ่ม Product API (8 ฟังก์ชัน)

---

## 📦 Receive Module

- [ ] ออกแบบ API: `tb_product_lot`
- [ ] ออกแบบ API: `tb_product_lot_transaction`

---

## 🧾 Order Module

- [ ] ออกแบบ API: รับคำสั่งซื้อ
- [ ] ตรวจสอบ stock ก่อนยืนยันคำสั่งซื้อ

---

## 🚚 Deliver Module

- [ ] ออกแบบ API: สร้างใบจัดส่งสินค้า
- [ ] อัปเดตสถานะการจัดส่ง

---

## 🔁 Return Module

- [ ] ออกแบบ API: คืนสินค้า
- [ ] ตรวจสอบสภาพสินค้า + คืนเงิน/เปลี่ยนสินค้า

---

## 🧾 Invoice Module

- [ ] ออกแบบ API: ใบแจ้งหนี้
- [ ] ติดตามสถานะการชำระเงิน

---

## ⚙️ System & Infra

- [x] JWT Middleware + Refresh Token
- [x] Token Blacklist (logout)
- [x] Graceful Shutdown
- [ ] เพิ่มระบบ Logging สำหรับ Audit
- [ ] วางโครงสร้าง Role/Permission (ระยะถัดไป)

---

## 🧪 Test Coverage & Validation

- [ ] ทดสอบ API ผ่าน Postman
- [ ] ทดสอบ Paging ทุก Module
- [ ] ทดสอบ Soft Delete/Remove ทุกฟังก์ชัน
- [ ] ทดสอบ Refresh Token Flow

---

## ⏳ Backlog / Idea

- [ ] UI Web Admin สำหรับทดสอบ API
- [ ] Export CSV รายงานจากระบบ
- [ ] ระบบ Cron Job แจ้งเตือน Invoice ค้างชำระ
