package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func SelectAllVendor(c *fiber.Ctx) error {
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

func SelectPageVendor(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	query := `
		SELECT vendor_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_vendor_type
		WHERE is_delete = 0
		ORDER BY update_date DESC
		LIMIT ? OFFSET ?
	`
	rows, err := config.DB.Query(query, limit, offset)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch vendor types by page",
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

func SelectVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT vendor_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_vendor_type
		WHERE vendor_type_id = ? AND is_delete = 0
	`
	var vt models.VendorType
	err := config.DB.QueryRow(query, id).Scan(&vt.VendorTypeID, &vt.TypeName, &vt.Description, &vt.UpdateBy, &vt.UpdateDate, &vt.IDStatus)
	if err != nil {
		log.Println(err)
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Vendor type not found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    vt,
	})
}

func SelectVendorByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT vendor_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_vendor_type
		WHERE type_name LIKE CONCAT('%', ?, '%') AND is_delete = 0
	`
	rows, err := config.DB.Query(query, name)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch vendor types by name",
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

func InsertVendor(c *fiber.Ctx) error {
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
		VALUES (?, ?, ?)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, vt.TypeName, vt.Description, vt.UpdateBy)
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

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor type inserted successfully",
		Data:    nil,
	})
}

func UpdateVendorByID(c *fiber.Ctx) error {
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
		SET type_name = COALESCE(NULLIF(?, ''), type_name),
			description = ?,
			update_by = ?
		WHERE vendor_type_id = ? AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, vt.TypeName, vt.Description, vt.UpdateBy, id)
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

func DeleteVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_vendor_type SET is_delete = 1 WHERE vendor_type_id = ?
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, id)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to delete vendor type",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor type soft-deleted successfully",
		Data:    nil,
	})
}

func RemoveVendorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		DELETE FROM tb_vendor_type WHERE vendor_type_id = ?
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, id)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to remove vendor type",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Vendor type removed successfully",
		Data:    nil,
	})
}
