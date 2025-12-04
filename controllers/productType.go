package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ---------- Helpers ----------
func resolveUser(c *fiber.Ctx) string {
	u := strings.TrimSpace(c.Query("user"))
	if u == "" {
		u = "UNKNOWN"
	}
	return u
}

// ---------- 1) Select All (JOIN group) ----------
func SelectAllProductType(c *fiber.Ctx) error {
	rows, err := config.DB.Query(`
		SELECT  p.product_type_id,
				p.type_name,
				p.product_type_group_id,
				ISNULL(g.group_id, '') AS group_id,
				ISNULL(g.group_name, '') AS group_name,
				p.description,
				p.update_by,
				p.update_date,
				p.id_status,
				p.is_delete
		FROM tb_product_type AS p
		LEFT JOIN tb_product_type_group AS g
			ON g.group_id = p.product_type_group_id
		   AND g.is_delete = 0
		WHERE p.is_delete = 0
		ORDER BY p.update_date DESC, p.product_type_id ASC`)
	if err != nil {
		log.Println("[SelectAllProductType] query:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch", Data: nil})
	}
	defer rows.Close()

	var out []models.ProductType
	for rows.Next() {
		var m models.ProductType
		var upd sql.NullTime
		if err := rows.Scan(
			&m.ProductTypeID,
			&m.TypeName,
			&m.ProductTypeGroupID,
			&m.GroupId,
			&m.GroupName,
			&m.Description,
			&m.UpdateBy,
			&upd,
			&m.IDStatus,
			&m.IsDelete,
		); err != nil {
			log.Println("[SelectAllProductType] scan:", err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read", Data: nil})
		}
		if upd.Valid {
			t := upd.Time
			m.UpdateDate = &t
		}
		out = append(out, m)
	}

	return c.JSON(models.ApiResponse{Status: "success", Message: "", Data: out})
}

// ---------- 2) Select Paging (JOIN group) ----------
func SelectPageProductType(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	rows, err := config.DB.Query(`
		SELECT  p.product_type_id,
				p.type_name,
				p.product_type_group_id,
				ISNULL(g.group_id, '') AS group_id,
				ISNULL(g.group_name, '') AS group_name,
				p.description,
				p.update_by,
				p.update_date,
				p.id_status,
				p.is_delete
		FROM tb_product_type AS p
		LEFT JOIN tb_product_type_group AS g
			ON g.group_id = p.product_type_group_id
		   AND g.is_delete = 0
		WHERE p.is_delete = 0
		ORDER BY p.update_date DESC, p.product_type_id ASC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY`,
		sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println("[SelectPageProductType] query:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch", Data: nil})
	}
	defer rows.Close()

	items := []models.ProductType{} // Initialize as empty slice
	for rows.Next() {
		var m models.ProductType
		var upd sql.NullTime
		if err := rows.Scan(
			&m.ProductTypeID,
			&m.TypeName,
			&m.ProductTypeGroupID,
			&m.GroupId,
			&m.GroupName,
			&m.Description,
			&m.UpdateBy,
			&upd,
			&m.IDStatus,
			&m.IsDelete,
		); err != nil {
			log.Println("[SelectPageProductType] scan:", err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read", Data: nil})
		}
		if upd.Valid {
			t := upd.Time
			m.UpdateDate = &t
		}
		items = append(items, m)
	}

	var total int
	if err := config.DB.QueryRow(`SELECT COUNT(*) FROM tb_product_type AS p WHERE p.is_delete = 0`).Scan(&total); err != nil {
		log.Println("[SelectPageProductType] count:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to count", Data: nil})
	}

	return c.JSON(models.ApiResponse{
		Status: "success", Message: "",
		Data: fiber.Map{"page": page, "limit": limit, "total": total, "items": items},
	})
}

// ---------- 3) Select By ID (JOIN group) ----------
func SelectProductTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")

	row := config.DB.QueryRow(`
		SELECT  p.product_type_id,
				p.type_name,
				p.product_type_group_id,
				ISNULL(g.group_id, '') AS group_id,
				ISNULL(g.group_name, '') AS group_name,
				p.description,
				p.update_by,
				p.update_date,
				p.id_status,
				p.is_delete
		FROM tb_product_type AS p
		LEFT JOIN tb_product_type_group AS g
			ON g.group_id = p.product_type_group_id
		   AND g.is_delete = 0
		WHERE p.product_type_id = @ID
		  AND p.is_delete = 0`, sql.Named("ID", id))

	var m models.ProductType
	var upd sql.NullTime
	if err := row.Scan(
		&m.ProductTypeID,
		&m.TypeName,
		&m.ProductTypeGroupID,
		&m.GroupId,
		&m.GroupName,
		&m.Description,
		&m.UpdateBy,
		&upd,
		&m.IDStatus,
		&m.IsDelete,
	); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "Not found", Data: nil})
		}
		log.Println("[SelectProductTypeByID] scan:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read", Data: nil})
	}
	if upd.Valid {
		t := upd.Time
		m.UpdateDate = &t
	}

	return c.JSON(models.ApiResponse{Status: "success", Message: "", Data: m})
}

// ---------- 4) Select By Name (LIKE) (JOIN group) ----------
func SelectProductTypeByName(c *fiber.Ctx) error {
	name := c.Params("name")

	rows, err := config.DB.Query(`
		SELECT  p.product_type_id,
				p.type_name,
				p.product_type_group_id,
				ISNULL(g.group_id, '') AS group_id,
				ISNULL(g.group_name, '') AS group_name,
				p.description,
				p.update_by,
				p.update_date,
				p.id_status,
				p.is_delete
		FROM tb_product_type AS p
		LEFT JOIN tb_product_type_group AS g
			ON g.group_id = p.product_type_group_id
		   AND g.is_delete = 0
		WHERE p.type_name LIKE '%' + @Name + '%'
		  AND p.is_delete = 0
		ORDER BY p.type_name ASC, p.product_type_id ASC`, sql.Named("Name", name))
	if err != nil {
		log.Println("[SelectProductTypeByName] query:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch", Data: nil})
	}
	defer rows.Close()

	var out []models.ProductType
	for rows.Next() {
		var m models.ProductType
		var upd sql.NullTime
		if err := rows.Scan(
			&m.ProductTypeID,
			&m.TypeName,
			&m.ProductTypeGroupID,
			&m.GroupId,
			&m.GroupName,
			&m.Description,
			&m.UpdateBy,
			&upd,
			&m.IDStatus,
			&m.IsDelete,
		); err != nil {
			log.Println("[SelectProductTypeByName] scan:", err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read", Data: nil})
		}
		if upd.Valid {
			t := upd.Time
			m.UpdateDate = &t
		}
		out = append(out, m)
	}

	if len(out) == 0 {
		return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "No matching product type found", Data: nil})
	}
	return c.JSON(models.ApiResponse{Status: "success", Message: "", Data: out})
}

// ---------- 5) Insert (รองรับ GroupID string) ----------
func InsertProductType(c *fiber.Ctx) error {
	var m models.ProductType
	if err := c.BodyParser(&m); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid request body"})
	}
	if strings.TrimSpace(m.TypeName) == "" {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "type_name is required"})
	}

	// Determine GroupID (ProductTypeGroupID or TypeNameGroupCode)
	groupID := strings.TrimSpace(m.ProductTypeGroupID)
	if groupID == "" {
		groupID = strings.TrimSpace(m.TypeNameGroupCode)
	}
	if groupID == "" {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Require product_type_group_id"})
	}

	if strings.TrimSpace(m.UpdateBy) == "" {
		m.UpdateBy = resolveUser(c)
	}

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			// Validate Group Exists
			var exists int
			err := tx.QueryRow("SELECT 1 FROM tb_product_type_group WHERE group_id = @ID AND is_delete = 0", sql.Named("ID", groupID)).Scan(&exists)
			if err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("invalid group_id: %s", groupID)
				}
				return err
			}

			// Insert
			_, err = tx.Exec(`
				INSERT INTO tb_product_type (product_type_group_id, type_name, description, update_by)
				VALUES (@GroupID, @TypeName, @Description, @UpdateBy)`,
				sql.Named("GroupID", groupID),
				sql.Named("TypeName", m.TypeName),
				sql.Named("Description", m.Description),
				sql.Named("UpdateBy", m.UpdateBy),
			)
			return err
		},
	})

	if err != nil {
		if strings.Contains(err.Error(), "invalid group_id") {
			return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: err.Error()})
		}
		log.Println("[InsertProductType] exec:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to insert", Data: fiber.Map{"type_name": m.TypeName}})
	}

	return c.Status(201).JSON(models.ApiResponse{Status: "success", Message: "Product type added successfully", Data: fiber.Map{"type_name": m.TypeName}})
}

// ---------- 6) Update By ID (รองรับ GroupID string) ----------
func UpdateProductTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var m models.ProductType
	if err := c.BodyParser(&m); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid request body"})
	}
	if strings.TrimSpace(m.UpdateBy) == "" {
		m.UpdateBy = resolveUser(c)
	}

	// Determine GroupID if provided
	groupID := strings.TrimSpace(m.ProductTypeGroupID)
	if groupID == "" {
		groupID = strings.TrimSpace(m.TypeNameGroupCode)
	}

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			// If GroupID is provided, validate and update it
			if groupID != "" {
				// Validate Group Exists
				var exists int
				err := tx.QueryRow("SELECT 1 FROM tb_product_type_group WHERE group_id = @ID AND is_delete = 0", sql.Named("ID", groupID)).Scan(&exists)
				if err != nil {
					if err == sql.ErrNoRows {
						return fmt.Errorf("invalid group_id: %s", groupID)
					}
					return err
				}

				// Update with GroupID
				res, err := tx.Exec(`
					UPDATE tb_product_type
					SET product_type_group_id = @GroupID,
						type_name             = COALESCE(NULLIF(@TypeName, ''), type_name),
						description           = @Description,
						update_by             = @UpdateBy,
						id_status             = @IDStatus
					WHERE product_type_id     = @ID
					  AND is_delete           = 0`,
					sql.Named("GroupID", groupID),
					sql.Named("TypeName", m.TypeName),
					sql.Named("Description", m.Description),
					sql.Named("UpdateBy", m.UpdateBy),
					sql.Named("IDStatus", m.IDStatus),
					sql.Named("ID", id),
				)
				if err != nil {
					return err
				}
				if aff, _ := res.RowsAffected(); aff == 0 {
					return sql.ErrNoRows
				}
				return nil
			} else {
				// Update without changing GroupID
				res, err := tx.Exec(`
					UPDATE tb_product_type
					SET type_name   = COALESCE(NULLIF(@TypeName,''), type_name),
						description = @Description,
						update_by   = @UpdateBy,
						id_status   = @IDStatus
					WHERE product_type_id = @ID
					  AND is_delete       = 0`,
					sql.Named("TypeName", m.TypeName),
					sql.Named("Description", m.Description),
					sql.Named("UpdateBy", m.UpdateBy),
					sql.Named("IDStatus", m.IDStatus),
					sql.Named("ID", id),
				)
				if err != nil {
					return err
				}
				if aff, _ := res.RowsAffected(); aff == 0 {
					return sql.ErrNoRows
				}
				return nil
			}
		},
	})

	if err != nil {
		if strings.Contains(err.Error(), "invalid group_id") {
			return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: err.Error()})
		}
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "Not found", Data: fiber.Map{"product_type_id": id}})
		}
		log.Println("[UpdateProductTypeByID] exec:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to update", Data: fiber.Map{"product_type_id": id}})
	}

	return c.JSON(models.ApiResponse{Status: "success", Message: "Product type updated successfully", Data: fiber.Map{"product_type_id": id}})
}

// ---------- 7) Delete (Soft) ----------
func DeleteProductTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	username := resolveUser(c)

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(`
				UPDATE tb_product_type
				SET is_delete = 1, update_by = @UpdateBy
				WHERE product_type_id = @ID`,
				sql.Named("ID", id),
				sql.Named("UpdateBy", username),
			)
			if err != nil {
				return err
			}
			if aff, _ := res.RowsAffected(); aff == 0 {
				return sql.ErrNoRows
			}
			return nil
		},
	})
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "Not found", Data: fiber.Map{"product_type_id": id}})
	}
	if err != nil {
		log.Println("[DeleteProductTypeByID] exec:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to soft delete"})
	}
	return c.JSON(models.ApiResponse{Status: "success", Message: "Product type marked as deleted", Data: fiber.Map{"product_type_id": id, "update_by": username}})
}

// ---------- 8) Remove (Hard) ----------
func RemoveProductTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			res, err := tx.Exec(`DELETE FROM tb_product_type WHERE product_type_id = @ID`, sql.Named("ID", id))
			if err != nil {
				return err
			}
			if aff, _ := res.RowsAffected(); aff == 0 {
				return sql.ErrNoRows
			}
			return nil
		},
	})
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "Not found", Data: fiber.Map{"product_type_id": id}})
	}
	if err != nil {
		log.Println("[RemoveProductTypeByID] exec:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to hard delete"})
	}
	return c.JSON(models.ApiResponse{Status: "success", Message: "Product type removed successfully", Data: fiber.Map{"product_type_id": id}})
}
