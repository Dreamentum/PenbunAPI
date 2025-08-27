package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

// 🔹 Select All (Business Response + แปลงเวลาเป็น string)
func SelectAllProductTypeGroup(c *fiber.Ctx) error {
	query := `
		SELECT
			group_code,
			group_name,
			description,
			update_by,
			CONVERT(varchar(19), update_date, 120) AS update_date, -- "YYYY-MM-DD HH:MM:SS"
			id_status,
			is_delete
		FROM tb_product_type_group
		WHERE is_delete = 0
		ORDER BY group_name ASC;
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println("[PTG][select/all] query error:", err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch product type groups",
			Data:    nil,
		})
	}
	defer rows.Close()

	result := make([]models.ProductTypeGroup, 0, 64)
	for rows.Next() {
		var g models.ProductTypeGroup
		if err := rows.Scan(
			&g.GroupCode,
			&g.GroupName,
			&g.Description, // *string
			&g.UpdateBy,    // *string
			&g.UpdateDate,  // *string ← safe เพราะเรา CONVERT มาแล้ว
			&g.IDStatus,    // string
			&g.IsDelete,    // bool
		); err != nil {
			log.Println("[PTG][select/all] scan error:", err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		result = append(result, g) // ✅ อย่าลืม append
	}

	if err := rows.Err(); err != nil {
		log.Println("[PTG][select/all] rows err:", err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Read rows failed",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "OK",
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
			Message: "Failed to fetch paged data",
			Data:    nil})
	}
	defer rows.Close()

	var result []models.ProductTypeGroup
	for rows.Next() {
		var g models.ProductTypeGroup
		if err := rows.Scan(&g.GroupCode, &g.GroupName, &g.Description, &g.UpdateBy, &g.UpdateDate, &g.IDStatus, &g.IsDelete); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil})
		}
		result = append(result, g)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_product_type_group WHERE is_delete = 0`).Scan(&total)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to count data",
			Data:    nil})
	}

	return c.JSON(models.ApiResponse{
		Status: "success",
		Data: fiber.Map{
			"page":   page,
			"limit":  limit,
			"total":  total,
			"groups": result}})
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
	if err := row.Scan(&g.GroupCode, &g.GroupName, &g.Description, &g.UpdateBy, &g.UpdateDate, &g.IDStatus, &g.IsDelete); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Group not found",
				Data:    nil})
		}
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to read data",
			Data:    nil})
	}
	return c.JSON(models.ApiResponse{
		Status: "success",
		Data:   g})
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
			Message: "Failed to search group",
			Data:    nil})
	}
	defer rows.Close()

	var result []models.ProductTypeGroup
	for rows.Next() {
		var g models.ProductTypeGroup
		if err := rows.Scan(&g.GroupCode, &g.GroupName, &g.Description, &g.UpdateBy, &g.UpdateDate, &g.IDStatus, &g.IsDelete); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil})
		}
		result = append(result, g)
	}

	if len(result) == 0 {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "No matching group found",
			Data:    nil})
	}
	return c.JSON(models.ApiResponse{
		Status: "success",
		Data:   result})
}

// 🔹 Insert
func InsertProductTypeGroup(c *fiber.Ctx) error {
	var g models.ProductTypeGroup
	if err := c.BodyParser(&g); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid body",
			Data:    nil})
	}

	query := `
		INSERT INTO tb_product_type_group (group_code, group_name, description, update_by)
		VALUES (@GroupCode, @GroupName, @Description, @UpdateBy)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("GroupCode", g.GroupCode),
				sql.Named("GroupName", g.GroupName),
				sql.Named("Description", g.Description),
				sql.Named("UpdateBy", g.UpdateBy))
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Insert failed",
			Data:    nil})
	}
	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Group inserted",
		Data:    nil})
}

// 🔹 Update
func UpdateProductTypeGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var g models.ProductTypeGroup
	if err := c.BodyParser(&g); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid body",
			Data:    nil})
	}

	query := `
		UPDATE tb_product_type_group
		SET group_name = COALESCE(NULLIF(@GroupName, ''), group_name),
			description = @Description,
			update_by = @UpdateBy
		WHERE group_code = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("GroupName", g.GroupName),
				sql.Named("Description", g.Description),
				sql.Named("UpdateBy", g.UpdateBy),
				sql.Named("ID", id))
			return err
		}})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Update failed", Data: nil})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product Type Group updated"})
}

// 🔹 Soft Delete
func DeleteProductTypeGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Query("user") // 👈 รับค่าจาก query param ?user=ROOT or ADMIN

	if username == "" {
		username = "UNKNOWN"
	}
	query := `
		UPDATE tb_product_type_group
		SET is_delete = 1, update_by = @UpdateBy
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
		log.Printf("[ERROR] DeleteProductTypeGroupByID: %v\n", err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to soft delete product type",
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
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Hard delete failed", Data: nil})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Group hard deleted"})
}
