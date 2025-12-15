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

func SelectAllDiscount(c *fiber.Ctx) error {
	query := `
		SELECT d.discount_id, d.discount_type_id, d.discount_name, d.discount_code, d.description,
		       d.discount_value, d.is_percent, d.min_order_amount,
		       d.start_date, d.end_date,
		       d.update_by, d.update_date, d.is_active,
		       dt.discount_type_name
		FROM tb_discount d
		LEFT JOIN tb_discount_type dt ON d.discount_type_id = dt.discount_type_id
		WHERE d.is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch discounts",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.Discount
	for rows.Next() {
		var item models.Discount
		var upd sql.NullTime
		var start sql.NullTime
		var end sql.NullTime
		if err := rows.Scan(
			&item.DiscountID, &item.DiscountTypeID, &item.DiscountName, &item.DiscountCode, &item.Description,
			&item.DiscountValue, &item.IsPercent, &item.MinOrderAmount,
			&start, &end,
			&item.UpdateBy, &upd, &item.IsActive, &item.DiscountTypeName,
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
		if start.Valid {
			t := start.Time.Format("2006-01-02T15:04:05")
			item.StartDate = &t
		}
		if end.Valid {
			t := end.Time.Format("2006-01-02T15:04:05")
			item.EndDate = &t
		}
		list = append(list, item)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    list,
	})
}

func SelectPageDiscount(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT d.discount_id, d.discount_type_id, d.discount_name, d.discount_code, d.description,
		       d.discount_value, d.is_percent, d.min_order_amount,
		       d.start_date, d.end_date,
		       d.update_by, d.update_date, d.is_active,
		       dt.discount_type_name
		FROM tb_discount d
		LEFT JOIN tb_discount_type dt ON d.discount_type_id = dt.discount_type_id
		WHERE d.is_delete = 0
		ORDER BY d.update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch discounts",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.Discount
	for rows.Next() {
		var item models.Discount
		var upd sql.NullTime
		var start sql.NullTime
		var end sql.NullTime
		if err := rows.Scan(
			&item.DiscountID, &item.DiscountTypeID, &item.DiscountName, &item.DiscountCode, &item.Description,
			&item.DiscountValue, &item.IsPercent, &item.MinOrderAmount,
			&start, &end,
			&item.UpdateBy, &upd, &item.IsActive, &item.DiscountTypeName,
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
		if start.Valid {
			t := start.Time.Format("2006-01-02T15:04:05")
			item.StartDate = &t
		}
		if end.Valid {
			t := end.Time.Format("2006-01-02T15:04:05")
			item.EndDate = &t
		}
		list = append(list, item)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_discount WHERE is_delete = 0`).Scan(&total)
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
			"discount": list,
		},
	})
}

func SelectDiscountByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT d.discount_id, d.discount_type_id, d.discount_name, d.discount_code, d.description,
		       d.discount_value, d.is_percent, d.min_order_amount,
		       d.start_date, d.end_date,
		       d.update_by, d.update_date, d.is_active,
		       dt.discount_type_name
		FROM tb_discount d
		LEFT JOIN tb_discount_type dt ON d.discount_type_id = dt.discount_type_id
		WHERE d.discount_id = @ID AND d.is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var item models.Discount
	var upd sql.NullTime
	var start sql.NullTime
	var end sql.NullTime
	if err := row.Scan(
		&item.DiscountID, &item.DiscountTypeID, &item.DiscountName, &item.DiscountCode, &item.Description,
		&item.DiscountValue, &item.IsPercent, &item.MinOrderAmount,
		&start, &end,
		&item.UpdateBy, &upd, &item.IsActive, &item.DiscountTypeName,
	); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Discount not found",
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
	if start.Valid {
		t := start.Time.Format("2006-01-02T15:04:05")
		item.StartDate = &t
	}
	if end.Valid {
		t := end.Time.Format("2006-01-02T15:04:05")
		item.EndDate = &t
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    item,
	})
}

func SelectDiscountByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT d.discount_id, d.discount_type_id, d.discount_name, d.discount_code, d.description,
		       d.discount_value, d.is_percent, d.min_order_amount,
		       d.start_date, d.end_date,
		       d.update_by, d.update_date, d.is_active,
		       dt.discount_type_name
		FROM tb_discount d
		LEFT JOIN tb_discount_type dt ON d.discount_type_id = dt.discount_type_id
		WHERE d.discount_name LIKE '%' + @Name + '%' AND d.is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch discounts",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.Discount
	for rows.Next() {
		var item models.Discount
		var upd sql.NullTime
		var start sql.NullTime
		var end sql.NullTime
		if err := rows.Scan(
			&item.DiscountID, &item.DiscountTypeID, &item.DiscountName, &item.DiscountCode, &item.Description,
			&item.DiscountValue, &item.IsPercent, &item.MinOrderAmount,
			&start, &end,
			&item.UpdateBy, &upd, &item.IsActive, &item.DiscountTypeName,
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
		if start.Valid {
			t := start.Time.Format("2006-01-02T15:04:05")
			item.StartDate = &t
		}
		if end.Valid {
			t := end.Time.Format("2006-01-02T15:04:05")
			item.EndDate = &t
		}
		list = append(list, item)
	}

	if len(list) == 0 {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "No matching discount found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    list,
	})
}

func InsertDiscount(c *fiber.Ctx) error {
	var item models.Discount
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	if item.DiscountTypeID == "" {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "discount_type_id is required"})
	}
	if strings.TrimSpace(item.DiscountName) == "" {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "discount_name is required"})
	}

	if item.UpdateBy == nil || *item.UpdateBy == "" {
		u := utils.ResolveUser(c)
		item.UpdateBy = &u
	}

	query := `
		INSERT INTO tb_discount (
			discount_type_id, discount_name, discount_code, description,
			discount_value, is_percent, min_order_amount,
			start_date, end_date, update_by
		)
		VALUES (
			@TypeID, @Name, @Code, @Desc,
			@Value, @IsPercent, @MinOrder,
			@StartDate, @EndDate, @UpdateBy
		)
	`
	
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeID", item.DiscountTypeID),
				sql.Named("Name", item.DiscountName),
				sql.Named("Code", item.DiscountCode),
				sql.Named("Desc", item.Description),
				sql.Named("Value", item.DiscountValue),
				sql.Named("IsPercent", item.IsPercent),
				sql.Named("MinOrder", item.MinOrderAmount),
				sql.Named("StartDate", item.StartDate), // Driver converts string to datetime if format is correct
				sql.Named("EndDate", item.EndDate),
				sql.Named("UpdateBy", item.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to insert discount",
			Data:    nil,
		})
	}

	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Discount added successfully",
		Data:    nil,
	})
}

func UpdateDiscountByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.Discount
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
		UPDATE tb_discount
		SET discount_type_id = COALESCE(NULLIF(@TypeID, ''), discount_type_id),
			discount_name = COALESCE(NULLIF(@Name, ''), discount_name),
			discount_code = COALESCE(@Code, discount_code),
			description = COALESCE(@Desc, description),
			discount_value = COALESCE(@Value, discount_value),
			is_percent = COALESCE(@IsPercent, is_percent),
			min_order_amount = COALESCE(@MinOrder, min_order_amount),
			start_date = COALESCE(@StartDate, start_date),
			end_date = COALESCE(@EndDate, end_date),
			update_by = @UpdateBy,
			is_active = COALESCE(@IsActive, is_active)
		WHERE discount_id = @ID AND is_delete = 0
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(query,
				sql.Named("TypeID", item.DiscountTypeID),
				sql.Named("Name", item.DiscountName),
				sql.Named("Code", item.DiscountCode),
				sql.Named("Desc", item.Description),
				sql.Named("Value", item.DiscountValue),
				sql.Named("IsPercent", item.IsPercent),
				sql.Named("MinOrder", item.MinOrderAmount),
				sql.Named("StartDate", item.StartDate),
				sql.Named("EndDate", item.EndDate),
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
			Message: "Discount not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update discount",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Discount updated successfully",
		Data:    nil,
	})
}

func DeleteDiscountByID(c *fiber.Ctx) error {
	id := c.Params("id")
	username := utils.ResolveUser(c)

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(`
				UPDATE tb_discount
				SET is_delete = 1,
					is_active = 0,
					update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME),
					update_by = @UpdateBy
				WHERE discount_id = @ID AND is_delete = 0`,
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
			Message: "Discount not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to delete discount",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Discount deleted successfully",
		Data:    nil,
	})
}

func RemoveDiscountByID(c *fiber.Ctx) error {
	id := c.Params("id")
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(`DELETE FROM tb_discount WHERE discount_id = @ID`, sql.Named("ID", id))
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
			Message: "Discount not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to remove discount",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Discount removed successfully",
		Data:    nil,
	})
}
