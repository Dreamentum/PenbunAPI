package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllCustomers(c *fiber.Ctx) error {
	query := `
		SELECT customer_code, customer_name, biz_id, customer_type_id, contract_id, discount_id,
		       first_name, last_name, email, phone1, phone2, tax_id, address, district, province, zip_code,
		       note, reference1, reference2, update_by, update_date, id_status, is_delete
		FROM tb_customer
		WHERE is_delete = 0
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

	var customers []models.Customer
	for rows.Next() {
		var cust models.Customer
		if err := rows.Scan(
			&cust.CustomerCode, &cust.CustomerName, &cust.BizID, &cust.CustomerTypeID, &cust.ContractID, &cust.DiscountID,
			&cust.FirstName, &cust.LastName, &cust.Email, &cust.Phone1, &cust.Phone2, &cust.TaxID, &cust.Address,
			&cust.District, &cust.Province, &cust.ZipCode, &cust.Note, &cust.Reference1, &cust.Reference2,
			&cust.UpdateBy, &cust.UpdateDate, &cust.IDStatus, &cust.IsDelete,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read customer data",
				Data:    nil,
			})
		}
		customers = append(customers, cust)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    customers,
	})
}

func SelectPageCustomers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT customer_code, customer_name, biz_id, customer_type_id, contract_id, discount_id,
		       first_name, last_name, email, phone1, phone2, tax_id, address, district, province, zip_code,
		       note, reference1, reference2, update_by, update_date, id_status, is_delete
		FROM tb_customer
		WHERE is_delete = 0
		ORDER BY update_date DESC
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

	var customers []models.Customer
	for rows.Next() {
		var cust models.Customer
		if err := rows.Scan(
			&cust.CustomerCode, &cust.CustomerName, &cust.BizID, &cust.CustomerTypeID, &cust.ContractID, &cust.DiscountID,
			&cust.FirstName, &cust.LastName, &cust.Email, &cust.Phone1, &cust.Phone2, &cust.TaxID, &cust.Address,
			&cust.District, &cust.Province, &cust.ZipCode, &cust.Note, &cust.Reference1, &cust.Reference2,
			&cust.UpdateBy, &cust.UpdateDate, &cust.IDStatus, &cust.IsDelete,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read customer data",
				Data:    nil,
			})
		}
		customers = append(customers, cust)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_customer WHERE is_delete = 0`).Scan(&total)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to count customers",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status: "success",
		Data: fiber.Map{
			"page":     page,
			"limit":    limit,
			"total":    total,
			"customer": customers,
		},
	})
}

func SelectCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT customer_code, customer_name, biz_id, customer_type_id, contract_id, discount_id,
		       first_name, last_name, email, phone1, phone2, tax_id, address, district, province, zip_code,
		       note, reference1, reference2, update_by, update_date, id_status, is_delete
		FROM tb_customer
		WHERE customer_code = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))
	var cust models.Customer
	if err := row.Scan(
		&cust.CustomerCode, &cust.CustomerName, &cust.BizID, &cust.CustomerTypeID, &cust.ContractID, &cust.DiscountID,
		&cust.FirstName, &cust.LastName, &cust.Email, &cust.Phone1, &cust.Phone2, &cust.TaxID, &cust.Address,
		&cust.District, &cust.Province, &cust.ZipCode, &cust.Note, &cust.Reference1, &cust.Reference2,
		&cust.UpdateBy, &cust.UpdateDate, &cust.IDStatus, &cust.IsDelete,
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
			Message: "Failed to read customer",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    cust,
	})
}

func SelectCustomerByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT customer_code, customer_name, biz_id, customer_type_id, contract_id, discount_id,
		       first_name, last_name, email, phone1, phone2, tax_id, address, district, province, zip_code,
		       note, reference1, reference2, update_by, update_date, id_status, is_delete
		FROM tb_customer
		WHERE customer_name LIKE '%' + @Name + '%' AND is_delete = 0
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

	var results []models.Customer
	for rows.Next() {
		var cust models.Customer
		if err := rows.Scan(
			&cust.CustomerCode, &cust.CustomerName, &cust.BizID, &cust.CustomerTypeID, &cust.ContractID, &cust.DiscountID,
			&cust.FirstName, &cust.LastName, &cust.Email, &cust.Phone1, &cust.Phone2, &cust.TaxID, &cust.Address,
			&cust.District, &cust.Province, &cust.ZipCode, &cust.Note, &cust.Reference1, &cust.Reference2,
			&cust.UpdateBy, &cust.UpdateDate, &cust.IDStatus, &cust.IsDelete,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read customer",
				Data:    nil,
			})
		}
		results = append(results, cust)
	}

	if len(results) == 0 {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "No matching customer found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    results,
	})
}

func InsertCustomer(c *fiber.Ctx) error {
	var cust models.Customer
	if err := c.BodyParser(&cust); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request",
			Data:    nil,
		})
	}

	query := `
		INSERT INTO tb_customer (
			customer_name, biz_id, customer_type_id, contract_id, discount_id,
			first_name, last_name, email, phone1, phone2, tax_id,
			address, district, province, zip_code, note,
			reference1, reference2, update_by
		)
		VALUES (
			@CustomerName, @BizID, @CustomerTypeID, @ContractID, @DiscountID,
			@FirstName, @LastName, @Email, @Phone1, @Phone2, @TaxID,
			@Address, @District, @Province, @ZipCode, @Note,
			@Reference1, @Reference2, @UpdateBy
		)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("CustomerName", cust.CustomerName),
				sql.Named("BizID", cust.BizID),
				sql.Named("CustomerTypeID", cust.CustomerTypeID),
				sql.Named("ContractID", cust.ContractID),
				sql.Named("DiscountID", cust.DiscountID),
				sql.Named("FirstName", cust.FirstName),
				sql.Named("LastName", cust.LastName),
				sql.Named("Email", cust.Email),
				sql.Named("Phone1", cust.Phone1),
				sql.Named("Phone2", cust.Phone2),
				sql.Named("TaxID", cust.TaxID),
				sql.Named("Address", cust.Address),
				sql.Named("District", cust.District),
				sql.Named("Province", cust.Province),
				sql.Named("ZipCode", cust.ZipCode),
				sql.Named("Note", cust.Note),
				sql.Named("Reference1", cust.Reference1),
				sql.Named("Reference2", cust.Reference2),
				sql.Named("UpdateBy", cust.UpdateBy),
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
	var cust models.Customer
	if err := c.BodyParser(&cust); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		UPDATE tb_customer
		SET customer_name = @CustomerName,
		    biz_id = @BizID,
		    customer_type_id = @CustomerTypeID,
		    contract_id = @ContractID,
		    discount_id = @DiscountID,
		    first_name = @FirstName,
		    last_name = @LastName,
		    email = @Email,
		    phone1 = @Phone1,
		    phone2 = @Phone2,
		    tax_id = @TaxID,
		    address = @Address,
		    district = @District,
		    province = @Province,
		    zip_code = @ZipCode,
		    note = @Note,
		    reference1 = @Reference1,
		    reference2 = @Reference2,
		    update_by = @UpdateBy
		WHERE customer_code = @ID AND is_delete = 0
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("CustomerName", cust.CustomerName),
				sql.Named("BizID", cust.BizID),
				sql.Named("CustomerTypeID", cust.CustomerTypeID),
				sql.Named("ContractID", cust.ContractID),
				sql.Named("DiscountID", cust.DiscountID),
				sql.Named("FirstName", cust.FirstName),
				sql.Named("LastName", cust.LastName),
				sql.Named("Email", cust.Email),
				sql.Named("Phone1", cust.Phone1),
				sql.Named("Phone2", cust.Phone2),
				sql.Named("TaxID", cust.TaxID),
				sql.Named("Address", cust.Address),
				sql.Named("District", cust.District),
				sql.Named("Province", cust.Province),
				sql.Named("ZipCode", cust.ZipCode),
				sql.Named("Note", cust.Note),
				sql.Named("Reference1", cust.Reference1),
				sql.Named("Reference2", cust.Reference2),
				sql.Named("UpdateBy", cust.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
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
	query := `
		UPDATE tb_customer
		SET is_delete = 1,
		    id_status = 0,
		    update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE customer_code = @ID
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
			Status:  "error",
			Message: "Failed to soft delete customer",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Customer marked as deleted",
		Data:    nil,
	})
}

func RemoveCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_customer WHERE customer_code = @ID`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, sql.Named("ID", id))
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to hard delete customer",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Customer removed successfully",
		Data:    nil,
	})
}
