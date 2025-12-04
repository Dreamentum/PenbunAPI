package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllDiscountType(c *fiber.Ctx) error {
	query := `
		SELECT discount_type_id, type_name, discount_unit_type, description,
		       update_by, update_date, id_status
		FROM tb_discount_type
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch discount types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.DiscountType
	for rows.Next() {
		var dt models.DiscountType
		if err := rows.Scan(
			&dt.DiscountTypeID, &dt.TypeName, &dt.DiscountUnitType,
			&dt.Description, &dt.UpdateBy, &dt.UpdateDate, &dt.IDStatus,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		list = append(list, dt)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    list,
	})
}

func SelectPageDiscountType(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT discount_type_id, type_name, discount_unit_type, description,
		       update_by, update_date, id_status
		FROM tb_discount_type
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch discount types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.DiscountType
	for rows.Next() {
		var dt models.DiscountType
		if err := rows.Scan(
			&dt.DiscountTypeID, &dt.TypeName, &dt.DiscountUnitType,
			&dt.Description, &dt.UpdateBy, &dt.UpdateDate, &dt.IDStatus,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		list = append(list, dt)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_discount_type WHERE is_delete = 0`).Scan(&total)
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
			"page":         page,
			"limit":        limit,
			"total":        total,
			"discountType": list,
		},
	})
}

func SelectDiscountTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT discount_type_id, type_name, discount_unit_type, description,
		       update_by, update_date, id_status
		FROM tb_discount_type
		WHERE discount_type_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var dt models.DiscountType
	if err := row.Scan(
		&dt.DiscountTypeID, &dt.TypeName, &dt.DiscountUnitType,
		&dt.Description, &dt.UpdateBy, &dt.UpdateDate, &dt.IDStatus,
	); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Discount type not found",
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

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    dt,
	})
}

func SelectDiscountTypeByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT discount_type_id, type_name, discount_unit_type, description,
		       update_by, update_date, id_status
		FROM tb_discount_type
		WHERE type_name LIKE '%' + @Name + '%' AND is_delete = 0
	`

	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch discount types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var results []models.DiscountType
	for rows.Next() {
		var dt models.DiscountType
		if err := rows.Scan(
			&dt.DiscountTypeID, &dt.TypeName, &dt.DiscountUnitType,
			&dt.Description, &dt.UpdateBy, &dt.UpdateDate, &dt.IDStatus,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		results = append(results, dt)
	}

	if len(results) == 0 {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "No matching discount type found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    results,
	})
}

func InsertDiscountType(c *fiber.Ctx) error {
	var dt models.DiscountType
	if err := c.BodyParser(&dt); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request",
			Data:    nil,
		})
	}

	query := `
		INSERT INTO tb_discount_type (type_name, discount_unit_type, description, update_by)
		VALUES (@TypeName, @UnitType, @Description, @UpdateBy)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeName", dt.TypeName),
				sql.Named("UnitType", dt.DiscountUnitType),
				sql.Named("Description", dt.Description),
				sql.Named("UpdateBy", dt.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to insert discount type",
			Data:    nil,
		})
	}
	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Discount type added successfully",
		Data:    nil,
	})
}

func UpdateDiscountTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var dt models.DiscountType
	if err := c.BodyParser(&dt); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		UPDATE tb_discount_type
		SET type_name = COALESCE(NULLIF(@TypeName, ''), type_name),
			discount_unit_type = COALESCE(NULLIF(@UnitType, ''), discount_unit_type),
			description = @Description,
			update_by = @UpdateBy
		WHERE discount_type_id = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeName", dt.TypeName),
				sql.Named("UnitType", dt.DiscountUnitType),
				sql.Named("Description", dt.Description),
				sql.Named("UpdateBy", dt.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update discount type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Discount type updated successfully",
		Data:    nil,
	})
}

func DeleteDiscountTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_discount_type
		SET is_delete = 1,
			update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE discount_type_id = @ID
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
			Message: "Failed to soft delete discount type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Discount type marked as deleted",
		Data:    nil,
	})
}

func RemoveDiscountTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_discount_type WHERE discount_type_id = @ID`
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
			Message: "Failed to hard delete discount type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Discount type removed successfully",
		Data:    nil,
	})
}
