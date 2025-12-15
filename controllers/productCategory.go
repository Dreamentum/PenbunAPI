package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllProductCategory(c *fiber.Ctx) error {
	query := `
		SELECT product_category_id, category_name, category_code, description, update_by, update_date, is_active, is_delete
		FROM tb_product_category
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch product categories"})
	}
	defer rows.Close()

	var result []models.ProductCategory
	for rows.Next() {
		var pc models.ProductCategory
		if err := rows.Scan(&pc.ProductCategoryID, &pc.CategoryName, &pc.CategoryCode, &pc.Description, &pc.UpdateBy, &pc.UpdateDate, &pc.IsActive, &pc.IsDelete); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data"})
		}
		result = append(result, pc)
	}

	return c.JSON(models.ApiResponse{Status: "success", Data: result})
}

func SelectPageProductCategory(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT product_category_id, category_name, category_code, description, update_by, update_date, is_active, is_delete
		FROM tb_product_category
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch product categories"})
	}
	defer rows.Close()

	var result []models.ProductCategory
	for rows.Next() {
		var pc models.ProductCategory
		if err := rows.Scan(&pc.ProductCategoryID, &pc.CategoryName, &pc.CategoryCode, &pc.Description, &pc.UpdateBy, &pc.UpdateDate, &pc.IsActive, &pc.IsDelete); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data"})
		}
		result = append(result, pc)
	}

	var total int
	if err := config.DB.QueryRow(`SELECT COUNT(*) FROM tb_product_category WHERE is_delete = 0`).Scan(&total); err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to count records"})
	}

	return c.JSON(models.ApiResponse{
		Status: "success",
		Data:   fiber.Map{"page": page, "limit": limit, "total": total, "productCategory": result},
	})
}

func SelectProductCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT product_category_id, category_name, category_code, description, update_by, update_date, is_active, is_delete
		FROM tb_product_category
		WHERE product_category_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var pc models.ProductCategory
	if err := row.Scan(&pc.ProductCategoryID, &pc.CategoryName, &pc.CategoryCode, &pc.Description, &pc.UpdateBy, &pc.UpdateDate, &pc.IsActive, &pc.IsDelete); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "Product category not found"})
		}
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data"})
	}

	return c.JSON(models.ApiResponse{Status: "success", Data: pc})
}

func SelectProductCategoryByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT product_category_id, category_name, category_code, description, update_by, update_date, is_active, is_delete
		FROM tb_product_category
		WHERE category_name LIKE '%' + @Name + '%' AND is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to search product categories"})
	}
	defer rows.Close()

	var result []models.ProductCategory
	for rows.Next() {
		var pc models.ProductCategory
		if err := rows.Scan(&pc.ProductCategoryID, &pc.CategoryName, &pc.CategoryCode, &pc.Description, &pc.UpdateBy, &pc.UpdateDate, &pc.IsActive, &pc.IsDelete); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data"})
		}
		result = append(result, pc)
	}

	if len(result) == 0 {
		return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "No matching product category found"})
	}

	return c.JSON(models.ApiResponse{Status: "success", Data: result})
}

func InsertProductCategory(c *fiber.Ctx) error {
	var pc models.ProductCategory
	if err := c.BodyParser(&pc); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid request body"})
	}

	query := `
		INSERT INTO tb_product_category (category_name, category_code, description, update_by)
		VALUES (@CategoryName, @CategoryCode, @Description, @UpdateBy)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("CategoryName", pc.CategoryName),
				sql.Named("CategoryCode", pc.CategoryCode),
				sql.Named("Description", pc.Description),
				sql.Named("UpdateBy", pc.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to insert product category"})
	}

	return c.Status(201).JSON(models.ApiResponse{Status: "success", Message: "Product category inserted successfully"})
}

func UpdateProductCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var pc models.ProductCategory
	if err := c.BodyParser(&pc); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid request body"})
	}

	query := `
		UPDATE tb_product_category
		SET category_name = COALESCE(NULLIF(@CategoryName, ''), category_name),
			category_code = COALESCE(NULLIF(@CategoryCode, ''), category_code),
		    description = COALESCE(@Description, description),
		    update_by = @UpdateBy,
		    is_active = COALESCE(@Status, is_active)
		WHERE product_category_id = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(query,
				sql.Named("CategoryName", pc.CategoryName),
				sql.Named("CategoryCode", pc.CategoryCode),
				sql.Named("Description", pc.Description),
				sql.Named("UpdateBy", pc.UpdateBy),
				sql.Named("Status", pc.IsActive),
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
				return errors.New("product category not found")
			}
			return nil
		},
	})
	if err != nil {
		if err.Error() == "product category not found" {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Product category not found",
				Data:    nil,
			})
		}
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to update product category"})
	}

	return c.JSON(models.ApiResponse{Status: "success", Message: "Product category updated successfully"})
}

func DeleteProductCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_product_category
		SET is_delete = 1,
		    update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE product_category_id = @ID
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, sql.Named("ID", id))
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to soft delete product category"})
	}

	return c.JSON(models.ApiResponse{Status: "success", Message: "Product category deleted (soft)"})
}

func RemoveProductCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_product_category WHERE product_category_id = @ID`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, sql.Named("ID", id))
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to hard delete product category"})
	}

	return c.JSON(models.ApiResponse{Status: "success", Message: "Product category deleted (hard)"})
}
