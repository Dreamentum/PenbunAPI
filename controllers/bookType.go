package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func SelectAllBookTypes(c *fiber.Ctx) error {
	query := `SELECT book_type_id, type_name, description, update_by, update_date, id_status 
			  FROM tb_book_type WHERE is_delete = 0`
	rows, err := config.DB.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch book types"})
	}
	defer rows.Close()

	var bookTypes []models.BookType
	for rows.Next() {
		var b models.BookType
		if err := rows.Scan(&b.BookTypeID, &b.TypeName, &b.Description, &b.UpdateBy, &b.UpdateDate, &b.IDStatus); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to read data"})
		}
		bookTypes = append(bookTypes, b)
	}
	return c.JSON(bookTypes)
}

func SelectPageBookTypes(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `SELECT book_type_id, type_name, description, update_by, update_date 
			  FROM tb_book_type WHERE is_delete = 0 
			  ORDER BY update_date DESC OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch book types"})
	}
	defer rows.Close()

	var bookTypes []models.BookType
	for rows.Next() {
		var b models.BookType
		if err := rows.Scan(&b.BookTypeID, &b.TypeName, &b.Description, &b.UpdateBy, &b.UpdateDate); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to read data"})
		}
		bookTypes = append(bookTypes, b)
	}
	return c.JSON(fiber.Map{"page": page, "limit": limit, "data": bookTypes})
}

func SelectBookTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `SELECT book_type_id, type_name, description, update_by, update_date, id_status 
			  FROM tb_book_type WHERE book_type_id = @ID AND is_delete = 0`
	row := config.DB.QueryRow(query, sql.Named("ID", id))
	var b models.BookType
	if err := row.Scan(&b.BookTypeID, &b.TypeName, &b.Description, &b.UpdateBy, &b.UpdateDate, &b.IDStatus); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Book type not found"})
	}
	return c.JSON(b)
}

func InsertBookType(c *fiber.Ctx) error {
	var b models.BookType
	if err := c.BodyParser(&b); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	query := `INSERT INTO tb_book_type (type_name, description, update_by) VALUES (@TypeName, @Description, @UpdateBy)`
	_, err := config.DB.Exec(query, sql.Named("TypeName", b.TypeName), sql.Named("Description", b.Description), sql.Named("UpdateBy", b.UpdateBy))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to insert book type"})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Book type inserted successfully"})
}

func UpdateBookTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var b models.BookType
	if err := c.BodyParser(&b); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	query := `UPDATE tb_book_type SET type_name = @TypeName, description = @Description, update_by = @UpdateBy 
			  WHERE book_type_id = @ID AND is_delete = 0`
	_, err := config.DB.Exec(query, sql.Named("TypeName", b.TypeName), sql.Named("Description", b.Description), sql.Named("UpdateBy", b.UpdateBy), sql.Named("ID", id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update book type"})
	}
	return c.JSON(fiber.Map{"message": "Book type updated successfully"})
}

func DeleteBookTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `UPDATE tb_book_type SET is_delete = 1, update_date = GETDATE() WHERE book_type_id = @ID`
	_, err := config.DB.Exec(query, sql.Named("ID", id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete book type"})
	}
	return c.JSON(fiber.Map{"message": "Book type marked as deleted"})
}

func RemoveBookTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_book_type WHERE book_type_id = @ID`
	_, err := config.DB.Exec(query, sql.Named("ID", id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to remove book type"})
	}
	return c.JSON(fiber.Map{"message": "Book type removed successfully"})
}
