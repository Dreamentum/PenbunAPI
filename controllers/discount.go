package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllDiscount(c *fiber.Ctx) error {
	query := `
		SELECT discount_id, discount_name, discount_type, discount_value,
		       start_date, end_date, note, update_by, update_date, id_status
		FROM tb_discount
		WHERE is_delete = 0
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
		var d models.Discount
		if err := rows.Scan(
			&d.DiscountID, &d.DiscountName, &d.DiscountType, &d.DiscountValue,
			&d.StartDate, &d.EndDate, &d.Note, &d.UpdateBy, &d.UpdateDate, &d.IDStatus,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		list = append(list, d)
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
		SELECT discount_id, discount_name, discount_type, discount_value,
		       start_date, end_date, note, update_by, update_date, id_status
		FROM tb_discount
		WHERE is_delete = 0
		ORDER BY update_date DESC
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
		var d models.Discount
		if err := rows.Scan(
			&d.DiscountID, &d.DiscountName, &d.DiscountType, &d.DiscountValue,
			&d.StartDate, &d.EndDate, &d.Note, &d.UpdateBy, &d.UpdateDate, &d.IDStatus,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		list = append(list, d)
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
		SELECT discount_id, discount_name, discount_type, discount_value,
		       start_date, end_date, note, update_by, update_date, id_status
		FROM tb_discount
		WHERE discount_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var d models.Discount
	if err := row.Scan(
		&d.DiscountID, &d.DiscountName, &d.DiscountType, &d.DiscountValue,
		&d.StartDate, &d.EndDate, &d.Note, &d.UpdateBy, &d.UpdateDate, &d.IDStatus,
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

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    d,
	})
}

func SelectDiscountByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT discount_id, discount_name, discount_type, discount_value,
		       start_date, end_date, note, update_by, update_date, id_status
		FROM tb_discount
		WHERE discount_name LIKE '%' + @Name + '%' AND is_delete = 0
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
		var d models.Discount
		if err := rows.Scan(
			&d.DiscountID, &d.DiscountName, &d.DiscountType, &d.DiscountValue,
			&d.StartDate, &d.EndDate, &d.Note, &d.UpdateBy, &d.UpdateDate, &d.IDStatus,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		list = append(list, d)
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
	var d models.Discount
	if err := c.BodyParser(&d); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}
	query := `
		INSERT INTO tb_discount (discount_name, discount_type, discount_value,
			start_date, end_date, note, update_by)
		VALUES (@Name, @Type, @Value, @Start, @End, @Note, @UpdateBy)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("Name", d.DiscountName),
				sql.Named("Type", d.DiscountType),
				sql.Named("Value", d.DiscountValue),
				sql.Named("Start", d.StartDate),
				sql.Named("End", d.EndDate),
				sql.Named("Note", d.Note),
				sql.Named("UpdateBy", d.UpdateBy),
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
	var d models.Discount
	if err := c.BodyParser(&d); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}
	query := `
		UPDATE tb_discount
		SET discount_name = COALESCE(NULLIF(@Name, ''), discount_name),
			discount_type = COALESCE(NULLIF(@Type, ''), discount_type),
			discount_value = @Value,
			start_date = @Start,
			end_date = @End,
			note = @Note,
			update_by = @UpdateBy
		WHERE discount_id = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("Name", d.DiscountName),
				sql.Named("Type", d.DiscountType),
				sql.Named("Value", d.DiscountValue),
				sql.Named("Start", d.StartDate),
				sql.Named("End", d.EndDate),
				sql.Named("Note", d.Note),
				sql.Named("UpdateBy", d.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
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
	query := `
		UPDATE tb_discount
		SET is_delete = 1,
		    update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE discount_id = @ID
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
			Message: "Failed to soft delete discount",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Discount marked as deleted",
		Data:    nil,
	})
}

func RemoveDiscountByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_discount WHERE discount_id = @ID`
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
			Message: "Failed to hard delete discount",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Discount removed successfully",
		Data:    nil,
	})
}
