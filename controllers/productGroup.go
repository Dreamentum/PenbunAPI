package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllProductGroup(c *fiber.Ctx) error {
	query := `
		SELECT g.product_group_id, g.product_category_id, g.product_group_name, g.description, 
		       g.update_by, g.update_date, g.is_active,
		       c.category_name
		FROM tb_product_group g
		LEFT JOIN tb_product_category c ON g.product_category_id = c.product_category_id
		WHERE g.is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch product groups",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.ProductGroup
	for rows.Next() {
		var item models.ProductGroup
		var upd sql.NullTime
		if err := rows.Scan(&item.ProductGroupID, &item.ProductCategoryID, &item.ProductGroupName, &item.Description, 
			&item.UpdateBy, &upd, &item.IsActive, &item.CategoryName); err != nil {
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

func SelectPageProductGroup(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT g.product_group_id, g.product_category_id, g.product_group_name, g.description, 
		       g.update_by, g.update_date, g.is_active,
		       c.category_name
		FROM tb_product_group g
		LEFT JOIN tb_product_category c ON g.product_category_id = c.product_category_id
		WHERE g.is_delete = 0
		ORDER BY g.update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch product groups",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.ProductGroup
	for rows.Next() {
		var item models.ProductGroup
		var upd sql.NullTime
		if err := rows.Scan(&item.ProductGroupID, &item.ProductCategoryID, &item.ProductGroupName, &item.Description, 
			&item.UpdateBy, &upd, &item.IsActive, &item.CategoryName); err != nil {
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
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_product_group WHERE is_delete = 0`).Scan(&total)
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
			"productGroup": list,
		},
	})
}

func SelectProductGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT g.product_group_id, g.product_category_id, g.product_group_name, g.description, 
		       g.update_by, g.update_date, g.is_active,
		       c.category_name
		FROM tb_product_group g
		LEFT JOIN tb_product_category c ON g.product_category_id = c.product_category_id
		WHERE g.product_group_id = @ID AND g.is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var item models.ProductGroup
	var upd sql.NullTime
	if err := row.Scan(&item.ProductGroupID, &item.ProductCategoryID, &item.ProductGroupName, &item.Description, 
		&item.UpdateBy, &upd, &item.IsActive, &item.CategoryName); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Product group not found",
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

func SelectProductGroupByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT g.product_group_id, g.product_category_id, g.product_group_name, g.description, 
		       g.update_by, g.update_date, g.is_active,
		       c.category_name
		FROM tb_product_group g
		LEFT JOIN tb_product_category c ON g.product_category_id = c.product_category_id
		WHERE g.product_group_name LIKE '%' + @Name + '%' AND g.is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch product groups",
			Data:    nil,
		})
	}
	defer rows.Close()

	var results []models.ProductGroup
	for rows.Next() {
		var item models.ProductGroup
		var upd sql.NullTime
		if err := rows.Scan(&item.ProductGroupID, &item.ProductCategoryID, &item.ProductGroupName, &item.Description, 
			&item.UpdateBy, &upd, &item.IsActive, &item.CategoryName); err != nil {
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
			Message: "No matching product group found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    results,
	})
}

func InsertProductGroup(c *fiber.Ctx) error {
	var item models.ProductGroup
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	if item.ProductCategoryID == "" {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "product_category_id is required",
			Data:    nil,
		})
	}

	query := `
		INSERT INTO tb_product_group (product_category_id, product_group_name, description, update_by)
		VALUES (@CategoryID, @GroupName, @Description, @UpdateBy)
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("CategoryID", item.ProductCategoryID),
				sql.Named("GroupName", item.ProductGroupName),
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
			Message: "Failed to insert product group",
			Data:    nil,
		})
	}

	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product group added successfully",
		Data:    nil,
	})
}

func UpdateProductGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.ProductGroup
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		UPDATE tb_product_group
		SET product_category_id = COALESCE(NULLIF(@CategoryID, ''), product_category_id),
			product_group_name = COALESCE(NULLIF(@GroupName, ''), product_group_name),
			description = COALESCE(@Description, description),
			update_by = @UpdateBy,
			is_active = COALESCE(@IsActive, is_active)
		WHERE product_group_id = @ID AND is_delete = 0
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(query,
				sql.Named("CategoryID", item.ProductCategoryID),
				sql.Named("GroupName", item.ProductGroupName),
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
			Message: "Product group not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update product group",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product group updated successfully",
		Data:    nil,
	})
}

func DeleteProductGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_product_group
		SET is_delete = 1,
			is_active = 0,
			update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE product_group_id = @ID AND is_delete = 0
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
			Message: "Product group not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to delete product group",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product group deleted successfully",
		Data:    nil,
	})
}

func RemoveProductGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_product_group WHERE product_group_id = @ID`
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
			Message: "Product group not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to remove product group",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product group removed successfully",
		Data:    nil,
	})
}
