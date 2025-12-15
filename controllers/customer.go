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

func SelectAllCustomers(c *fiber.Ctx) error {
	query := `
		SELECT c.customer_id, c.customer_type_id, c.customer_name, c.tax_id, c.branch_name,
		       c.contact_person, c.phone1, c.phone2, c.email, c.line_id,
		       c.address, c.sub_district, c.district, c.province, c.zip_code,
		       c.credit_limit, c.credit_term_day, c.note,
		       c.update_by, c.update_date, c.is_active,
		       ct.customer_type_name
		FROM tb_customer c
		LEFT JOIN tb_customer_type ct ON c.customer_type_id = ct.customer_type_id
		WHERE c.is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch customers",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.Customer
	for rows.Next() {
		var item models.Customer
		var upd sql.NullTime
		if err := rows.Scan(
			&item.CustomerID, &item.CustomerTypeID, &item.CustomerName, &item.TaxID, &item.BranchName,
			&item.ContactPerson, &item.Phone1, &item.Phone2, &item.Email, &item.LineID,
			&item.Address, &item.SubDistrict, &item.District, &item.Province, &item.ZipCode,
			&item.CreditLimit, &item.CreditTermDay, &item.Note,
			&item.UpdateBy, &upd, &item.IsActive, &item.CustomerTypeName,
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

func SelectPageCustomers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT c.customer_id, c.customer_type_id, c.customer_name, c.tax_id, c.branch_name,
		       c.contact_person, c.phone1, c.phone2, c.email, c.line_id,
		       c.address, c.sub_district, c.district, c.province, c.zip_code,
		       c.credit_limit, c.credit_term_day, c.note,
		       c.update_by, c.update_date, c.is_active,
		       ct.customer_type_name
		FROM tb_customer c
		LEFT JOIN tb_customer_type ct ON c.customer_type_id = ct.customer_type_id
		WHERE c.is_delete = 0
		ORDER BY c.update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch customers",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.Customer
	for rows.Next() {
		var item models.Customer
		var upd sql.NullTime
		if err := rows.Scan(
			&item.CustomerID, &item.CustomerTypeID, &item.CustomerName, &item.TaxID, &item.BranchName,
			&item.ContactPerson, &item.Phone1, &item.Phone2, &item.Email, &item.LineID,
			&item.Address, &item.SubDistrict, &item.District, &item.Province, &item.ZipCode,
			&item.CreditLimit, &item.CreditTermDay, &item.Note,
			&item.UpdateBy, &upd, &item.IsActive, &item.CustomerTypeName,
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
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_customer WHERE is_delete = 0`).Scan(&total)
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
			"page":     page,
			"limit":    limit,
			"total":    total,
			"customer": list,
		},
	})
}

func SelectCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT c.customer_id, c.customer_type_id, c.customer_name, c.tax_id, c.branch_name,
		       c.contact_person, c.phone1, c.phone2, c.email, c.line_id,
		       c.address, c.sub_district, c.district, c.province, c.zip_code,
		       c.credit_limit, c.credit_term_day, c.note,
		       c.update_by, c.update_date, c.is_active,
		       ct.customer_type_name
		FROM tb_customer c
		LEFT JOIN tb_customer_type ct ON c.customer_type_id = ct.customer_type_id
		WHERE c.customer_id = @ID AND c.is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var item models.Customer
	var upd sql.NullTime
	if err := row.Scan(
		&item.CustomerID, &item.CustomerTypeID, &item.CustomerName, &item.TaxID, &item.BranchName,
		&item.ContactPerson, &item.Phone1, &item.Phone2, &item.Email, &item.LineID,
		&item.Address, &item.SubDistrict, &item.District, &item.Province, &item.ZipCode,
		&item.CreditLimit, &item.CreditTermDay, &item.Note,
		&item.UpdateBy, &upd, &item.IsActive, &item.CustomerTypeName,
	); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Customer not found",
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

func SelectCustomerByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT c.customer_id, c.customer_type_id, c.customer_name, c.tax_id, c.branch_name,
		       c.contact_person, c.phone1, c.phone2, c.email, c.line_id,
		       c.address, c.sub_district, c.district, c.province, c.zip_code,
		       c.credit_limit, c.credit_term_day, c.note,
		       c.update_by, c.update_date, c.is_active,
		       ct.customer_type_name
		FROM tb_customer c
		LEFT JOIN tb_customer_type ct ON c.customer_type_id = ct.customer_type_id
		WHERE c.customer_name LIKE '%' + @Name + '%' AND c.is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch customers",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.Customer
	for rows.Next() {
		var item models.Customer
		var upd sql.NullTime
		if err := rows.Scan(
			&item.CustomerID, &item.CustomerTypeID, &item.CustomerName, &item.TaxID, &item.BranchName,
			&item.ContactPerson, &item.Phone1, &item.Phone2, &item.Email, &item.LineID,
			&item.Address, &item.SubDistrict, &item.District, &item.Province, &item.ZipCode,
			&item.CreditLimit, &item.CreditTermDay, &item.Note,
			&item.UpdateBy, &upd, &item.IsActive, &item.CustomerTypeName,
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
			Message: "No matching customer found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    list,
	})
}

func InsertCustomer(c *fiber.Ctx) error {
	var item models.Customer
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	if item.CustomerTypeID == "" {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "customer_type_id is required"})
	}
	if strings.TrimSpace(item.CustomerName) == "" {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "customer_name is required"})
	}

	if item.UpdateBy == nil || *item.UpdateBy == "" {
		u := utils.ResolveUser(c)
		item.UpdateBy = &u
	}

	query := `
		INSERT INTO tb_customer (
			customer_type_id, customer_name, tax_id, branch_name,
			contact_person, phone1, phone2, email, line_id,
			address, sub_district, district, province, zip_code,
			credit_limit, credit_term_day, note, update_by
		)
		VALUES (
			@TypeID, @Name, @TaxID, @Branch,
			@Contact, @Phone1, @Phone2, @Email, @LineID,
			@Address, @SubDistrict, @District, @Province, @ZipCode,
			COALESCE(@CreditLimit, 0), COALESCE(@CreditTerm, 0), @Note, @UpdateBy
		)
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeID", item.CustomerTypeID),
				sql.Named("Name", item.CustomerName),
				sql.Named("TaxID", item.TaxID),
				sql.Named("Branch", item.BranchName),
				sql.Named("Contact", item.ContactPerson),
				sql.Named("Phone1", item.Phone1),
				sql.Named("Phone2", item.Phone2),
				sql.Named("Email", item.Email),
				sql.Named("LineID", item.LineID),
				sql.Named("Address", item.Address),
				sql.Named("SubDistrict", item.SubDistrict),
				sql.Named("District", item.District),
				sql.Named("Province", item.Province),
				sql.Named("ZipCode", item.ZipCode),
				sql.Named("CreditLimit", item.CreditLimit),
				sql.Named("CreditTerm", item.CreditTermDay),
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
			Message: "Failed to insert customer",
			Data:    nil,
		})
	}

	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Customer added successfully",
		Data:    nil,
	})
}

func UpdateCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.Customer
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
		UPDATE tb_customer
		SET customer_type_id = COALESCE(NULLIF(@TypeID, ''), customer_type_id),
			customer_name = COALESCE(NULLIF(@Name, ''), customer_name),
			tax_id = COALESCE(@TaxID, tax_id),
			branch_name = COALESCE(@Branch, branch_name),
			contact_person = COALESCE(@Contact, contact_person),
			phone1 = COALESCE(@Phone1, phone1),
			phone2 = COALESCE(@Phone2, phone2),
			email = COALESCE(@Email, email),
			line_id = COALESCE(@LineID, line_id),
			address = COALESCE(@Address, address),
			sub_district = COALESCE(@SubDistrict, sub_district),
			district = COALESCE(@District, district),
			province = COALESCE(@Province, province),
			zip_code = COALESCE(@ZipCode, zip_code),
			credit_limit = COALESCE(@CreditLimit, credit_limit),
			credit_term_day = COALESCE(@CreditTerm, credit_term_day),
			note = COALESCE(@Note, note),
			update_by = @UpdateBy,
			is_active = COALESCE(@IsActive, is_active)
		WHERE customer_id = @ID AND is_delete = 0
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(query,
				sql.Named("TypeID", item.CustomerTypeID),
				sql.Named("Name", item.CustomerName),
				sql.Named("TaxID", item.TaxID),
				sql.Named("Branch", item.BranchName),
				sql.Named("Contact", item.ContactPerson),
				sql.Named("Phone1", item.Phone1),
				sql.Named("Phone2", item.Phone2),
				sql.Named("Email", item.Email),
				sql.Named("LineID", item.LineID),
				sql.Named("Address", item.Address),
				sql.Named("SubDistrict", item.SubDistrict),
				sql.Named("District", item.District),
				sql.Named("Province", item.Province),
				sql.Named("ZipCode", item.ZipCode),
				sql.Named("CreditLimit", item.CreditLimit),
				sql.Named("CreditTerm", item.CreditTermDay),
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
			Message: "Customer not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update customer",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Customer updated successfully",
		Data:    nil,
	})
}

func DeleteCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	username := utils.ResolveUser(c)

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(`
				UPDATE tb_customer
				SET is_delete = 1,
					is_active = 0,
					update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME),
					update_by = @UpdateBy
				WHERE customer_id = @ID AND is_delete = 0`,
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
			Message: "Customer not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to delete customer",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Customer deleted successfully",
		Data:    nil,
	})
}

func RemoveCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(`DELETE FROM tb_customer WHERE customer_id = @ID`, sql.Named("ID", id))
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
			Message: "Customer not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to remove customer",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Customer removed successfully",
		Data:    nil,
	})
}
