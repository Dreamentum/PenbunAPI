package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func executeDiscountTypeTransaction(db *sql.DB, queries []func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	for _, query := range queries {
		if err := query(tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func SelectAllDiscountType(c *fiber.Ctx) error {
	query := `
		SELECT discount_type_id, type_name, discount_unit_type, description, update_by, update_date, id_status
		FROM tb_discount_type
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch discount types"})
	}
	defer rows.Close()

	var discountTypes []models.DiscountType
	for rows.Next() {
		var dt models.DiscountType
		if err := rows.Scan(&dt.DiscountTypeID, &dt.TypeName, &dt.DiscountUnitType, &dt.Description, &dt.UpdateBy, &dt.UpdateDate, &dt.IDStatus); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
		}
		discountTypes = append(discountTypes, dt)
	}
	return c.JSON(discountTypes)
}

func SelectPageDiscountType(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT discount_type_id, type_name, discount_unit_type, description, update_by, update_date, id_status
		FROM tb_discount_type
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch discount types"})
	}
	defer rows.Close()

	var discountTypes []models.DiscountType
	for rows.Next() {
		var dt models.DiscountType
		if err := rows.Scan(&dt.DiscountTypeID, &dt.TypeName, &dt.DiscountUnitType, &dt.Description, &dt.UpdateBy, &dt.UpdateDate, &dt.IDStatus); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
		}
		discountTypes = append(discountTypes, dt)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_discount_type WHERE is_delete = 0`).Scan(&total)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count records"})
	}

	return c.JSON(fiber.Map{
		"page":          page,
		"limit":         limit,
		"total":         total,
		"discountTypes": discountTypes,
	})
}

func SelectDiscountTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT discount_type_id, type_name, discount_unit_type, description, update_by, update_date, id_status
		FROM tb_discount_type
		WHERE discount_type_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))
	var dt models.DiscountType
	if err := row.Scan(&dt.DiscountTypeID, &dt.TypeName, &dt.DiscountUnitType, &dt.Description, &dt.UpdateBy, &dt.UpdateDate, &dt.IDStatus); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Discount type not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
	}
	return c.JSON(dt)
}

func InsertDiscountType(c *fiber.Ctx) error {
	var dt models.DiscountType
	if err := c.BodyParser(&dt); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	query := `
		INSERT INTO tb_discount_type (type_name, discount_unit_type, description, update_by)
		VALUES (@TypeName, @DiscountUnitType, @Description, @UpdateBy)
	`
	err := executeDiscountTypeTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeName", dt.TypeName),
				sql.Named("DiscountUnitType", dt.DiscountUnitType),
				sql.Named("Description", dt.Description),
				sql.Named("UpdateBy", dt.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert discount type"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Discount type added successfully"})
}

func UpdateDiscountTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var dt models.DiscountType
	if err := c.BodyParser(&dt); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	query := `
		UPDATE tb_discount_type
		SET type_name = COALESCE(NULLIF(@TypeName, ''), type_name),
			discount_unit_type = COALESCE(NULLIF(@DiscountUnitType, ''), discount_unit_type),
			description = @Description,
			update_by = @UpdateBy
		WHERE discount_type_id = @ID AND is_delete = 0
	`
	err := executeDiscountTypeTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeName", dt.TypeName),
				sql.Named("DiscountUnitType", dt.DiscountUnitType),
				sql.Named("Description", dt.Description),
				sql.Named("UpdateBy", dt.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update discount type"})
	}
	return c.JSON(fiber.Map{"message": "Discount type updated successfully"})
}

func DeleteDiscountTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_discount_type
		SET is_delete = 1, update_date = GETDATE()
		WHERE discount_type_id = @ID
	`
	_, err := config.DB.Exec(query, sql.Named("ID", id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete discount type"})
	}
	return c.JSON(fiber.Map{"message": "Discount type marked as deleted"})
}

func RemoveDiscountTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_discount_type WHERE discount_type_id = @ID`
	_, err := config.DB.Exec(query, sql.Named("ID", id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove discount type"})
	}
	return c.JSON(fiber.Map{"message": "Discount type removed successfully"})
}
