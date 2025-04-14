package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// 1. Select All
func SelectAllCustomerTypes(c *fiber.Ctx) error {
	query := `
		SELECT customer_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_customer_type
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}
	defer rows.Close()

	var list []models.CustomerType
	for rows.Next() {
		var item models.CustomerType
		if err := rows.Scan(&item.CustomerTypeID, &item.TypeName, &item.Description, &item.UpdateBy, &item.UpdateDate, &item.IDStatus); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to read row"})
		}
		list = append(list, item)
	}
	return c.JSON(list)
}

// 2. Select Paging
func SelectPageCustomerTypes(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT customer_type_id, type_name, description, update_by, update_date
		FROM tb_customer_type
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}
	defer rows.Close()

	var list []models.CustomerType
	for rows.Next() {
		var item models.CustomerType
		if err := rows.Scan(&item.CustomerTypeID, &item.TypeName, &item.Description, &item.UpdateBy, &item.UpdateDate); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to read row"})
		}
		list = append(list, item)
	}
	return c.JSON(fiber.Map{"page": page, "limit": limit, "data": list})
}

// 3. Select By ID
func SelectCustomerTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT customer_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_customer_type
		WHERE customer_type_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))
	var item models.CustomerType
	if err := row.Scan(&item.CustomerTypeID, &item.TypeName, &item.Description, &item.UpdateBy, &item.UpdateDate, &item.IDStatus); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}
	return c.JSON(item)
}

// 4. Insert
func InsertCustomerType(c *fiber.Ctx) error {
	var item models.CustomerType
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}
	query := `
		INSERT INTO tb_customer_type (type_name, description, update_by)
		VALUES (@TypeName, @Description, @UpdateBy)
	`
	_, err := config.DB.Exec(query,
		sql.Named("TypeName", item.TypeName),
		sql.Named("Description", item.Description),
		sql.Named("UpdateBy", item.UpdateBy),
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to insert"})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Created"})
}

// 5. Update By ID
func UpdateCustomerTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.CustomerType
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}
	query := `
		UPDATE tb_customer_type
		SET type_name = @TypeName, description = @Description, update_by = @UpdateBy
		WHERE customer_type_id = @ID AND is_delete = 0
	`
	_, err := config.DB.Exec(query,
		sql.Named("TypeName", item.TypeName),
		sql.Named("Description", item.Description),
		sql.Named("UpdateBy", item.UpdateBy),
		sql.Named("ID", id),
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update"})
	}
	return c.JSON(fiber.Map{"message": "Updated"})
}

// 6. Delete (Soft)
func DeleteCustomerTypeByID(c *fiber.Ctx) error {
	customerTypeID := c.Params("id")

	query := `
		UPDATE tb_customer_type
		SET is_delete = 1,
		    id_status = 0,
		    update_date = GETDATE()
		WHERE customer_type_id = @CustomerTypeID
	`

	_, err := config.DB.Exec(query, sql.Named("CustomerTypeID", customerTypeID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to soft delete customer type"})
	}

	return c.JSON(fiber.Map{"message": "Soft deleted"})
}

// 7. Remove (Hard)
func RemoveCustomerTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_customer_type WHERE customer_type_id = @ID`
	_, err := config.DB.Exec(query, sql.Named("ID", id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to remove"})
	}
	return c.JSON(fiber.Map{"message": "Hard deleted"})
}
