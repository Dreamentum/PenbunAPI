package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

// Select All
func SelectAllProductTypeGroup(c *fiber.Ctx) error {
	rows, err := config.DB.Query(`
		SELECT group_code, group_name, description, update_by, update_date, id_status, is_delete
		FROM tb_product_type_group
		WHERE is_delete = 0`)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch", Data: nil})
	}
	defer rows.Close()

	var out []models.ProductTypeGroup
	for rows.Next() {
		var g models.ProductTypeGroup
		var upd sql.NullTime
		if err := rows.Scan(&g.GroupCode, &g.GroupName, &g.Description, &g.UpdateBy, &upd, &g.IDStatus, &g.IsDelete); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read", Data: nil})
		}
		if upd.Valid {
			t := upd.Time
			g.UpdateDate = &t
		}
		out = append(out, g)
	}
	return c.JSON(models.ApiResponse{Status: "success", Data: out})
}

// Select Paging
func SelectPageProductTypeGroup(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit
	rows, err := config.DB.Query(`
		SELECT group_code, group_name, description, update_by, update_date, id_status, is_delete
		FROM tb_product_type_group
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY`,
		sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch", Data: nil})
	}
	defer rows.Close()

	var items []models.ProductTypeGroup
	for rows.Next() {
		var g models.ProductTypeGroup
		var upd sql.NullTime
		if err := rows.Scan(&g.GroupCode, &g.GroupName, &g.Description, &g.UpdateBy, &upd, &g.IDStatus, &g.IsDelete); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read", Data: nil})
		}
		if upd.Valid {
			t := upd.Time
			g.UpdateDate = &t
		}
		items = append(items, g)
	}
	var total int
	if err := config.DB.QueryRow(`SELECT COUNT(*) FROM tb_product_type_group WHERE is_delete = 0`).Scan(&total); err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to count", Data: nil})
	}
	return c.JSON(models.ApiResponse{
		Status: "success",
		Data:   fiber.Map{"page": page, "limit": limit, "total": total, "groups": items},
	})
}

// Select By ID
func SelectProductTypeGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	row := config.DB.QueryRow(`
		SELECT group_code, group_name, description, update_by, update_date, id_status, is_delete
		FROM tb_product_type_group
		WHERE group_code = @ID AND is_delete = 0`, sql.Named("ID", id))
	var g models.ProductTypeGroup
	var upd sql.NullTime
	if err := row.Scan(&g.GroupCode, &g.GroupName, &g.Description, &g.UpdateBy, &upd, &g.IDStatus, &g.IsDelete); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "Not found"})
		}
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read"})
	}
	if upd.Valid {
		t := upd.Time
		g.UpdateDate = &t
	}
	return c.JSON(models.ApiResponse{Status: "success", Data: g})
}

// Select By Name (LIKE)
func SelectProductTypeGroupByName(c *fiber.Ctx) error {
	name := c.Params("name")
	rows, err := config.DB.Query(`
		SELECT group_code, group_name, description, update_by, update_date, id_status, is_delete
		FROM tb_product_type_group
		WHERE group_name LIKE '%' + @Name + '%' AND is_delete = 0`, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch", Data: nil})
	}
	defer rows.Close()
	var out []models.ProductTypeGroup
	for rows.Next() {
		var g models.ProductTypeGroup
		var upd sql.NullTime
		if err := rows.Scan(&g.GroupCode, &g.GroupName, &g.Description, &g.UpdateBy, &upd, &g.IDStatus, &g.IsDelete); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read", Data: nil})
		}
		if upd.Valid {
			t := upd.Time
			g.UpdateDate = &t
		}
		out = append(out, g)
	}
	if len(out) == 0 {
		return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "No matching group found"})
	}
	return c.JSON(models.ApiResponse{Status: "success", Data: out})
}

// Insert — Thin: ไม่อ่าน OUTPUT; Trigger/DEFAULT จัดการ
func InsertProductTypeGroup(c *fiber.Ctx) error {
	var g models.ProductTypeGroup
	if err := c.BodyParser(&g); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid request body"})
	}
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(`
				INSERT INTO tb_product_type_group (group_name, description, update_by)
				VALUES (@GroupName, @Description, @UpdateBy)`,
				sql.Named("GroupName", g.GroupName),
				sql.Named("Description", g.Description),
				sql.Named("UpdateBy", g.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to insert"})
	}
	return c.Status(201).JSON(models.ApiResponse{
		Status: "success", Message: "Product type group added successfully",
		Data: fiber.Map{"group_name": g.GroupName},
	})
}

// Update — ให้ Trigger อัปเดตเวลาเอง
func UpdateProductTypeGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var g models.ProductTypeGroup
	if err := c.BodyParser(&g); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid request body"})
	}
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(`
				UPDATE tb_product_type_group SET
					group_name = COALESCE(NULLIF(@GroupName,''), group_name),
					description = @Description,
					update_by = @UpdateBy
				WHERE group_code = @ID AND is_delete = 0`,
				sql.Named("GroupName", g.GroupName),
				sql.Named("Description", g.Description),
				sql.Named("UpdateBy", g.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to update"})
	}
	return c.JSON(models.ApiResponse{Status: "success", Message: "Product type group updated successfully"})
}

// Soft Delete — ให้ Trigger อัปเดตเวลาเอง
func DeleteProductTypeGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Query("user")
	if username == "" {
		username = "UNKNOWN"
	}
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(`
				UPDATE tb_product_type_group
				SET is_delete = 1, update_by = @UpdateBy
				WHERE group_code = @ID`,
				sql.Named("ID", id), sql.Named("UpdateBy", username),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to soft delete"})
	}
	return c.JSON(models.ApiResponse{Status: "success", Message: "Product type group marked as deleted"})
}

// Hard Delete
func RemoveProductTypeGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(`DELETE FROM tb_product_type_group WHERE group_code = @ID`, sql.Named("ID", id))
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to hard delete"})
	}
	return c.JSON(models.ApiResponse{Status: "success", Message: "Product type group removed successfully"})
}
