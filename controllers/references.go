package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

// GetReference retrieves a specific reference by parameter
func GetReference(c *fiber.Ctx) error {
	// รับค่าพารามิเตอร์ เช่น ref_id
	parameter := c.Query("parameter")  // ชื่อฟิลด์ใน WHERE
	value := c.Query("value")          // ค่าที่จะค้นหา

	if parameter == "" || value == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing parameter or value"})
	}

	// สร้างคำสั่ง SQL โดยใช้ parameter และ value
	query := "SELECT row_id, ref_id, ref_int, ref_text, update_by, update_date FROM tb_reference WHERE " + parameter + " = @Value"

	var reference models.Reference

	err := config.DB.QueryRow(query, sql.Named("Value", value)).Scan(
		&reference.RowID,
		&reference.RefID,
		&reference.RefInt,
		&reference.RefText,
		&reference.UpdateBy,
		&reference.UpdateDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Reference not found"})
		}
		log.Println("[ERROR] Query failed:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database query failed"})
	}

	// ส่งข้อมูลกลับในรูปแบบ JSON
	return c.JSON(reference)
}
