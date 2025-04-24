package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func SelectAllCustomers(c *fiber.Ctx) error {
	query := `
	SELECT customer_code, customer_name, biz_id, customer_type_id, contract_id, discount_id,
	       first_name, last_name, email, phone1, phone2, tax_id, address, district, province, zip_code,
	       note, reference1, reference2, update_by, update_date, id_status
	FROM tb_customer
	WHERE is_delete = 0
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(
			&customer.CustomerCode, &customer.CustomerName, &customer.BizID, &customer.CustomerTypeID,
			&customer.ContractID, &customer.DiscountID,
			&customer.FirstName, &customer.LastName, &customer.Email, &customer.Phone1, &customer.Phone2,
			&customer.TaxID, &customer.Address, &customer.District, &customer.Province, &customer.ZipCode,
			&customer.Note, &customer.Reference1, &customer.Reference2,
			&customer.UpdateBy, &customer.UpdateDate, &customer.IDStatus,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan data"})
		}
		customers = append(customers, customer)
	}

	return c.JSON(customers)
}

func SelectPageCustomers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	offset := (page - 1) * limit

	query := `
	SELECT customer_code, customer_name, biz_id, customer_type_id, contract_id, discount_id,
	       first_name, last_name, email, phone1, phone2, tax_id, address, district, province, zip_code,
	       note, reference1, reference2, update_by, update_date, id_status
	FROM tb_customer
	WHERE is_delete = 0
	ORDER BY update_date DESC
	OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query,
		sql.Named("Offset", offset),
		sql.Named("Limit", limit),
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(
			&customer.CustomerCode, &customer.CustomerName, &customer.BizID, &customer.CustomerTypeID,
			&customer.ContractID, &customer.DiscountID,
			&customer.FirstName, &customer.LastName, &customer.Email, &customer.Phone1, &customer.Phone2,
			&customer.TaxID, &customer.Address, &customer.District, &customer.Province, &customer.ZipCode,
			&customer.Note, &customer.Reference1, &customer.Reference2,
			&customer.UpdateBy, &customer.UpdateDate, &customer.IDStatus,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan data"})
		}
		customers = append(customers, customer)
	}

	return c.JSON(fiber.Map{
		"page":  page,
		"limit": limit,
		"data":  customers,
		"total": len(customers), // Optional, you can add COUNT() if needed
	})
}

func SelectCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
	SELECT customer_code, customer_name, biz_id, customer_type_id, contract_id, discount_id,
	       first_name, last_name, email, phone1, phone2, tax_id, address, district, province, zip_code,
	       note, reference1, reference2, update_by, update_date, id_status
	FROM tb_customer
	WHERE customer_code = @CustomerCode AND is_delete = 0
	`
	var customer models.Customer
	err := config.DB.QueryRow(query, sql.Named("CustomerCode", id)).Scan(
		&customer.CustomerCode, &customer.CustomerName, &customer.BizID, &customer.CustomerTypeID,
		&customer.ContractID, &customer.DiscountID,
		&customer.FirstName, &customer.LastName, &customer.Email, &customer.Phone1, &customer.Phone2,
		&customer.TaxID, &customer.Address, &customer.District, &customer.Province, &customer.ZipCode,
		&customer.Note, &customer.Reference1, &customer.Reference2,
		&customer.UpdateBy, &customer.UpdateDate, &customer.IDStatus,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Customer not found"})
	}

	return c.JSON(customer)
}

func InsertCustomer(c *fiber.Ctx) error {
	var customer models.Customer
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	query := `
	INSERT INTO tb_customer (
		customer_name, biz_id, customer_type_id, contract_id, discount_id,
		first_name, last_name, email, phone1, phone2, tax_id, address, district, province, zip_code,
		note, reference1, reference2, update_by
	) VALUES (
		@CustomerName, @BizID, @CustomerTypeID, @ContractID, @DiscountID,
		@FirstName, @LastName, @Email, @Phone1, @Phone2, @TaxID, @Address, @District, @Province, @ZipCode,
		@Note, @Reference1, @Reference2, @UpdateBy
	)
	`
	_, err := config.DB.Exec(query,
		sql.Named("CustomerName", customer.CustomerName),
		sql.Named("BizID", customer.BizID),
		sql.Named("CustomerTypeID", customer.CustomerTypeID),
		sql.Named("ContractID", customer.ContractID),
		sql.Named("DiscountID", customer.DiscountID),
		sql.Named("FirstName", customer.FirstName),
		sql.Named("LastName", customer.LastName),
		sql.Named("Email", customer.Email),
		sql.Named("Phone1", customer.Phone1),
		sql.Named("Phone2", customer.Phone2),
		sql.Named("TaxID", customer.TaxID),
		sql.Named("Address", customer.Address),
		sql.Named("District", customer.District),
		sql.Named("Province", customer.Province),
		sql.Named("ZipCode", customer.ZipCode),
		sql.Named("Note", customer.Note),
		sql.Named("Reference1", customer.Reference1),
		sql.Named("Reference2", customer.Reference2),
		sql.Named("UpdateBy", customer.UpdateBy),
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Insert failed"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Customer inserted"})
}

func UpdateCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var customer models.Customer
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	query := `
	UPDATE tb_customer SET
		customer_name = COALESCE(NULLIF(@CustomerName, ''), customer_name),
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
	WHERE customer_code = @CustomerCode AND is_delete = 0
	`

	_, err := config.DB.Exec(query,
		sql.Named("CustomerName", customer.CustomerName),
		sql.Named("BizID", customer.BizID),
		sql.Named("CustomerTypeID", customer.CustomerTypeID),
		sql.Named("ContractID", customer.ContractID),
		sql.Named("DiscountID", customer.DiscountID),
		sql.Named("FirstName", customer.FirstName),
		sql.Named("LastName", customer.LastName),
		sql.Named("Email", customer.Email),
		sql.Named("Phone1", customer.Phone1),
		sql.Named("Phone2", customer.Phone2),
		sql.Named("TaxID", customer.TaxID),
		sql.Named("Address", customer.Address),
		sql.Named("District", customer.District),
		sql.Named("Province", customer.Province),
		sql.Named("ZipCode", customer.ZipCode),
		sql.Named("Note", customer.Note),
		sql.Named("Reference1", customer.Reference1),
		sql.Named("Reference2", customer.Reference2),
		sql.Named("UpdateBy", customer.UpdateBy),
		sql.Named("CustomerCode", id),
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Update failed"})
	}

	return c.JSON(fiber.Map{"message": "Customer updated"})
}

func DeleteCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
	UPDATE tb_customer
	SET is_delete = 1, update_date = GETDATE()
	WHERE customer_code = @CustomerCode
	`
	_, err := config.DB.Exec(query, sql.Named("CustomerCode", id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Delete failed"})
	}

	return c.JSON(fiber.Map{"message": "Customer soft deleted"})
}

func RemoveCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
	DELETE FROM tb_customer
	WHERE customer_code = @CustomerCode
	`
	_, err := config.DB.Exec(query, sql.Named("CustomerCode", id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Remove failed"})
	}

	return c.JSON(fiber.Map{"message": "Customer permanently removed"})
}
