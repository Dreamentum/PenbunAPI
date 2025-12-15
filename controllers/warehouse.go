package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllWarehouse(c *fiber.Ctx) error {
	query := `
		SELECT warehouse_id, warehouse_code, warehouse_name, description, is_main_dc, allow_negative_stock, update_by, update_date, is_active
		FROM tb_warehouse
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch warehouses",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.Warehouse
	for rows.Next() {
		var item models.Warehouse
		var upd sql.NullTime
		if err := rows.Scan(&item.WarehouseID, &item.WarehouseCode, &item.WarehouseName, &item.Description, &item.IsMainDC, &item.AllowNegativeStock, &item.UpdateBy, &upd, &item.IsActive); err != nil {
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

func SelectPageWarehouse(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT warehouse_id, warehouse_code, warehouse_name, description, is_main_dc, allow_negative_stock, update_by, update_date, is_active
		FROM tb_warehouse
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch warehouses",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.Warehouse
	for rows.Next() {
		var item models.Warehouse
		var upd sql.NullTime
		if err := rows.Scan(&item.WarehouseID, &item.WarehouseCode, &item.WarehouseName, &item.Description, &item.IsMainDC, &item.AllowNegativeStock, &item.UpdateBy, &upd, &item.IsActive); err != nil {
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
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_warehouse WHERE is_delete = 0`).Scan(&total)
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
			"page":      page,
			"limit":     limit,
			"total":     total,
			"warehouse": list,
		},
	})
}

func SelectWarehouseByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT warehouse_id, warehouse_code, warehouse_name, description, is_main_dc, allow_negative_stock, update_by, update_date, is_active
		FROM tb_warehouse
		WHERE warehouse_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var item models.Warehouse
	var upd sql.NullTime
	if err := row.Scan(&item.WarehouseID, &item.WarehouseCode, &item.WarehouseName, &item.Description, &item.IsMainDC, &item.AllowNegativeStock, &item.UpdateBy, &upd, &item.IsActive); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Warehouse not found",
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

func SelectWarehouseByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT warehouse_id, warehouse_code, warehouse_name, description, is_main_dc, allow_negative_stock, update_by, update_date, is_active
		FROM tb_warehouse
		WHERE warehouse_name LIKE '%' + @Name + '%' AND is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch warehouses",
			Data:    nil,
		})
	}
	defer rows.Close()

	var results []models.Warehouse
	for rows.Next() {
		var item models.Warehouse
		var upd sql.NullTime
		if err := rows.Scan(&item.WarehouseID, &item.WarehouseCode, &item.WarehouseName, &item.Description, &item.IsMainDC, &item.AllowNegativeStock, &item.UpdateBy, &upd, &item.IsActive); err != nil {
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
			Message: "No matching warehouse found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    results,
	})
}

func InsertWarehouse(c *fiber.Ctx) error {
	var item models.Warehouse
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		INSERT INTO tb_warehouse (warehouse_code, warehouse_name, description, is_main_dc, allow_negative_stock, update_by)
		VALUES (@WarehouseCode, @WarehouseName, @Description, COALESCE(@IsMainDC, 0), COALESCE(@AllowNegativeStock, 0), @UpdateBy)
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("WarehouseCode", item.WarehouseCode),
				sql.Named("WarehouseName", item.WarehouseName),
				sql.Named("Description", item.Description),
				sql.Named("IsMainDC", item.IsMainDC),
				sql.Named("AllowNegativeStock", item.AllowNegativeStock),
				sql.Named("UpdateBy", item.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to insert warehouse",
			Data:    nil,
		})
	}

	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Warehouse added successfully",
		Data:    nil,
	})
}

func UpdateWarehouseByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.Warehouse
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		UPDATE tb_warehouse
		SET warehouse_code = COALESCE(NULLIF(@WarehouseCode, ''), warehouse_code),
			warehouse_name = COALESCE(NULLIF(@WarehouseName, ''), warehouse_name),
			description = COALESCE(@Description, description),
			is_main_dc = COALESCE(@IsMainDC, is_main_dc),
			allow_negative_stock = COALESCE(@AllowNegativeStock, allow_negative_stock),
			update_by = @UpdateBy,
			is_active = COALESCE(@IsActive, is_active)
		WHERE warehouse_id = @ID AND is_delete = 0
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(query,
				sql.Named("WarehouseCode", item.WarehouseCode),
				sql.Named("WarehouseName", item.WarehouseName),
				sql.Named("Description", item.Description),
				sql.Named("IsMainDC", item.IsMainDC),
				sql.Named("AllowNegativeStock", item.AllowNegativeStock),
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
			Message: "Warehouse not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update warehouse",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Warehouse updated successfully",
		Data:    nil,
	})
}

func DeleteWarehouseByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_warehouse
		SET is_delete = 1,
			is_active = 0,
			update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE warehouse_id = @ID AND is_delete = 0
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
			Message: "Warehouse not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to delete warehouse",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Warehouse deleted successfully",
		Data:    nil,
	})
}

func RemoveWarehouseByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_warehouse WHERE warehouse_id = @ID`
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
			Message: "Warehouse not found",
			Data:    nil,
		})
	}
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to remove warehouse",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Warehouse removed successfully",
		Data:    nil,
	})
}
