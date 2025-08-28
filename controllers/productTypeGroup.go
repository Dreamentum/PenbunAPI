package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// --- helpers ----------------------------------------------------

func strPtr(ns sql.NullString) *string {
	if ns.Valid {
		s := ns.String
		return &s
	}
	return nil
}

func timePtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		t := nt.Time
		return &t
	}
	return nil
}

func scanGroup(row *sql.Row) (models.ProductTypeGroup, error) {
	var g models.ProductTypeGroup
	var desc, updBy sql.NullString
	var updDate sql.NullTime
	// id_status เป็น NVARCHAR → string
	if err := row.Scan(
		&g.GroupCode,
		&g.GroupName,
		&desc,
		&updBy,
		&updDate,
		&g.IDStatus,
		&g.IsDelete,
	); err != nil {
		return g, err
	}
	g.Description = strPtr(desc)
	g.UpdateBy = strPtr(updBy)
	g.UpdateDate = timePtr(updDate)
	return g, nil
}

func scanGroups(rows *sql.Rows) ([]models.ProductTypeGroup, error) {
	var list []models.ProductTypeGroup
	for rows.Next() {
		var g models.ProductTypeGroup
		var desc, updBy sql.NullString
		var updDate sql.NullTime
		if err := rows.Scan(
			&g.GroupCode,
			&g.GroupName,
			&desc,
			&updBy,
			&updDate,
			&g.IDStatus,
			&g.IsDelete,
		); err != nil {
			return nil, err
		}
		g.Description = strPtr(desc)
		g.UpdateBy = strPtr(updBy)
		g.UpdateDate = timePtr(updDate)
		list = append(list, g)
	}
	return list, nil
}

// 🔹 Select All
func SelectAllProductTypeGroup(c *fiber.Ctx) error {
	query := `
		SELECT group_code, group_name, description, update_by, update_date, id_status, is_delete
		FROM tb_product_type_group
		WHERE is_delete = 0
		ORDER BY group_name ASC
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println("[SelectAllProductTypeGroup] query error:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch product type groups", Data: nil})
	}
	defer rows.Close()

	result, err := scanGroups(rows)
	if err != nil {
		log.Println("[SelectAllProductTypeGroup] scan error:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data", Data: nil})
	}
	return c.JSON(models.ApiResponse{Status: "success", Data: result})
}

// 🔹 Select Paging
// 🔹 Select Paging (สไตล์เดียวกับ SelectPageVendor)
func SelectPageProductTypeGroup(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT g.group_code, g.group_name, g.description,
		       g.update_by, g.update_date, g.id_status, g.is_delete
		FROM tb_product_type_group g
		WHERE g.is_delete = 0
		ORDER BY g.update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch product type groups",
			Data:    nil,
		})
	}
	defer rows.Close()

	var result []models.ProductTypeGroup
	for rows.Next() {
		var g models.ProductTypeGroup
		if err := rows.Scan(
			&g.GroupCode, &g.GroupName, &g.Description,
			&g.UpdateBy, &g.UpdateDate, &g.IDStatus, &g.IsDelete,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		result = append(result, g)
	}

	var total int
	if err := config.DB.QueryRow(`SELECT COUNT(*) FROM tb_product_type_group WHERE is_delete = 0`).Scan(&total); err != nil {
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
			"page":   page,
			"limit":  limit,
			"total":  total,
			"groups": result, // ใช้ key 'groups' ให้ตรงกับของเดิมใน controller นี้
		},
	})
}

// 🔹 Select By ID
func SelectProductTypeGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT group_code, group_name, description, update_by, update_date, id_status, is_delete
		FROM tb_product_type_group
		WHERE group_code = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))
	g, err := scanGroup(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "Group not found", Data: nil})
		}
		log.Println("[SelectProductTypeGroupByID] scan error:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data", Data: nil})
	}
	return c.JSON(models.ApiResponse{Status: "success", Data: g})
}

// 🔹 Select By Name
func SelectProductTypeGroupByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT group_code, group_name, description, update_by, update_date, id_status, is_delete
		FROM tb_product_type_group
		WHERE group_name LIKE '%' + @Name + '%' AND is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println("[SelectProductTypeGroupByName] query error:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to search group", Data: nil})
	}
	defer rows.Close()

	result, err := scanGroups(rows)
	if err != nil {
		log.Println("[SelectProductTypeGroupByName] scan error:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read data", Data: nil})
	}

	if len(result) == 0 {
		return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "No matching group found", Data: nil})
	}
	return c.JSON(models.ApiResponse{Status: "success", Data: result})
}

// 🔹 Insert (ใช้ prefix + Trigger สร้าง group_code)
func InsertProductTypeGroup(c *fiber.Ctx) error {
	var g models.ProductTypeGroup
	if err := c.BodyParser(&g); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid body", Data: nil})
	}
	if g.GroupName == "" {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "group_name is required", Data: nil})
	}
	// ถ้าไม่ได้ส่ง prefix มา กำหนดค่าปริยายสำหรับกลุ่มประเภทสินค้า
	if g.Prefix == "" {
		g.Prefix = "PRODTG"
	}

	// ห้ามส่ง group_code มาเอง — ให้ Trigger สร้าง
	query := `
		INSERT INTO tb_product_type_group (prefix, group_name, description, update_by)
		OUTPUT inserted.group_code
		VALUES (@Prefix, @GroupName, NULLIF(@Description, ''), @UpdateBy)
	`
	var newCode string
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			// ใช้ sql.NullString เพื่อรองรับ nil
			var desc sql.NullString
			if g.Description != nil && *g.Description != "" {
				desc = sql.NullString{String: *g.Description, Valid: true}
			}
			var updBy sql.NullString
			if g.UpdateBy != nil && *g.UpdateBy != "" {
				updBy = sql.NullString{String: *g.UpdateBy, Valid: true}
			}
			return tx.QueryRow(query,
				sql.Named("Prefix", g.Prefix),
				sql.Named("GroupName", g.GroupName),
				sql.Named("Description", desc),
				sql.Named("UpdateBy", updBy),
			).Scan(&newCode)
		},
	})
	if err != nil {
		log.Println("[InsertProductTypeGroup] insert error:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Insert failed", Data: nil})
	}
	return c.Status(201).JSON(models.ApiResponse{Status: "success", Message: "Group inserted", Data: fiber.Map{"group_code": newCode}})
}

// 🔹 Update
func UpdateProductTypeGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var g models.ProductTypeGroup
	if err := c.BodyParser(&g); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid body", Data: nil})
	}

	query := `
		UPDATE tb_product_type_group
		SET group_name = COALESCE(NULLIF(@GroupName, ''), group_name),
			description = COALESCE(NULLIF(@Description, ''), description),
			update_by  = @UpdateBy
		WHERE group_code = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			var desc sql.NullString
			if g.Description != nil && *g.Description != "" {
				desc = sql.NullString{String: *g.Description, Valid: true}
			}
			var updBy sql.NullString
			if g.UpdateBy != nil && *g.UpdateBy != "" {
				updBy = sql.NullString{String: *g.UpdateBy, Valid: true}
			}
			_, err := tx.Exec(query,
				sql.Named("GroupName", g.GroupName),
				sql.Named("Description", desc),
				sql.Named("UpdateBy", updBy),
				sql.Named("ID", id),
			)
			return err
		}})
	if err != nil {
		log.Println("[UpdateProductTypeGroupByID] update error:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Update failed", Data: nil})
	}
	return c.JSON(models.ApiResponse{Status: "success", Message: "Group updated"})
}

// 🔹 Soft Delete
func DeleteProductTypeGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Query("user")
	if username == "" {
		username = "UNKNOWN"
	}

	query := `
		UPDATE tb_product_type_group
		SET is_delete = 1,
		    update_by = @UpdateBy,
		    update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE group_code = @ID
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("ID", id),
				sql.Named("UpdateBy", username),
			)
			return err
		},
	})
	if err != nil {
		log.Printf("[DeleteProductTypeGroupByID] soft delete error: %v\n", err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to soft delete product type group",
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product type group deleted (soft)",
	})
}

// 🔹 Hard Delete
func RemoveProductTypeGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_product_type_group WHERE group_code = @ID`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, sql.Named("ID", id))
			return err
		}})
	if err != nil {
		log.Println("[RemoveProductTypeGroupByID] hard delete error:", err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Hard delete failed", Data: nil})
	}
	return c.JSON(models.ApiResponse{Status: "success", Message: "Group hard deleted"})
}
