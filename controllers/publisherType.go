package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// Select All Publisher Types
func SelectAllPublisherTypes(c *fiber.Ctx) error {
	query := `
		SELECT publisher_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_publisher_type
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch publisher types"})
	}
	defer rows.Close()

	var publisherTypes []models.PublisherType

	for rows.Next() {
		var publisherType models.PublisherType
		if err := rows.Scan(
			&publisherType.PublisherTypeID, &publisherType.TypeName, &publisherType.Description,
			&publisherType.UpdateBy, &publisherType.UpdateDate, &publisherType.IDStatus,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
		}
		publisherTypes = append(publisherTypes, publisherType)
	}

	return c.JSON(publisherTypes)
}

// Select Publisher Types By Paging
func SelectPagePublisherTypes(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT publisher_type_id, type_name, description, update_by, update_date
		FROM tb_publisher_type
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch publisher types"})
	}
	defer rows.Close()

	var publisherTypes []models.PublisherType

	for rows.Next() {
		var publisherType models.PublisherType
		if err := rows.Scan(
			&publisherType.PublisherTypeID, &publisherType.TypeName, &publisherType.Description,
			&publisherType.UpdateBy, &publisherType.UpdateDate,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
		}
		publisherTypes = append(publisherTypes, publisherType)
	}

	return c.JSON(fiber.Map{
		"page":       page,
		"limit":      limit,
		"data":       publisherTypes,
		"totalItems": len(publisherTypes), // Optional: Adjust if total count is needed
	})
}

// Select Publisher Type By ID
func SelectPublisherTypeByID(c *fiber.Ctx) error {
	publisherTypeID := c.Params("id")
	query := `
		SELECT publisher_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_publisher_type
		WHERE publisher_type_id = @PublisherTypeID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("PublisherTypeID", publisherTypeID))

	var publisherType models.PublisherType
	if err := row.Scan(
		&publisherType.PublisherTypeID, &publisherType.TypeName, &publisherType.Description,
		&publisherType.UpdateBy, &publisherType.UpdateDate, &publisherType.IDStatus,
	); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Publisher type not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
	}

	return c.JSON(publisherType)
}

// Insert Publisher Type
func InsertPublisherType(c *fiber.Ctx) error {
	var publisherType models.PublisherType
	if err := c.BodyParser(&publisherType); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	query := `
		INSERT INTO tb_publisher_type (type_name, description, update_by)
		VALUES (@TypeName, @Description, @UpdateBy)
	`
	_, err := config.DB.Exec(query,
		sql.Named("TypeName", publisherType.TypeName),
		sql.Named("Description", publisherType.Description),
		sql.Named("UpdateBy", publisherType.UpdateBy),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert publisher type"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Publisher type added successfully"})
}

// Update Publisher Type By ID
func UpdatePublisherTypeByID(c *fiber.Ctx) error {
	publisherTypeID := c.Params("id")

	var updatedPublisherType models.PublisherType
	if err := c.BodyParser(&updatedPublisherType); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	query := `
		UPDATE tb_publisher_type
		SET type_name = @TypeName, description = @Description, update_by = @UpdateBy
		WHERE publisher_type_id = @PublisherTypeID AND is_delete = 0
	`
	_, err := config.DB.Exec(query,
		sql.Named("TypeName", updatedPublisherType.TypeName),
		sql.Named("Description", updatedPublisherType.Description),
		sql.Named("UpdateBy", updatedPublisherType.UpdateBy),
		sql.Named("PublisherTypeID", publisherTypeID),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update publisher type"})
	}

	return c.JSON(fiber.Map{"message": "Publisher type updated successfully"})
}

// Delete Publisher Type (Soft Delete)
func DeletePublisherTypeByID(c *fiber.Ctx) error {
	publisherTypeID := c.Params("id")
	query := `
		UPDATE tb_publisher_type
		SET is_delete = 1, update_date = GETDATE()
		WHERE publisher_type_id = @PublisherTypeID
	`
	_, err := config.DB.Exec(query, sql.Named("PublisherTypeID", publisherTypeID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete publisher type"})
	}

	return c.JSON(fiber.Map{"message": "Publisher type marked as deleted"})
}

// Remove Publisher Type (Hard Delete)
func RemovePublisherTypeByID(c *fiber.Ctx) error {
	publisherTypeID := c.Params("id")
	query := `
		DELETE FROM tb_publisher_type WHERE publisher_type_id = @PublisherTypeID
	`
	_, err := config.DB.Exec(query, sql.Named("PublisherTypeID", publisherTypeID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove publisher type"})
	}

	return c.JSON(fiber.Map{"message": "Publisher type removed successfully"})
}
