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
				ISNULL(g.group_code, '') AS group_code,
				ISNULL(g.group_name, '') AS group_name,
				p.description,
				p.update_by,
				p.update_date,
				p.id_status,
				p.is_delete
		FROM tb_product_type AS p
		LEFT JOIN tb_product_type_group AS g
			ON g.autoID = p.product_type_group_id
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
			&m.GroupCode,
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
				ISNULL(g.group_code, '') AS group_code,
				ISNULL(g.group_name, '') AS group_name,
				p.description,
				p.update_by,
				p.update_date,
				p.id_status,
				p.is_delete
		FROM tb_product_type AS p
		LEFT JOIN tb_product_type_group AS g
			ON g.autoID = p.product_type_group_id
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

	var items []models.ProductType
	for rows.Next() {
		var m models.ProductType
		var upd sql.NullTime
		if err := rows.Scan(
			&m.ProductTypeID,
			&m.TypeName,
			&m.ProductTypeGroupID,
			&m.GroupCode,
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
				ISNULL(g.group_code, '') AS group_code,
				ISNULL(g.group_name, '') AS group_name,
				p.description,
				p.update_by,
				p.update_date,
				p.id_status,
				p.is_delete
		FROM tb_product_type AS p
		LEFT JOIN tb_product_type_group AS g
			ON g.autoID = p.product_type_group_id
		   AND g.is_delete = 0
		WHERE p.product_type_id = @ID
		  AND p.is_delete = 0`, sql.Named("ID", id))

	var m models.ProductType
	var upd sql.NullTime
	if err := row.Scan(
		&m.ProductTypeID,
		&m.TypeName,
		&m.ProductTypeGroupID,
		&m.GroupCode,
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
				ISNULL(g.group_code, '') AS group_code,
				ISNULL(g.group_name, '') AS group_name,
				p.description,
				p.update_by,
				p.update_date,
				p.id_status,
				p.is_delete
		FROM tb_product_type AS p
		LEFT JOIN tb_product_type_group AS g
			ON g.autoID = p.product_type_group_id
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
			&m.GroupCode,
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

// ---------- 5) Insert (รองรับ ID หรือ CODE) ----------
func InsertProductType(c *fiber.Ctx) error {
	var m models.ProductType
	if err := c.BodyParser(&m); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid request body"})
	}
	if strings.TrimSpace(m.TypeName) == "" {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "type_name is required"})
	}
	if m.ProductTypeGroupID == 0 && strings.TrimSpace(m.TypeNameGroupCode) == "" {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Require product_type_group_id or type_name_group_code"})
	}
	if strings.TrimSpace(m.UpdateBy) == "" {
		m.UpdateBy = resolveUser(c)
	}

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			// case 1: ใช้ ID ตรง ๆ
			if m.ProductTypeGroupID > 0 {
				_, err := tx.Exec(`
					INSERT INTO tb_product_type (product_type_group_id, type_name, description, update_by)
					VALUES (@GroupID, @TypeName, @Description, @UpdateBy)`,
					sql.Named("GroupID", m.ProductTypeGroupID),
					sql.Named("TypeName", m.TypeName),
					sql.Named("Description", m.Description),
					sql.Named("UpdateBy", m.UpdateBy),
				)
				return err
			}

			// case 2: ใช้ CODE → map เป็น autoID (ถ้าไม่พบจะไม่ insert)
			res, err := tx.Exec(`
				INSERT INTO tb_product_type (product_type_group_id, type_name, description, update_by)
				SELECT g.autoID, @TypeName, @Description, @UpdateBy
				FROM tb_product_type_group AS g
				WHERE g.group_code = @GroupCode
				  AND g.is_delete = 0`,
				sql.Named("TypeName", m.TypeName),
				sql.Named("Description", m.Description),
				sql.Named("UpdateBy", m.UpdateBy),
				sql.Named("GroupCode", m.TypeNameGroupCode),
			)
			if err != nil {
				return err
			}
			if aff, _ := res.RowsAffected(); aff == 0 {
				return sql.ErrNoRows // ถือว่า group_code ไม่ถูกต้อง
			}
			return nil
		},
	})
	if err == sql.ErrNoRows {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid group_code", Data: fiber.Map{"group_code": m.TypeNameGroupCode}})
	}
	if err != nil {
		log.Println("[InsertProductType] exec:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to insert", Data: fiber.Map{"type_name": m.TypeName}})
	}

	return c.Status(201).JSON(models.ApiResponse{Status: "success", Message: "Product type added successfully", Data: fiber.Map{"type_name": m.TypeName}})
}

// ---------- 6) Update By ID (รองรับ ID หรือ CODE) ----------
func UpdateProductTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var m models.ProductType
	if err := c.BodyParser(&m); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid request body"})
	}
	if strings.TrimSpace(m.UpdateBy) == "" {
		m.UpdateBy = resolveUser(c)
	}

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			// 6.1) เปลี่ยนกลุ่มด้วย ID โดยตรง (ถ้า > 0)
			if m.ProductTypeGroupID > 0 {
				res, err := tx.Exec(`
					UPDATE tb_product_type
					SET product_type_group_id = @GroupID,
						type_name             = COALESCE(NULLIF(@TypeName, ''), type_name),
						description           = @Description,
						update_by             = @UpdateBy
					WHERE product_type_id     = @ID
					  AND is_delete           = 0`,
					sql.Named("GroupID", m.ProductTypeGroupID),
					sql.Named("TypeName", m.TypeName),
					sql.Named("Description", m.Description),
					sql.Named("UpdateBy", m.UpdateBy),
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

			// 6.2) เปลี่ยนกลุ่มด้วย CODE (ถ้าส่งมา)
			if strings.TrimSpace(m.TypeNameGroupCode) != "" {
				res, err := tx.Exec(`
					UPDATE p
					SET p.product_type_group_id = g.autoID,
						p.type_name             = COALESCE(NULLIF(@TypeName,''), p.type_name),
						p.description           = @Description,
						p.update_by             = @UpdateBy
					FROM tb_product_type AS p
					INNER JOIN tb_product_type_group AS g
							ON g.group_code = @GroupCode
						   AND g.is_delete = 0
					WHERE p.product_type_id = @ID
					  AND p.is_delete       = 0`,
					sql.Named("TypeName", m.TypeName),
					sql.Named("Description", m.Description),
					sql.Named("UpdateBy", m.UpdateBy),
					sql.Named("GroupCode", m.TypeNameGroupCode),
					sql.Named("ID", id),
				)
				if err != nil {
					return err
				}
				if aff, _ := res.RowsAffected(); aff == 0 {
					return sql.ErrNoRows // ไม่พบ p หรือ group_code ไม่ถูกต้อง
				}
				return nil
			}

			// 6.3) ไม่เปลี่ยนกลุ่ม → อัปเดตเฉพาะข้อมูลอื่น
			res, err := tx.Exec(`
				UPDATE tb_product_type
				SET type_name   = COALESCE(NULLIF(@TypeName,''), type_name),
					description = @Description,
					update_by   = @UpdateBy
				WHERE product_type_id = @ID
				  AND is_delete       = 0`,
				sql.Named("TypeName", m.TypeName),
				sql.Named("Description", m.Description),
				sql.Named("UpdateBy", m.UpdateBy),
				sql.Named("ID", id),
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
		return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "Not found or invalid group", Data: fiber.Map{"product_type_id": id}})
	}
	if err != nil {
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
