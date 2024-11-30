package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"github.com/gofiber/fiber/v2"

	"github.com/sirupsen/logrus" // Logrus สำหรับบันทึก Log
)

func CreateBook(c *fiber.Ctx) error {
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		config.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("[BOOK] Failed to parse book data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	query := "INSERT INTO books (title, author, description) VALUES (?, ?, ?)"
	_, err := config.DB.Exec(query, book.Title, book.Author, book.Description)
	if err != nil {
		config.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"book":  book,
		}).Error("[BOOK] Failed to insert book")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert book"})
	}

	config.Logger.WithFields(logrus.Fields{
		"title": book.Title,
		"author": book.Author,
	}).Info("[BOOK] Book created successfully")
	return c.JSON(fiber.Map{"message": "Book created successfully"})
}

func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")

	var book models.Book
	query := "SELECT id, title, author, description, created_at FROM books WHERE id = ?"
	err := config.DB.QueryRow(query, id).Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.CreatedAt)

	if err != nil {
		if err.Error() == "[BOOK] sql: no rows in result set" {
			config.Logger.WithFields(logrus.Fields{
				"id": id,
			}).Warn("Book not found")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
		}
		config.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("[BOOK] Failed to get book")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get book"})
	}

	config.Logger.WithFields(logrus.Fields{
		"id":    book.ID,
		"title": book.Title,
	}).Info("[BOOK] Book retrieved successfully")
	return c.JSON(book)
}

func GetAllBooks(c *fiber.Ctx) error {
	var books []models.Book
	query := "SELECT id, title, author, description, created_at FROM books"
	rows, err := config.DB.Query(query)
	if err != nil {
		config.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("[BOOK] Failed to get books")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get books"})
	}
	defer rows.Close()

	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.CreatedAt)
		if err != nil {
			config.Logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("[BOOK] Failed to scan book row")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve book data"})
		}
		books = append(books, book)
	}

	config.Logger.Info("[BOOK] All books retrieved successfully")
	return c.JSON(books)
}

func UpdateBook(c *fiber.Ctx) error {
	// รับ ID ของหนังสือจากพารามิเตอร์
	id := c.Params("id")

	// รับข้อมูลหนังสือจาก Request Body
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		config.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"book":  book,
		}).Error("[BOOK] Failed to parse book data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// อัปเดตข้อมูลหนังสือในฐานข้อมูล
	query := "UPDATE books SET title = ?, author = ?, description = ? WHERE id = ?"
	_, err := config.DB.Exec(query, book.Title, book.Author, book.Description, id)
	if err != nil {
		config.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
			"title": book.Title,
			"author": book.Author,
		}).Error("[BOOK] Failed to update book")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update book"})
	}

	// บันทึก Log เมื่อการอัปเดตสำเร็จ
	config.Logger.WithFields(logrus.Fields{
		"id":      id,
		"title":   book.Title,
		"author":  book.Author,
		"message": "Book updated successfully",
	}).Info("Book updated")

	return c.JSON(fiber.Map{"message": "Book updated successfully"})
}

func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	query := "DELETE FROM books WHERE id = ?"
	_, err := config.DB.Exec(query, id)
	if err != nil {
		config.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("[BOOK] Failed to delete book")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete book"})
	}

	config.Logger.WithFields(logrus.Fields{
		"id": id,
	}).Info("Book deleted successfully")
	return c.JSON(fiber.Map{"message": "Book deleted successfully"})
}


