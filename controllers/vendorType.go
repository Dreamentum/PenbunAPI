package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllVendorType(c *fiber.Ctx) error {
	query := `
		SELECT vendor_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_vendor_type
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
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
		if err := rows.Scan(&vt.VendorTypeID, &vt.TypeName, &vt.Description, &vt.UpdateBy, &vt.UpdateDate, &vt.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		result = append(result, vt)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    result,
	})
}

func SelectPageVendorType(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT vendor_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_vendor_type
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
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
		if err := rows.Scan(&vt.VendorTypeID, &vt.TypeName, &vt.Description, &vt.UpdateBy, &vt.UpdateDate, &vt.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		result = append(result, vt)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_vendor_type WHERE is_delete = 0`).Scan(&total)
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
			"page":       page,
			"limit":      limit,
			"total":      total,
			"vendorType": result,
		},
	})
}

func SelectVendorTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT vendor_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_vendor_type
		WHERE vendor_type_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var vt models.VendorType
	if err := row.Scan(&vt.VendorTypeID, &vt.TypeName, &vt.Description, &vt.UpdateBy, &vt.UpdateDate, &vt.IDStatus); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Vendor type not found",
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

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    vt,
	})
}

func SelectVendorTypeByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT vendor_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_vendor_type
		WHERE type_name LIKE '%' + @Name + '%' AND is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
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
		if err := rows.Scan(&vt.VendorTypeID, &vt.TypeName, &vt.Description, &vt.UpdateBy, &vt.UpdateDate, &vt.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
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

func InsertVendorType(c *fiber.Ctx) error {
	var vt models.VendorType
	if err := c.BodyParser(&vt); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		INSERT INTO tb_vendor_type (type_name, description, update_by)
		VALUES (@TypeName, @Description, @UpdateBy)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeName", vt.TypeName),
				sql.Named("Description", vt.Description),
				sql.Named("UpdateBy", vt.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to insert vendor type",
			Data:    nil,
		})
	}

	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor type added successfully",
		Data:    nil,
	})
}

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

	query := `
		UPDATE tb_vendor_type
		SET type_name = COALESCE(NULLIF(@TypeName, ''), type_name),
			description = @Description,
			update_by = @UpdateBy
		WHERE vendor_type_id = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeName", vt.TypeName),
				sql.Named("Description", vt.Description),
				sql.Named("UpdateBy", vt.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update vendor type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor type updated successfully",
		Data:    nil,
	})
}

func DeleteVendorTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_vendor_type
		SET is_delete = 1,
			update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE vendor_type_id = @ID
	`
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
			Message: "Failed to soft delete vendor type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor type marked as deleted",
		Data:    nil,
	})
}

func RemoveVendorTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_vendor_type WHERE vendor_type_id = @ID`
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
			Message: "Failed to hard delete vendor type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor type removed successfully",
		Data:    nil,
	})
}
