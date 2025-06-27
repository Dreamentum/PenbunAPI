package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllVendor(c *fiber.Ctx) error {
	query := `
		SELECT vendor_code, vendor_type_id, discount_id, vendor_name,
			contact_name1, contact_name2, email, phone1, phone2,
			address, district, province, zip_code, note,
			update_by, update_date, id_status
		FROM tb_vendor
		WHERE is_delete = 0
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to fetch vendors", Data: nil,
		})
	}
	defer rows.Close()

	var result []models.Vendor
	for rows.Next() {
		var v models.Vendor
		if err := rows.Scan(
			&v.VendorCode, &v.VendorTypeID, &v.DiscountID, &v.VendorName,
			&v.ContactName1, &v.ContactName2, &v.Email, &v.Phone1, &v.Phone2,
			&v.Address, &v.District, &v.Province, &v.ZipCode, &v.Note,
			&v.UpdateBy, &v.UpdateDate, &v.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status: "error", Message: "Failed to read data", Data: nil,
			})
		}
		result = append(result, v)
	}

	return c.JSON(models.ApiResponse{
		Status: "success", Message: "", Data: result,
	})
}

func SelectPageVendor(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT vendor_code, vendor_type_id, discount_id, vendor_name,
			contact_name1, contact_name2, email, phone1, phone2,
			address, district, province, zip_code, note,
			update_by, update_date, id_status
		FROM tb_vendor
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to fetch vendors", Data: nil,
		})
	}
	defer rows.Close()

	var result []models.Vendor
	for rows.Next() {
		var v models.Vendor
		if err := rows.Scan(
			&v.VendorCode, &v.VendorTypeID, &v.DiscountID, &v.VendorName,
			&v.ContactName1, &v.ContactName2, &v.Email, &v.Phone1, &v.Phone2,
			&v.Address, &v.District, &v.Province, &v.ZipCode, &v.Note,
			&v.UpdateBy, &v.UpdateDate, &v.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status: "error", Message: "Failed to read data", Data: nil,
			})
		}
		result = append(result, v)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_vendor WHERE is_delete = 0`).Scan(&total)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to count records", Data: nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status: "success",
		Data: fiber.Map{
			"page":   page,
			"limit":  limit,
			"total":  total,
			"vendor": result,
		},
	})
}

func SelectVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT vendor_code, vendor_type_id, discount_id, vendor_name,
			contact_name1, contact_name2, email, phone1, phone2,
			address, district, province, zip_code, note,
			update_by, update_date, id_status
		FROM tb_vendor
		WHERE vendor_code = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var v models.Vendor
	if err := row.Scan(
		&v.VendorCode, &v.VendorTypeID, &v.DiscountID, &v.VendorName,
		&v.ContactName1, &v.ContactName2, &v.Email, &v.Phone1, &v.Phone2,
		&v.Address, &v.District, &v.Province, &v.ZipCode, &v.Note,
		&v.UpdateBy, &v.UpdateDate, &v.IDStatus); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status: "error", Message: "Vendor not found", Data: nil,
			})
		}
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to read data", Data: nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status: "success", Message: "", Data: v,
	})
}

func SelectVendorByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT vendor_code, vendor_type_id, discount_id, vendor_name,
			contact_name1, contact_name2, email, phone1, phone2,
			address, district, province, zip_code, note,
			update_by, update_date, id_status
		FROM tb_vendor
		WHERE vendor_name LIKE '%' + @Name + '%' AND is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to fetch vendors", Data: nil,
		})
	}
	defer rows.Close()

	var result []models.Vendor
	for rows.Next() {
		var v models.Vendor
		if err := rows.Scan(
			&v.VendorCode, &v.VendorTypeID, &v.DiscountID, &v.VendorName,
			&v.ContactName1, &v.ContactName2, &v.Email, &v.Phone1, &v.Phone2,
			&v.Address, &v.District, &v.Province, &v.ZipCode, &v.Note,
			&v.UpdateBy, &v.UpdateDate, &v.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status: "error", Message: "Failed to read data", Data: nil,
			})
		}
		result = append(result, v)
	}

	if len(result) == 0 {
		return c.Status(404).JSON(models.ApiResponse{
			Status: "error", Message: "No matching vendor found", Data: nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status: "success", Message: "", Data: result,
	})
}

func InsertVendor(c *fiber.Ctx) error {
	var v models.Vendor
	if err := c.BodyParser(&v); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status: "error", Message: "Invalid request body", Data: nil,
		})
	}

	query := `
		INSERT INTO tb_vendor (
			vendor_type_id, discount_id, vendor_name, contact_name1, contact_name2,
			email, phone1, phone2, address, district, province, zip_code, note, update_by
		)
		VALUES (
			@TypeID, @DiscountID, @Name, @Contact1, @Contact2,
			@Email, @Phone1, @Phone2, @Address, @District, @Province, @Zip, @Note, @UpdateBy
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
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to insert vendor", Data: nil,
		})
	}

	return c.Status(201).JSON(models.ApiResponse{
		Status: "success", Message: "Vendor added successfully", Data: nil,
	})
}

func UpdateVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var v models.Vendor
	if err := c.BodyParser(&v); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status: "error", Message: "Invalid request body", Data: nil,
		})
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
			update_by = @UpdateBy
		WHERE vendor_code = @ID AND is_delete = 0
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
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to update vendor", Data: nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status: "success", Message: "Vendor updated successfully", Data: nil,
	})
}

func DeleteVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_vendor
		SET is_delete = 1,
			update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE vendor_code = @ID
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, sql.Named("ID", id))
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to soft delete vendor", Data: nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status: "success", Message: "Vendor marked as deleted", Data: nil,
	})
}

func RemoveVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_vendor WHERE vendor_code = @ID`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, sql.Named("ID", id))
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to hard delete vendor", Data: nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status: "success", Message: "Vendor removed successfully", Data: nil,
	})
}
