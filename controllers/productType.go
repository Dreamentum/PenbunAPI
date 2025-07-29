package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllProductType(c *fiber.Ctx) error {
	query := `
		SELECT product_type_id, type_name, type_group_name, description, update_by, update_date, id_status
		FROM tb_product_type
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch product types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var result []models.ProductType
	for rows.Next() {
		var pt models.ProductType
		if err := rows.Scan(&pt.ProductTypeID, &pt.TypeName, &pt.TypeGroupName, &pt.Description, &pt.UpdateBy, &pt.UpdateDate, &pt.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read product types",
				Data:    nil,
			})
		}
		result = append(result, pt)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    result,
	})
}

func SelectPageProductType(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT product_type_id, type_name, type_group_name, description, update_by, update_date, id_status
		FROM tb_product_type
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch product types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var result []models.ProductType
	for rows.Next() {
		var pt models.ProductType
		if err := rows.Scan(&pt.ProductTypeID, &pt.TypeName, &pt.TypeGroupName, &pt.Description, &pt.UpdateBy, &pt.UpdateDate, &pt.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read product types",
				Data:    nil,
			})
		}
		result = append(result, pt)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_product_type WHERE is_delete = 0`).Scan(&total)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to count product types",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status: "success",
		Data: fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"productType": result,
		},
	})
}

func SelectProductTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT product_type_id, type_name, type_group_name, description, update_by, update_date, id_status
		FROM tb_product_type
		WHERE product_type_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var pt models.ProductType
	if err := row.Scan(&pt.ProductTypeID, &pt.TypeName, &pt.TypeGroupName, &pt.Description, &pt.UpdateBy, &pt.UpdateDate, &pt.IDStatus); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Product type not found",
				Data:    nil,
			})
		}
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to read product type",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    pt,
	})
}

func SelectProductTypeByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT product_type_id, type_name, type_group_name, description, update_by, update_date, id_status
		FROM tb_product_type
		WHERE type_name LIKE '%' + @Name + '%' AND is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to search product types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var result []models.ProductType
	for rows.Next() {
		var pt models.ProductType
		if err := rows.Scan(&pt.ProductTypeID, &pt.TypeName, &pt.TypeGroupName, &pt.Description, &pt.UpdateBy, &pt.UpdateDate, &pt.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read product types",
				Data:    nil,
			})
		}
		result = append(result, pt)
	}

	if len(result) == 0 {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "No matching product type found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    result,
	})
}

func InsertProductType(c *fiber.Ctx) error {
	var pt models.ProductType
	if err := c.BodyParser(&pt); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		INSERT INTO tb_product_type (type_name, type_group_name, description, update_by)
		VALUES (@TypeName, @TypeGroupName, @Description, @UpdateBy)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeName", pt.TypeName),
				sql.Named("TypeGroupName", pt.TypeGroupName),
				sql.Named("Description", pt.Description),
				sql.Named("UpdateBy", pt.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to insert product type",
			Data:    nil,
		})
	}

	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product type inserted successfully",
		Data:    nil,
	})
}

func UpdateProductTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var pt models.ProductType
	if err := c.BodyParser(&pt); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		UPDATE tb_product_type
		SET type_name = COALESCE(NULLIF(@TypeName, ''), type_name),
			type_group_name = @TypeGroupName,
			description = @Description,
			update_by = @UpdateBy
		WHERE product_type_id = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeName", pt.TypeName),
				sql.Named("TypeGroupName", pt.TypeGroupName),
				sql.Named("Description", pt.Description),
				sql.Named("UpdateBy", pt.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update product type",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product type updated successfully",
		Data:    nil,
	})
}

func DeleteProductTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_product_type
		SET is_delete = 1,
			update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE product_type_id = @ID
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
			Message: "Failed to soft delete product type",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product type deleted (soft)",
		Data:    nil,
	})
}

func RemoveProductTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_product_type WHERE product_type_id = @ID`
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
			Message: "Failed to hard delete product type",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product type deleted (hard)",
		Data:    nil,
	})
}
