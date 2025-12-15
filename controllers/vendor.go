package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SelectAllVendors(c *fiber.Ctx) error {
	query := `
		SELECT v.vendor_id, v.vendor_type_id, v.vendor_name, v.tax_id, v.branch_name,
		       v.contact_person, v.phone1, v.phone2, v.email, v.website,
		       v.address, v.sub_district, v.district, v.province, v.zip_code,
		       v.credit_term_day, v.currency, v.note,
		       v.update_by, v.update_date, v.is_active,
		       vt.type_name
		FROM tb_vendor v
		LEFT JOIN tb_vendor_type vt ON v.vendor_type_id = vt.vendor_type_id
		WHERE v.is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch vendors",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.Vendor
	for rows.Next() {
		var item models.Vendor
		var upd sql.NullTime
		if err := rows.Scan(
			&item.VendorID, &item.VendorTypeID, &item.VendorName, &item.TaxID, &item.BranchName,
			&item.ContactPerson, &item.Phone1, &item.Phone2, &item.Email, &item.Website,
			&item.Address, &item.SubDistrict, &item.District, &item.Province, &item.ZipCode,
			&item.CreditTermDay, &item.Currency, &item.Note,
			&item.UpdateBy, &upd, &item.IsActive, &item.VendorTypeName,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		if upd.Valid {
			t := upd.Time.Format("2006-01-02T15:04:05")
			item.UpdateDate = &t
		}
		list = append(list, item)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    list,
	})
}

func SelectPageVendors(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT v.vendor_id, v.vendor_type_id, v.vendor_name, v.tax_id, v.branch_name,
		       v.contact_person, v.phone1, v.phone2, v.email, v.website,
		       v.address, v.sub_district, v.district, v.province, v.zip_code,
		       v.credit_term_day, v.currency, v.note,
		       v.update_by, v.update_date, v.is_active,
		       vt.type_name
		FROM tb_vendor v
		LEFT JOIN tb_vendor_type vt ON v.vendor_type_id = vt.vendor_type_id
		WHERE v.is_delete = 0
		ORDER BY v.update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch vendors",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.Vendor
	for rows.Next() {
		var item models.Vendor
		var upd sql.NullTime
		if err := rows.Scan(
			&item.VendorID, &item.VendorTypeID, &item.VendorName, &item.TaxID, &item.BranchName,
			&item.ContactPerson, &item.Phone1, &item.Phone2, &item.Email, &item.Website,
			&item.Address, &item.SubDistrict, &item.District, &item.Province, &item.ZipCode,
			&item.CreditTermDay, &item.Currency, &item.Note,
			&item.UpdateBy, &upd, &item.IsActive, &item.VendorTypeName,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		if upd.Valid {
			t := upd.Time.Format("2006-01-02T15:04:05")
			item.UpdateDate = &t
		}
		list = append(list, item)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_vendor WHERE is_delete = 0`).Scan(&total)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to count records",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status: "success",
		Data: fiber.Map{
			"page":   page,
			"limit":  limit,
			"total":  total,
			"vendors": list,
		},
	})
}

func SelectVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT v.vendor_id, v.vendor_type_id, v.vendor_name, v.tax_id, v.branch_name,
		       v.contact_person, v.phone1, v.phone2, v.email, v.website,
		       v.address, v.sub_district, v.district, v.province, v.zip_code,
		       v.credit_term_day, v.currency, v.note,
		       v.update_by, v.update_date, v.is_active,
		       vt.type_name
		FROM tb_vendor v
		LEFT JOIN tb_vendor_type vt ON v.vendor_type_id = vt.vendor_type_id
		WHERE v.vendor_id = @ID AND v.is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var item models.Vendor
	var upd sql.NullTime
	if err := row.Scan(
		&item.VendorID, &item.VendorTypeID, &item.VendorName, &item.TaxID, &item.BranchName,
		&item.ContactPerson, &item.Phone1, &item.Phone2, &item.Email, &item.Website,
		&item.Address, &item.SubDistrict, &item.District, &item.Province, &item.ZipCode,
		&item.CreditTermDay, &item.Currency, &item.Note,
		&item.UpdateBy, &upd, &item.IsActive, &item.VendorTypeName,
	); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Vendor not found",
				Data:    nil,
			})
		}
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to read data",
			Data:    nil,
		})
	}
	if upd.Valid {
		t := upd.Time.Format("2006-01-02T15:04:05")
		item.UpdateDate = &t
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    item,
	})
}

func SelectVendorByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT v.vendor_id, v.vendor_type_id, v.vendor_name, v.tax_id, v.branch_name,
		       v.contact_person, v.phone1, v.phone2, v.email, v.website,
		       v.address, v.sub_district, v.district, v.province, v.zip_code,
		       v.credit_term_day, v.currency, v.note,
		       v.update_by, v.update_date, v.is_active,
		       vt.type_name
		FROM tb_vendor v
		LEFT JOIN tb_vendor_type vt ON v.vendor_type_id = vt.vendor_type_id
		WHERE v.vendor_name LIKE '%' + @Name + '%' AND v.is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch vendors",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.Vendor
	for rows.Next() {
		var item models.Vendor
		var upd sql.NullTime
		if err := rows.Scan(
			&item.VendorID, &item.VendorTypeID, &item.VendorName, &item.TaxID, &item.BranchName,
			&item.ContactPerson, &item.Phone1, &item.Phone2, &item.Email, &item.Website,
			&item.Address, &item.SubDistrict, &item.District, &item.Province, &item.ZipCode,
			&item.CreditTermDay, &item.Currency, &item.Note,
			&item.UpdateBy, &upd, &item.IsActive, &item.VendorTypeName,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		if upd.Valid {
			t := upd.Time.Format("2006-01-02T15:04:05")
			item.UpdateDate = &t
		}
		list = append(list, item)
	}

	if len(list) == 0 {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "No matching vendor found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    list,
	})
}

func InsertVendor(c *fiber.Ctx) error {
	var item models.Vendor
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	if item.VendorTypeID == "" {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "vendor_type_id is required"})
	}
	if strings.TrimSpace(item.VendorName) == "" {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "vendor_name is required"})
	}

	if item.UpdateBy == nil || *item.UpdateBy == "" {
		u := utils.ResolveUser(c)
		item.UpdateBy = &u
	}

	query := `
		INSERT INTO tb_vendor (
			vendor_type_id, vendor_name, tax_id, branch_name,
			contact_person, phone1, phone2, email, website,
			address, sub_district, district, province, zip_code,
			credit_term_day, currency, note, update_by
		)
		VALUES (
			@TypeID, @Name, @TaxID, @Branch,
			@Contact, @Phone1, @Phone2, @Email, @Website,
			@Address, @SubDistrict, @District, @Province, @ZipCode,
			COALESCE(@CreditTerm, 30), COALESCE(@Currency, 'THB'), @Note, @UpdateBy
		)
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeID", item.VendorTypeID),
				sql.Named("Name", item.VendorName),
				sql.Named("TaxID", item.TaxID),
				sql.Named("Branch", item.BranchName),
				sql.Named("Contact", item.ContactPerson),
				sql.Named("Phone1", item.Phone1),
				sql.Named("Phone2", item.Phone2),
				sql.Named("Email", item.Email),
				sql.Named("Website", item.Website),
				sql.Named("Address", item.Address),
				sql.Named("SubDistrict", item.SubDistrict),
				sql.Named("District", item.District),
				sql.Named("Province", item.Province),
				sql.Named("ZipCode", item.ZipCode),
				sql.Named("CreditTerm", item.CreditTermDay),
				sql.Named("Currency", item.Currency),
				sql.Named("Note", item.Note),
				sql.Named("UpdateBy", item.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to insert vendor",
			Data:    nil,
		})
	}

	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor added successfully",
		Data:    nil,
	})
}

func UpdateVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.Vendor
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	if item.UpdateBy == nil || *item.UpdateBy == "" {
		u := utils.ResolveUser(c)
		item.UpdateBy = &u
	}

	query := `
		UPDATE tb_vendor
		SET vendor_type_id = COALESCE(NULLIF(@TypeID, ''), vendor_type_id),
			vendor_name = COALESCE(NULLIF(@Name, ''), vendor_name),
			tax_id = COALESCE(@TaxID, tax_id),
			branch_name = COALESCE(@Branch, branch_name),
			contact_person = COALESCE(@Contact, contact_person),
			phone1 = COALESCE(@Phone1, phone1),
			phone2 = COALESCE(@Phone2, phone2),
			email = COALESCE(@Email, email),
			website = COALESCE(@Website, website),
			address = COALESCE(@Address, address),
			sub_district = COALESCE(@SubDistrict, sub_district),
			district = COALESCE(@District, district),
			province = COALESCE(@Province, province),
			zip_code = COALESCE(@ZipCode, zip_code),
			credit_term_day = COALESCE(@CreditTerm, credit_term_day),
			currency = COALESCE(@Currency, currency),
			note = COALESCE(@Note, note),
			update_by = @UpdateBy,
			is_active = COALESCE(@IsActive, is_active)
		WHERE vendor_id = @ID AND is_delete = 0
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(query,
				sql.Named("TypeID", item.VendorTypeID),
				sql.Named("Name", item.VendorName),
				sql.Named("TaxID", item.TaxID),
				sql.Named("Branch", item.BranchName),
				sql.Named("Contact", item.ContactPerson),
				sql.Named("Phone1", item.Phone1),
				sql.Named("Phone2", item.Phone2),
				sql.Named("Email", item.Email),
				sql.Named("Website", item.Website),
				sql.Named("Address", item.Address),
				sql.Named("SubDistrict", item.SubDistrict),
				sql.Named("District", item.District),
				sql.Named("Province", item.Province),
				sql.Named("ZipCode", item.ZipCode),
				sql.Named("CreditTerm", item.CreditTermDay),
				sql.Named("Currency", item.Currency),
				sql.Named("Note", item.Note),
				sql.Named("UpdateBy", item.UpdateBy),
				sql.Named("IsActive", item.IsActive),
				sql.Named("ID", id),
			)
			if err != nil {
				return err
			}
			rows, err := res.RowsAffected()
			if err != nil {
				return err
			}
			if rows == 0 {
				return sql.ErrNoRows
			}
			return nil
		},
	})
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Vendor not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update vendor",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor updated successfully",
		Data:    nil,
	})
}

func DeleteVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	username := utils.ResolveUser(c)

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(`
				UPDATE tb_vendor
				SET is_delete = 1,
					is_active = 0,
					update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME),
					update_by = @UpdateBy
				WHERE vendor_id = @ID AND is_delete = 0`,
				sql.Named("ID", id),
				sql.Named("UpdateBy", username),
			)
			if err != nil {
				return err
			}
			rows, err := res.RowsAffected()
			if err != nil {
				return err
			}
			if rows == 0 {
				return sql.ErrNoRows
			}
			return nil
		},
	})
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Vendor not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to delete vendor",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor deleted successfully",
		Data:    nil,
	})
}

func RemoveVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(`DELETE FROM tb_vendor WHERE vendor_id = @ID`, sql.Named("ID", id))
			if err != nil {
				return err
			}
			rows, err := res.RowsAffected()
			if err != nil {
				return err
			}
			if rows == 0 {
				return sql.ErrNoRows
			}
			return nil
		},
	})
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Vendor not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to remove vendor",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor removed successfully",
		Data:    nil,
	})
}