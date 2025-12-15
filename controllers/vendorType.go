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

// ---------- 1) Select All ----------
func SelectAllVendorType(c *fiber.Ctx) error {
	query := `
        SELECT vendor_type_id, prefix, type_name, description, update_by, update_date, is_active, is_delete
        FROM tb_vendor_type
        WHERE is_delete = 0
        ORDER BY update_date DESC, vendor_type_id ASC
    `
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println("[SelectAllVendorType] query:", err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch vendor types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var result []models.VendorType
	for rows.Next() {
		var vt models.VendorType
		var upd sql.NullTime
		if err := rows.Scan(&vt.VendorTypeID, &vt.Prefix, &vt.TypeName, &vt.Description, &vt.UpdateBy, &upd, &vt.IsActive, &vt.IsDelete); err != nil {
			log.Println("[SelectAllVendorType] scan:", err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		if upd.Valid {
			t := upd.Time.Format("2006-01-02T15:04:05")
			vt.UpdateDate = &t
		}
		result = append(result, vt)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    result,
	})
}

// ---------- 2) Select Paging ----------
func SelectPageVendorType(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
        SELECT vendor_type_id, prefix, type_name, description, update_by, update_date, is_active, is_delete
        FROM tb_vendor_type
        WHERE is_delete = 0
        ORDER BY update_date DESC, vendor_type_id ASC
        OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
    `
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println("[SelectPageVendorType] query:", err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch vendor types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var items []models.VendorType
	for rows.Next() {
		var vt models.VendorType
		var upd sql.NullTime
		if err := rows.Scan(&vt.VendorTypeID, &vt.Prefix, &vt.TypeName, &vt.Description, &vt.UpdateBy, &upd, &vt.IsActive, &vt.IsDelete); err != nil {
			log.Println("[SelectPageVendorType] scan:", err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		if upd.Valid {
			t := upd.Time.Format("2006-01-02T15:04:05")
			vt.UpdateDate = &t
		}
		items = append(items, vt)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_vendor_type WHERE is_delete = 0`).Scan(&total)
	if err != nil {
		log.Println("[SelectPageVendorType] count:", err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to count records",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data: fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
			"items": items,
		},
	})
}

// ---------- 3) Select By ID ----------
func SelectVendorTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
        SELECT vendor_type_id, prefix, type_name, description, update_by, update_date, is_active, is_delete
        FROM tb_vendor_type
        WHERE vendor_type_id = @ID AND is_delete = 0
    `
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var vt models.VendorType
	var upd sql.NullTime
	if err := row.Scan(&vt.VendorTypeID, &vt.Prefix, &vt.TypeName, &vt.Description, &vt.UpdateBy, &upd, &vt.IsActive, &vt.IsDelete); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Vendor type not found",
				Data:    nil,
			})
		}
		log.Println("[SelectVendorTypeByID] scan:", err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to read data",
			Data:    nil,
		})
	}
	if upd.Valid {
		t := upd.Time.Format("2006-01-02T15:04:05")
		vt.UpdateDate = &t
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    vt,
	})
}

// ---------- 4) Select By Name (LIKE) ----------
func SelectVendorTypeByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
        SELECT vendor_type_id, prefix, type_name, description, update_by, update_date, is_active, is_delete
        FROM tb_vendor_type
        WHERE type_name LIKE '%' + @Name + '%' AND is_delete = 0
        ORDER BY type_name ASC, vendor_type_id ASC
    `
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println("[SelectVendorTypeByName] query:", err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch vendor types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var result []models.VendorType
	for rows.Next() {
		var vt models.VendorType
		var upd sql.NullTime
		if err := rows.Scan(&vt.VendorTypeID, &vt.Prefix, &vt.TypeName, &vt.Description, &vt.UpdateBy, &upd, &vt.IsActive, &vt.IsDelete); err != nil {
			log.Println("[SelectVendorTypeByName] scan:", err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		if upd.Valid {
			t := upd.Time.Format("2006-01-02T15:04:05")
			vt.UpdateDate = &t
		}
		result = append(result, vt)
	}

	if len(result) == 0 {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "No matching vendor type found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    result,
	})
}

// ---------- 5) Insert ----------
func InsertVendorType(c *fiber.Ctx) error {
	var vt models.VendorType
	if err := c.BodyParser(&vt); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	// Validation
	if strings.TrimSpace(vt.TypeName) == "" {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "type_name is required",
			Data:    nil,
		})
	}

	// Resolve update_by
	if vt.UpdateBy == nil || *vt.UpdateBy == "" {
		username := utils.ResolveUser(c)
		vt.UpdateBy = &username
	}

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			// Log incoming request for debugging
			log.Printf("[InsertVendorType] TypeName: %s, IsActive: %v, UpdateBy: %v", vt.TypeName, vt.IsActive, vt.UpdateBy)

			_, err := tx.Exec(`
                INSERT INTO tb_vendor_type (type_name, description, is_active, update_by)
                VALUES (@TypeName, @Description, @IsActive, @UpdateBy)`,
				sql.Named("TypeName", vt.TypeName),
				sql.Named("Description", vt.Description),
				sql.Named("IsActive", vt.IsActive),
				sql.Named("UpdateBy", vt.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println("[InsertVendorType] exec:", err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to insert vendor type",
			Data:    fiber.Map{"type_name": vt.TypeName},
		})
	}

	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor type added successfully",
		Data:    fiber.Map{"type_name": vt.TypeName},
	})
}

// ---------- 6) Update By ID ----------
func UpdateVendorTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var vt models.VendorType
	if err := c.BodyParser(&vt); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	// Resolve update_by
	if vt.UpdateBy == nil || *vt.UpdateBy == "" {
		username := utils.ResolveUser(c)
		vt.UpdateBy = &username
	}

	// Log incoming request for debugging
	log.Printf("[UpdateVendorTypeByID] ID: %s, TypeName: %s, IsActive: %v", id, vt.TypeName, vt.IsActive)

	query := `
        UPDATE tb_vendor_type
        SET type_name = COALESCE(NULLIF(@TypeName, ''), type_name),
            description = COALESCE(@Description, description),
            is_active = COALESCE(@IsActive, is_active),
            update_by = @UpdateBy
        WHERE vendor_type_id = @ID AND is_delete = 0
    `
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(query,
				sql.Named("TypeName", vt.TypeName),
				sql.Named("Description", vt.Description),
				sql.Named("IsActive", vt.IsActive),
				sql.Named("UpdateBy", vt.UpdateBy),
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
				return sql.ErrNoRows // Use standard error for 404
			}
			return nil
		},
	})
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Vendor type not found",
			Data:    fiber.Map{"vendor_type_id": id},
		})
	}
	if err != nil {
		log.Println("[UpdateVendorTypeByID] exec:", err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update vendor type",
			Data:    fiber.Map{"vendor_type_id": id},
		})
	}

	log.Printf("[UpdateVendorTypeByID] Successfully updated ID: %s with status: %v", id, vt.IsActive)

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor type updated successfully",
		Data:    fiber.Map{"vendor_type_id": id},
	})
}

// ---------- 7) Delete (Soft) ----------
func DeleteVendorTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	username := utils.ResolveUser(c)

	query := `
        UPDATE tb_vendor_type
        SET is_delete = 1, update_by = @UpdateBy
        WHERE vendor_type_id = @ID
    `
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(query,
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
			Message: "Vendor type not found",
			Data:    fiber.Map{"vendor_type_id": id},
		})
	}
	if err != nil {
		log.Println("[DeleteVendorTypeByID] exec:", err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to soft delete vendor type",
			Data:    fiber.Map{"vendor_type_id": id},
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor type marked as deleted",
		Data:    fiber.Map{"vendor_type_id": id, "update_by": username},
	})
}

// ---------- 8) Remove (Hard) ----------
func RemoveVendorTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_vendor_type WHERE vendor_type_id = @ID`
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
			Message: "Vendor type not found",
			Data:    fiber.Map{"vendor_type_id": id},
		})
	}
	if err != nil {
		log.Println("[RemoveVendorTypeByID] exec:", err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to hard delete vendor type",
			Data:    fiber.Map{"vendor_type_id": id},
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor type removed successfully",
		Data:    fiber.Map{"vendor_type_id": id},
	})
}