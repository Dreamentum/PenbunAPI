package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllProductFormatType(c *fiber.Ctx) error {
	query := `
		SELECT product_format_type_id, format_name, description, update_by, update_date, id_status, is_delete
		FROM tb_product_format_type
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch product format types"})
	}
	defer rows.Close()

	var result []models.ProductFormatType
	for rows.Next() {
		var ft models.ProductFormatType
		if err := rows.Scan(&ft.ProductFormatTypeID, &ft.FormatName, &ft.Description, &ft.UpdateBy, &ft.UpdateDate, &ft.IDStatus, &ft.IsDelete); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data"})
		}
		result = append(result, ft)
	}

	return c.JSON(models.ApiResponse{Status: "success", Data: result})
}

func SelectPageProductFormatType(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT product_format_type_id, format_name, description, update_by, update_date, id_status, is_delete
		FROM tb_product_format_type
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch product format types"})
	}
	defer rows.Close()

	var result []models.ProductFormatType
	for rows.Next() {
		var ft models.ProductFormatType
		if err := rows.Scan(&ft.ProductFormatTypeID, &ft.FormatName, &ft.Description, &ft.UpdateBy, &ft.UpdateDate, &ft.IDStatus, &ft.IsDelete); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data"})
		}
		result = append(result, ft)
	}

	var total int
	if err := config.DB.QueryRow(`SELECT COUNT(*) FROM tb_product_format_type WHERE is_delete = 0`).Scan(&total); err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to count records"})
	}

	return c.JSON(models.ApiResponse{
		Status: "success",
		Data:   fiber.Map{"page": page, "limit": limit, "total": total, "productFormatType": result},
	})
}

func SelectProductFormatTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT product_format_type_id, format_name, description, update_by, update_date, id_status, is_delete
		FROM tb_product_format_type
		WHERE product_format_type_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var ft models.ProductFormatType
	if err := row.Scan(&ft.ProductFormatTypeID, &ft.FormatName, &ft.Description, &ft.UpdateBy, &ft.UpdateDate, &ft.IDStatus, &ft.IsDelete); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "Product format type not found"})
		}
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data"})
	}

	return c.JSON(models.ApiResponse{Status: "success", Data: ft})
}

func SelectProductFormatTypeByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT product_format_type_id, format_name, description, update_by, update_date, id_status, is_delete
		FROM tb_product_format_type
		WHERE format_name LIKE '%' + @Name + '%' AND is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to search product format types"})
	}
	defer rows.Close()

	var result []models.ProductFormatType
	for rows.Next() {
		var ft models.ProductFormatType
		if err := rows.Scan(&ft.ProductFormatTypeID, &ft.FormatName, &ft.Description, &ft.UpdateBy, &ft.UpdateDate, &ft.IDStatus, &ft.IsDelete); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data"})
		}
		result = append(result, ft)
	}

	if len(result) == 0 {
		return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "No matching product format type found"})
	}

	return c.JSON(models.ApiResponse{Status: "success", Data: result})
}

func InsertProductFormatType(c *fiber.Ctx) error {
	var ft models.ProductFormatType
	if err := c.BodyParser(&ft); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid request body"})
	}

	query := `
		INSERT INTO tb_product_format_type (format_name, description, update_by)
		VALUES (@FormatName, @Description, @UpdateBy)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("FormatName", ft.FormatName),
				sql.Named("Description", ft.Description),
				sql.Named("UpdateBy", ft.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to insert product format type"})
	}

	return c.Status(201).JSON(models.ApiResponse{Status: "success", Message: "Product format type inserted successfully"})
}

func UpdateProductFormatTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var ft models.ProductFormatType
	if err := c.BodyParser(&ft); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid request body"})
	}

	query := `
		UPDATE tb_product_format_type
		SET format_name = COALESCE(NULLIF(@FormatName, ''), format_name),
		    description = @Description,
		    update_by = @UpdateBy
		WHERE product_format_type_id = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("FormatName", ft.FormatName),
				sql.Named("Description", ft.Description),
				sql.Named("UpdateBy", ft.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to update product format type"})
	}

	return c.JSON(models.ApiResponse{Status: "success", Message: "Product format type updated successfully"})
}

func DeleteProductFormatTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_product_format_type
		SET is_delete = 1,
		    update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE product_format_type_id = @ID
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, sql.Named("ID", id))
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to soft delete product format type"})
	}

	return c.JSON(models.ApiResponse{Status: "success", Message: "Product format type deleted (soft)"})
}

func RemoveProductFormatTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_product_format_type WHERE product_format_type_id = @ID`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, sql.Named("ID", id))
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to hard delete product format type"})
	}

	return c.JSON(models.ApiResponse{Status: "success", Message: "Product format type deleted (hard)"})
}
