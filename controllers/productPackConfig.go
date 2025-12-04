package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

// 1. Select All
func SelectAllProductPackConfig(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT autoID, product_pack_config_id, product_id, bundle_qty, unit_type_id, note, update_by, update_date, id_status FROM tb_product_pack_config WHERE is_delete = 0")
	if err != nil {
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch product pack configs",
		})
	}
	defer rows.Close()

	var configs []models.ProductPackConfig
	for rows.Next() {
		var cfg models.ProductPackConfig
		if err := rows.Scan(&cfg.AutoID, &cfg.ProductPackConfigID, &cfg.ProductID, &cfg.BundleQty, &cfg.UnitTypeID, &cfg.Note, &cfg.UpdateBy, &cfg.UpdateDate, &cfg.IDStatus); err != nil {
			continue
		}
		configs = append(configs, cfg)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product pack configs fetched successfully",
		Data:    configs,
	})
}

// 2. Select Page
func SelectPageProductPackConfig(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM tb_product_pack_config WHERE is_delete = 0").Scan(&total)
	if err != nil {
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to count product pack configs",
		})
	}

	// Get data
	query := `
		SELECT autoID, product_pack_config_id, product_id, bundle_qty, unit_type_id, note, update_by, update_date, id_status 
		FROM tb_product_pack_config 
		WHERE is_delete = 0 
		ORDER BY update_date DESC 
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY`

	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch product pack configs page",
		})
	}
	defer rows.Close()

	var configs []models.ProductPackConfig
	for rows.Next() {
		var cfg models.ProductPackConfig
		if err := rows.Scan(&cfg.AutoID, &cfg.ProductPackConfigID, &cfg.ProductID, &cfg.BundleQty, &cfg.UnitTypeID, &cfg.Note, &cfg.UpdateBy, &cfg.UpdateDate, &cfg.IDStatus); err != nil {
			continue
		}
		configs = append(configs, cfg)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product pack configs page fetched successfully",
		Data: fiber.Map{
			"total": total,
			"items": configs,
		},
	})
}

// 3. Select By ID
func SelectProductPackConfigByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var cfg models.ProductPackConfig

	query := `
		SELECT autoID, product_pack_config_id, product_id, bundle_qty, unit_type_id, note, update_by, update_date, id_status 
		FROM tb_product_pack_config 
		WHERE product_pack_config_id = @ID AND is_delete = 0`

	err := config.DB.QueryRow(query, sql.Named("ID", id)).Scan(
		&cfg.AutoID, &cfg.ProductPackConfigID, &cfg.ProductID, &cfg.BundleQty, &cfg.UnitTypeID, &cfg.Note, &cfg.UpdateBy, &cfg.UpdateDate, &cfg.IDStatus,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "fail",
				Message: "Product pack config not found",
			})
		}
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch product pack config",
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product pack config fetched successfully",
		Data:    cfg,
	})
}

// 4. Select By Name (Search by Note or ProductID)
func SelectProductPackConfigByName(c *fiber.Ctx) error {
	name := c.Params("name") // Using 'name' param for search query
	query := `
		SELECT autoID, product_pack_config_id, product_id, bundle_qty, unit_type_id, note, update_by, update_date, id_status 
		FROM tb_product_pack_config 
		WHERE (note LIKE '%' + @Name + '%' OR product_id LIKE '%' + @Name + '%') AND is_delete = 0`

	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to search product pack configs",
		})
	}
	defer rows.Close()

	var configs []models.ProductPackConfig
	for rows.Next() {
		var cfg models.ProductPackConfig
		if err := rows.Scan(&cfg.AutoID, &cfg.ProductPackConfigID, &cfg.ProductID, &cfg.BundleQty, &cfg.UnitTypeID, &cfg.Note, &cfg.UpdateBy, &cfg.UpdateDate, &cfg.IDStatus); err != nil {
			continue
		}
		configs = append(configs, cfg)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product pack configs search results",
		Data:    configs,
	})
}

// 5. Insert
func InsertProductPackConfig(c *fiber.Ctx) error {
	var cfg models.ProductPackConfig
	if err := c.BodyParser(&cfg); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "fail",
			Message: "Invalid request body",
		})
	}

	query := `
		INSERT INTO tb_product_pack_config (product_id, bundle_qty, unit_type_id, note, update_by, id_status)
		VALUES (@ProductID, @BundleQty, @UnitTypeID, @Note, @UpdateBy, @IDStatus)`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("ProductID", cfg.ProductID),
				sql.Named("BundleQty", cfg.BundleQty),
				sql.Named("UnitTypeID", cfg.UnitTypeID),
				sql.Named("Note", cfg.Note),
				sql.Named("UpdateBy", cfg.UpdateBy),
				sql.Named("IDStatus", cfg.IDStatus),
			)
			return err
		},
	})

	if err != nil {
		log.Println("Insert Error:", err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to create product pack config",
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product pack config created successfully",
		Data:    cfg,
	})
}

// 6. Update By ID
func UpdateProductPackConfigByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var cfg models.ProductPackConfig
	if err := c.BodyParser(&cfg); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "fail",
			Message: "Invalid request body",
		})
	}

	query := `
		UPDATE tb_product_pack_config 
		SET product_id = @ProductID,
			bundle_qty = @BundleQty,
			unit_type_id = @UnitTypeID,
			note = @Note,
			update_by = @UpdateBy,
			id_status = @IDStatus
		WHERE product_pack_config_id = @ID AND is_delete = 0`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			result, err := tx.Exec(query,
				sql.Named("ProductID", cfg.ProductID),
				sql.Named("BundleQty", cfg.BundleQty),
				sql.Named("UnitTypeID", cfg.UnitTypeID),
				sql.Named("Note", cfg.Note),
				sql.Named("UpdateBy", cfg.UpdateBy),
				sql.Named("IDStatus", cfg.IDStatus),
				sql.Named("ID", id),
			)
			if err != nil {
				return err
			}
			rowsAffected, _ := result.RowsAffected()
			if rowsAffected == 0 {
				return fmt.Errorf("record not found")
			}
			return nil
		},
	})

	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "fail",
				Message: "Product pack config not found",
			})
		}
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update product pack config",
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product pack config updated successfully",
		Data:    cfg,
	})
}

// 7. Delete By ID (Soft Delete)
func DeleteProductPackConfigByID(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Query("user")
	if username == "" {
		username = "UNKNOWN"
	}

	query := `UPDATE tb_product_pack_config SET is_delete = 1, update_by = @UpdateBy WHERE product_pack_config_id = @ID`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			result, err := tx.Exec(query, sql.Named("UpdateBy", username), sql.Named("ID", id))
			if err != nil {
				return err
			}
			rowsAffected, _ := result.RowsAffected()
			if rowsAffected == 0 {
				return fmt.Errorf("record not found")
			}
			return nil
		},
	})

	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "fail",
				Message: "Product pack config not found",
			})
		}
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to delete product pack config",
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product pack config deleted successfully",
	})
}

// 8. Remove By ID (Hard Delete)
func RemoveProductPackConfigByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_product_pack_config WHERE product_pack_config_id = @ID`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			result, err := tx.Exec(query, sql.Named("ID", id))
			if err != nil {
				return err
			}
			rowsAffected, _ := result.RowsAffected()
			if rowsAffected == 0 {
				return fmt.Errorf("record not found")
			}
			return nil
		},
	})

	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "fail",
				Message: "Product pack config not found",
			})
		}
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to remove product pack config",
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product pack config removed successfully",
	})
}
