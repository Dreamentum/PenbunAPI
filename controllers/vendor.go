package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"
	"math"

	"github.com/gofiber/fiber/v2"
)

// SelectAllVendor: ดึงข้อมูล Vendor ทั้งหมด (is_delete = 0)
func SelectAllVendor(c *fiber.Ctx) error {
	query := `
		SELECT v.vendor_id, v.vendor_name, v.vendor_type_id, t.type_name, 
			v.discount_id, d.discount_name,
			v.contact_name1, v.contact_name2, v.email, v.phone1, v.phone2,
			v.address, v.district, v.province, v.zip_code, v.note,
			v.update_by, v.update_date, v.id_status
		FROM tb_vendor v
		LEFT JOIN tb_vendor_type t ON v.vendor_type_id = t.vendor_type_id
		LEFT JOIN tb_discount d ON v.discount_id = d.discount_id
		WHERE v.is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println("SQL Query Error:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch vendors", Data: nil})
	}
	defer rows.Close()

	var result []models.Vendor
	for rows.Next() {
		var v models.Vendor
		var upd sql.NullTime
		if err := rows.Scan(
			&v.VendorID, &v.VendorName, &v.VendorTypeID, &v.TypeName,
			&v.DiscountID, &v.DiscountName,
			&v.ContactName1, &v.ContactName2, &v.Email, &v.Phone1, &v.Phone2,
			&v.Address, &v.District, &v.Province, &v.ZipCode, &v.Note,
			&v.UpdateBy, &upd, &v.IDStatus); err != nil {
			log.Println("Rows Scan Error:", err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data", Data: nil})
		}
		if upd.Valid { t := upd.Time.Format("2006-01-02T15:04:05"); v.UpdateDate = &t }
		result = append(result, v)
	}
	if err := rows.Err(); err != nil { log.Println("Rows Iteration Error:", err); return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Database iteration error", Data: nil}) }
	return c.JSON(models.ApiResponse{Status: "success", Message: "Vendors retrieved successfully", Data: result})
}


// SelectPageVendor: ดึงข้อมูล Vendor แบบแบ่งหน้า (Pagination)
func SelectPageVendor(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }
	offset := (page - 1) * limit

	query := `
		SELECT v.vendor_id, v.vendor_name, v.vendor_type_id, t.type_name,
			v.discount_id, d.discount_name,
			v.contact_name1, v.contact_name2, v.email, v.phone1, v.phone2,
			v.address, v.district, v.province, v.zip_code, v.note,
			v.update_by, v.update_date, v.id_status
		FROM tb_vendor v
		LEFT JOIN tb_vendor_type t ON v.vendor_type_id = t.vendor_type_id
		LEFT JOIN tb_discount d ON v.discount_id = d.discount_id
		WHERE v.is_delete = 0
		ORDER BY v.update_date DESC
		OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY
	`

	rows, err := config.DB.Query(query, offset, limit) 
	if err != nil {
		log.Println("SQL Query Error:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch vendors", Data: nil})
	}
	defer rows.Close()

	var result []models.Vendor
	for rows.Next() {
		var v models.Vendor
		var upd sql.NullTime
		if err := rows.Scan(
			&v.VendorID, &v.VendorName, &v.VendorTypeID, &v.TypeName,
			&v.DiscountID, &v.DiscountName,
			&v.ContactName1, &v.ContactName2, &v.Email, &v.Phone1, &v.Phone2,
			&v.Address, &v.District, &v.Province, &v.ZipCode, &v.Note,
			&v.UpdateBy, &upd, &v.IDStatus); err != nil {
			log.Println("Rows Scan Error:", err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data", Data: nil})
		}
		if upd.Valid { t := upd.Time.Format("2006-01-02T15:04:05"); v.UpdateDate = &t }
		result = append(result, v)
	}
	
	if err := rows.Err(); err != nil { log.Println("Rows Iteration Error:", err); return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Database iteration error", Data: nil}) }

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_vendor WHERE is_delete = 0`).Scan(&total)
	if err != nil {
		log.Println("Count Query Error:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to count records", Data: nil})
	}
	
	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return c.JSON(models.ApiResponse{
		Status: "success",
		Data: fiber.Map{
			"page": page, "limit": limit, "total": total, "total_page": totalPage, "vendor": result,
		},
	})
}


// SelectVendorByID: ดึงข้อมูล Vendor ตาม vendor_id
func SelectVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")

	query := `
		SELECT v.vendor_id, v.vendor_name, v.vendor_type_id, t.type_name,
			v.discount_id, d.discount_name,
			v.contact_name1, v.contact_name2, v.email, v.phone1, v.phone2,
			v.address, v.district, v.province, v.zip_code, v.note,
			v.update_by, v.update_date, v.id_status
		FROM tb_vendor v
		LEFT JOIN tb_vendor_type t ON v.vendor_type_id = t.vendor_type_id
		LEFT JOIN tb_discount d ON v.discount_id = d.discount_id
		WHERE v.vendor_id = @p1 AND v.is_delete = 0
	`
	
	row := config.DB.QueryRow(query, id)

	var v models.Vendor
	var upd sql.NullTime
	if err := row.Scan(
		&v.VendorID, &v.VendorName, &v.VendorTypeID, &v.TypeName,
		&v.DiscountID, &v.DiscountName,
		&v.ContactName1, &v.ContactName2, &v.Email, &v.Phone1, &v.Phone2,
		&v.Address, &v.District, &v.Province, &v.ZipCode, &v.Note,
		&v.UpdateBy, &upd, &v.IDStatus); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "Vendor not found", Data: nil})
		}
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data", Data: nil})
	}
	if upd.Valid { t := upd.Time.Format("2006-01-02T15:04:05"); v.UpdateDate = &t }

	return c.JSON(models.ApiResponse{Status: "success", Message: "", Data: v})
}


// SelectVendorByName: ดึงข้อมูล Vendor ตามชื่อ (LIKE search)
func SelectVendorByName(c *fiber.Ctx) error {
	name := c.Params("name")

	query := `
		SELECT v.vendor_id, v.vendor_type_id, t.type_name, v.discount_id, v.vendor_name,
			v.contact_name1, v.contact_name2, v.email, v.phone1, v.phone2,
			v.address, v.district, v.province, v.zip_code, v.note,
			v.update_by, v.update_date, v.id_status
		FROM tb_vendor v
		LEFT JOIN tb_vendor_type t ON v.vendor_type_id = t.vendor_type_id
		WHERE v.vendor_name LIKE '%' + @p1 + '%' AND v.is_delete = 0
	`

	rows, err := config.DB.Query(query, name)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch vendors", Data: nil})
	}
	defer rows.Close()

	var result []models.Vendor
	for rows.Next() {
		var v models.Vendor
		var upd sql.NullTime
		if err := rows.Scan(
			&v.VendorID, &v.VendorTypeID, &v.TypeName, &v.DiscountID, &v.VendorName,
			&v.ContactName1, &v.ContactName2, &v.Email, &v.Phone1, &v.Phone2,
			&v.Address, &v.District, &v.Province, &v.ZipCode, &v.Note,
			&v.UpdateBy, &upd, &v.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data", Data: nil})
		}
		if upd.Valid { t := upd.Time.Format("2006-01-02T15:04:05"); v.UpdateDate = &t }
		result = append(result, v)
	}

	if len(result) == 0 {
		return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "No matching vendor found", Data: nil})
	}

	return c.JSON(models.ApiResponse{Status: "success", Message: "", Data: result})
}


// InsertVendor: เพิ่ม Vendor ใหม่ (ใช้ Named Parameter)
func InsertVendor(c *fiber.Ctx) error {
	var v models.Vendor
	if err := c.BodyParser(&v); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid request body", Data: nil})
	}

	query := `
		INSERT INTO tb_vendor (
			vendor_type_id, discount_id, vendor_name, contact_name1, contact_name2,
			email, phone1, phone2, address, district, province, zip_code, note, update_by, id_status
		)
		VALUES (
			@TypeID, @DiscountID, @Name, @Contact1, @Contact2,
			@Email, @Phone1, @Phone2, @Address, @District, @Province, @Zip, @Note, @UpdateBy, @Status
		)
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeID", v.VendorTypeID),
				sql.Named("DiscountID", v.DiscountID),
				sql.Named("Name", v.VendorName),
				sql.Named("Contact1", v.ContactName1),
				sql.Named("Contact2", v.ContactName2),
				sql.Named("Email", v.Email),
				sql.Named("Phone1", v.Phone1),
				sql.Named("Phone2", v.Phone2),
				sql.Named("Address", v.Address),
				sql.Named("District", v.District),
				sql.Named("Province", v.Province),
				sql.Named("Zip", v.ZipCode),
				sql.Named("Note", v.Note),
				sql.Named("UpdateBy", v.UpdateBy),
				sql.Named("Status", v.IDStatus),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to insert vendor", Data: nil})
	}

	return c.Status(201).JSON(models.ApiResponse{Status: "success", Message: "Vendor added successfully", Data: nil})
}


// UpdateVendorByID: อัปเดตข้อมูล Vendor ตาม vendor_id
func UpdateVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var v models.Vendor
	if err := c.BodyParser(&v); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid request body", Data: nil})
	}

	query := `
		UPDATE tb_vendor SET
			vendor_type_id = COALESCE(NULLIF(@TypeID, ''), vendor_type_id),
			discount_id = COALESCE(NULLIF(@DiscountID, ''), discount_id),
			vendor_name = COALESCE(NULLIF(@Name, ''), vendor_name),
			contact_name1 = @Contact1,
			contact_name2 = @Contact2,
			email = @Email,
			phone1 = @Phone1,
			phone2 = @Phone2,
			address = @Address,
			district = @District,
			province = @Province,
			zip_code = @Zip,
			note = @Note,
			update_by = @UpdateBy,
			id_status = @Status
		WHERE vendor_id = @ID AND is_delete = 0
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeID", v.VendorTypeID),
				sql.Named("DiscountID", v.DiscountID),
				sql.Named("Name", v.VendorName),
				sql.Named("Contact1", v.ContactName1),
				sql.Named("Contact2", v.ContactName2),
				sql.Named("Email", v.Email),
				sql.Named("Phone1", v.Phone1),
				sql.Named("Phone2", v.Phone2),
				sql.Named("Address", v.Address),
				sql.Named("District", v.District),
				sql.Named("Province", v.Province),
				sql.Named("Zip", v.ZipCode),
				sql.Named("Note", v.Note),
				sql.Named("UpdateBy", v.UpdateBy),
				sql.Named("Status", v.IDStatus),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to update vendor", Data: nil})
	}

	return c.JSON(models.ApiResponse{Status: "success", Message: "Vendor updated successfully", Data: nil})
}


// DeleteVendorByID: Soft Delete
func DeleteVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_vendor
		SET is_delete = 1,
			update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE vendor_id = @ID
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, sql.Named("ID", id))
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to soft delete vendor", Data: nil})
	}
	return c.JSON(models.ApiResponse{Status: "success", Message: "Vendor marked as deleted", Data: nil})
}


// RemoveVendorByID: Hard Delete
func RemoveVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_vendor WHERE vendor_id = @ID`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, sql.Named("ID", id))
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to hard delete vendor", Data: nil})
	}
	return c.JSON(models.ApiResponse{Status: "success", Message: "Vendor removed successfully", Data: nil})
}