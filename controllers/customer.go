package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// Select All Customers
func SelectAllCustomers(c *fiber.Ctx) error {
	query := `
		SELECT customer_code, customer_name, biz_id, customer_type_id, contract_id, discount_id,
		       first_name, last_name, email, phone1, phone2, tax_id, address, district, province,
		       zip_code, note, reference1, reference2, update_by, update_date, id_status
		FROM tb_customer
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch customers"})
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		if err := rows.Scan(
			&customer.CustomerCode, &customer.CustomerName, &customer.BizID, &customer.CustomerTypeID, &customer.ContractID, &customer.DiscountID,
			&customer.FirstName, &customer.LastName, &customer.Email, &customer.Phone1, &customer.Phone2, &customer.TaxID, &customer.Address,
			&customer.District, &customer.Province, &customer.ZipCode, &customer.Note, &customer.Reference1, &customer.Reference2,
			&customer.UpdateBy, &customer.UpdateDate, &customer.IDStatus,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
		}
		customers = append(customers, customer)
	}

	return c.JSON(customers)
}

// Select Customers By Paging
func SelectPageCustomers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT customer_code, customer_name, biz_id, customer_type_id, contract_id, discount_id,
		       first_name, last_name, email, phone1, phone2, tax_id, address, district, province,
		       zip_code, note, reference1, reference2, update_by, update_date, id_status
		FROM tb_customer
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch customers"})
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		if err := rows.Scan(
			&customer.CustomerCode, &customer.CustomerName, &customer.BizID, &customer.CustomerTypeID, &customer.ContractID, &customer.DiscountID,
			&customer.FirstName, &customer.LastName, &customer.Email, &customer.Phone1, &customer.Phone2, &customer.TaxID, &customer.Address,
			&customer.District, &customer.Province, &customer.ZipCode, &customer.Note, &customer.Reference1, &customer.Reference2,
			&customer.UpdateBy, &customer.UpdateDate, &customer.IDStatus,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
		}
		customers = append(customers, customer)
	}

	return c.JSON(fiber.Map{
		"page":  page,
		"limit": limit,
		"data":  customers,
	})
}

// Select Customer By ID
func SelectCustomerByID(c *fiber.Ctx) error {
	customerCode := c.Params("id")

	query := `
		SELECT customer_code, customer_name, biz_id, customer_type_id, contract_id, discount_id,
		       first_name, last_name, email, phone1, phone2, tax_id, address, district, province,
		       zip_code, note, reference1, reference2, update_by, update_date, id_status
		FROM tb_customer
		WHERE customer_code = @CustomerCode AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("CustomerCode", customerCode))

	var customer models.Customer
	if err := row.Scan(
		&customer.CustomerCode, &customer.CustomerName, &customer.BizID, &customer.CustomerTypeID, &customer.ContractID, &customer.DiscountID,
		&customer.FirstName, &customer.LastName, &customer.Email, &customer.Phone1, &customer.Phone2, &customer.TaxID, &customer.Address,
		&customer.District, &customer.Province, &customer.ZipCode, &customer.Note, &customer.Reference1, &customer.Reference2,
		&customer.UpdateBy, &customer.UpdateDate, &customer.IDStatus,
	); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Customer not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
	}

	return c.JSON(customer)
}

// Insert Customer
func InsertCustomer(c *fiber.Ctx) error {
	var customer models.Customer
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	query := `
		INSERT INTO tb_customer (customer_name, biz_id, customer_type_id, contract_id, discount_id, first_name, last_name, email, phone1, phone2, tax_id, address, district, province, zip_code, note, reference1, reference2, update_by)
		VALUES (@CustomerName, @BizID, @CustomerTypeID, @ContractID, @DiscountID, @FirstName, @LastName, @Email, @Phone1, @Phone2, @TaxID, @Address, @District, @Province, @ZipCode, @Note, @Reference1, @Reference2, @UpdateBy)
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert customer"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Customer added successfully"})
}

// Update Customer By ID
func UpdateCustomerByID(c *fiber.Ctx) error {
	customerCode := c.Params("id")

	var customer models.Customer
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	query := `
		UPDATE tb_customer
		SET customer_name = @CustomerName, biz_id = @BizID, customer_type_id = @CustomerTypeID,
		    contract_id = @ContractID, discount_id = @DiscountID, first_name = @FirstName, last_name = @LastName,
		    email = @Email, phone1 = @Phone1, phone2 = @Phone2, tax_id = @TaxID, address = @Address,
		    district = @District, province = @Province, zip_code = @ZipCode, note = @Note,
		    reference1 = @Reference1, reference2 = @Reference2, update_by = @UpdateBy
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
		sql.Named("CustomerCode", customerCode),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update customer"})
	}

	return c.JSON(fiber.Map{"message": "Customer updated successfully"})
}

// Delete Customer By ID (Soft Delete)
func DeleteCustomerByID(c *fiber.Ctx) error {
	sCustomerID := c.Params("id") // customer_code ที่ส่งมา

	// ตรวจสอบว่ามีข้อมูลลูกค้านี้หรือไม่
	queryCheck := `
		SELECT COUNT(*) 
		FROM tb_customer 
		WHERE customer_code = @CustomerID AND is_delete = 0
	`
	var iCount int
	err := config.DB.QueryRow(queryCheck, sql.Named("CustomerID", sCustomerID)).Scan(&iCount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check customer",
		})
	}

	if iCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Customer not found",
		})
	}

	// ทำการ Soft Delete (update is_delete = 1)
	queryDelete := `
		UPDATE tb_customer
		SET is_delete = 1, update_date = GETDATE()
		WHERE customer_code = @CustomerID
	`

	_, err = config.DB.Exec(queryDelete, sql.Named("CustomerID", sCustomerID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete customer",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Customer marked as deleted successfully",
		"customer_id": sCustomerID,
	})
}

// Remove Customer (Hard Delete)
func RemoveCustomerByID(c *fiber.Ctx) error {
	customerCode := c.Params("id")
	query := `DELETE FROM tb_customer WHERE customer_code = @CustomerCode`
	_, err := config.DB.Exec(query, sql.Named("CustomerCode", customerCode))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove customer"})
	}

	return c.JSON(fiber.Map{"message": "Customer removed successfully"})
}
