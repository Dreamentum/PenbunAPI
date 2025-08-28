package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

// 🔹 Select All
func SelectAllProductTypeGroup(c *fiber.Ctx) error {
	query := `
		SELECT group_code, group_name, description, update_by, update_date, id_status, is_delete
		FROM tb_product_type_group
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
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
		var updDate sql.NullTime
		if err := rows.Scan(
			&g.GroupCode,
			&g.GroupName,
			&g.Description,
			&g.UpdateBy,
			&updDate,
			&g.IDStatus, // BIT -> bool
			&g.IsDelete,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		if updDate.Valid {
			t := updDate.Time
			g.UpdateDate = &t // *time.Time
		}
		result = append(result, g)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    result,
	})
}

// 🔹 Select Paging
func SelectPageProductTypeGroup(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT group_code, group_name, description, update_by, update_date, id_status, is_delete
		FROM tb_product_type_group
		WHERE is_delete = 0
		ORDER BY update_date DESC
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
		var updDate sql.NullTime
		if err := rows.Scan(
			&g.GroupCode,
			&g.GroupName,
			&g.Description,
			&g.UpdateBy,
			&updDate,
			&g.IDStatus,
			&g.IsDelete,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		if updDate.Valid {
			t := updDate.Time
			g.UpdateDate = &t
		}
		result = append(result, g)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_product_type_group WHERE is_delete = 0`).Scan(&total)
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
			"page":   page,
			"limit":  limit,
			"total":  total,
			"groups": result,
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

	var g models.ProductTypeGroup
	var updDate sql.NullTime
	if err := row.Scan(
		&g.GroupCode,
		&g.GroupName,
		&g.Description,
		&g.UpdateBy,
		&updDate,
		&g.IDStatus,
		&g.IsDelete,
	); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Product type group not found",
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
	if updDate.Valid {
		t := updDate.Time
		g.UpdateDate = &t
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    g,
	})
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
		var updDate sql.NullTime
		if err := rows.Scan(
			&g.GroupCode,
			&g.GroupName,
			&g.Description,
			&g.UpdateBy,
			&updDate,
			&g.IDStatus,
			&g.IsDelete,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		if updDate.Valid {
			t := updDate.Time
			g.UpdateDate = &t
		}
		result = append(result, g)
	}

	if len(result) == 0 {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "No matching product type group found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    result,
	})
}

// 🔹 Insert (ใช้ prefix + Trigger สร้าง group_code และส่งกลับ)
func InsertProductTypeGroup(c *fiber.Ctx) error {
	var g models.ProductTypeGroup
	if err := c.BodyParser(&g); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	// ถ้าไม่ได้ส่ง prefix มา จะใช้ค่าเริ่มต้น
	if g.Prefix == "" {
		g.Prefix = "PRODTG"
	}

	query := `
		INSERT INTO tb_product_type_group (prefix, group_name, description, update_by)
		OUTPUT inserted.group_code
		VALUES (@Prefix, @GroupName, @Description, @UpdateBy)
	`
	var newCode string
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			return tx.QueryRow(query,
				sql.Named("Prefix", g.Prefix),
				sql.Named("GroupName", g.GroupName),
				sql.Named("Description", g.Description),
				sql.Named("UpdateBy", g.UpdateBy),
			).Scan(&newCode)
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to insert product type group",
			Data:    nil,
		})
	}

	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product type group added successfully",
		Data:    fiber.Map{"group_code": newCode},
	})
}

// 🔹 Update (ใช้ COALESCE ตามมาตรฐาน)
func UpdateProductTypeGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var g models.ProductTypeGroup
	if err := c.BodyParser(&g); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		UPDATE tb_product_type_group
		SET group_name = COALESCE(NULLIF(@GroupName, ''), group_name),
			description = @Description,
			update_by  = @UpdateBy
		WHERE group_code = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
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
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update product type group",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product type group updated successfully",
		Data:    nil,
	})
}

// 🔹 Soft Delete (รองรับ ?user=... และ update_date TimeZone SE Asia)
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
		    update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'Se Asia Standard Time' AS DATETIME)
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
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to soft delete product type group",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product type group marked as deleted",
		Data:    nil,
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
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to hard delete product type group",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product type group removed successfully",
		Data:    nil,
	})
}
