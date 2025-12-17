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
		SELECT discount_type_id, discount_type_name, description, update_by, update_date, is_active
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
		var item models.DiscountType
		var upd sql.NullTime
		if err := rows.Scan(&item.DiscountTypeID, &item.DiscountTypeName, &item.Description, &item.UpdateBy, &upd, &item.IsActive); err != nil {
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

func SelectPageDiscountType(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT discount_type_id, discount_type_name, description, update_by, update_date, is_active
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
		var item models.DiscountType
		var upd sql.NullTime
		if err := rows.Scan(&item.DiscountTypeID, &item.DiscountTypeName, &item.Description, &item.UpdateBy, &upd, &item.IsActive); err != nil {
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
			"items":        list,
		},
	})
}

func SelectDiscountTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT discount_type_id, discount_type_name, description, update_by, update_date, is_active
		FROM tb_discount_type
		WHERE discount_type_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var item models.DiscountType
	var upd sql.NullTime
	if err := row.Scan(&item.DiscountTypeID, &item.DiscountTypeName, &item.Description, &item.UpdateBy, &upd, &item.IsActive); err != nil {
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

func SelectDiscountTypeByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT discount_type_id, discount_type_name, description, update_by, update_date, is_active
		FROM tb_discount_type
		WHERE discount_type_name LIKE '%' + @Name + '%' AND is_delete = 0
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
		var item models.DiscountType
		var upd sql.NullTime
		if err := rows.Scan(&item.DiscountTypeID, &item.DiscountTypeName, &item.Description, &item.UpdateBy, &upd, &item.IsActive); err != nil {
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
		results = append(results, item)
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
	var item models.DiscountType
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		INSERT INTO tb_discount_type (discount_type_name, description, update_by)
		VALUES (@DiscountTypeName, @Description, @UpdateBy)
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("DiscountTypeName", item.DiscountTypeName),
				sql.Named("Description", item.Description),
				sql.Named("UpdateBy", item.UpdateBy),
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
	var item models.DiscountType
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		UPDATE tb_discount_type
		SET discount_type_name = COALESCE(NULLIF(@DiscountTypeName, ''), discount_type_name),
			description = COALESCE(@Description, description),
			update_by = @UpdateBy,
			is_active = COALESCE(@IsActive, is_active)
		WHERE discount_type_id = @ID AND is_delete = 0
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(query,
				sql.Named("DiscountTypeName", item.DiscountTypeName),
				sql.Named("Description", item.Description),
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
			Message: "Discount type not found",
			Data:    nil,
		})
	}
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
			is_active = 0,
			update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE discount_type_id = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(query, sql.Named("ID", id))
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
			Message: "Discount type not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to delete discount type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Discount type deleted successfully",
		Data:    nil,
	})
}

func RemoveDiscountTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_discount_type WHERE discount_type_id = @ID`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(query, sql.Named("ID", id))
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
			Message: "Discount type not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to remove discount type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Discount type removed successfully",
		Data:    nil,
	})
}
