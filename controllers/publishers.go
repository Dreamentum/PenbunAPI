package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// Select All Publishers
func SelectAllPublishers(c *fiber.Ctx) error {
	query := `SELECT publisher_ID, publisher_name, note, update_date FROM tb_publisher WHERE is_delete = 0`
	rows, err := config.DB.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch publishers"})
	}
	defer rows.Close()

	var publishers []models.Publisher
	for rows.Next() {
		var publisher models.Publisher
		if err := rows.Scan(&publisher.PublisherID, &publisher.PublisherName, &publisher.Note, &publisher.UpdateDate); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
		}
		publishers = append(publishers, publisher)
	}

	return c.JSON(publishers)
}

// Select Publisher By ID
func SelectPublisherByID(c *fiber.Ctx) error {
	publisherID := c.Params("id")
	query := `SELECT publisher_ID, publisher_name, note, update_date FROM tb_publisher WHERE is_delete = 0 AND publisher_ID = @PublisherID`
	row := config.DB.QueryRow(query, sql.Named("PublisherID", publisherID))

	var publisher models.Publisher
	if err := row.Scan(&publisher.PublisherID, &publisher.PublisherName, &publisher.Note, &publisher.UpdateDate); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Publisher not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
	}

	return c.JSON(publisher)
}

func InsertPublisher(c *fiber.Ctx) error {
	var publisher models.Publisher
	if err := c.BodyParser(&publisher); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	query := `INSERT INTO tb_publisher (publisher_name, note) VALUES (@PublisherName, @Note)`
	_, err := config.DB.Exec(query,
		sql.Named("PublisherName", publisher.PublisherName),
		sql.Named("Note", publisher.Note),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert publisher"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Publisher added successfully"})
}

// Update Publisher By ID
func UpdatePublisherByID(c *fiber.Ctx) error {
	var publisher models.Publisher
	if err := c.BodyParser(&publisher); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	publisherID := c.Params("id")
	query := `UPDATE tb_publisher SET publisher_name = @PublisherName, note = @Note, update_date = GETDATE() WHERE publisher_ID = @PublisherID`
	_, err := config.DB.Exec(query,
		sql.Named("PublisherName", publisher.PublisherName),
		sql.Named("Note", publisher.Note),
		sql.Named("PublisherID", publisherID),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update publisher"})
	}

	return c.JSON(fiber.Map{"message": "Publisher updated successfully"})
}

// Delete Publisher (Update is_delete)
func DeletePublisher(c *fiber.Ctx) error {
	publisherID := c.Params("id")
	query := `UPDATE tb_publisher SET is_delete = 1, update_date = GETDATE() WHERE publisher_ID = @PublisherID`
	_, err := config.DB.Exec(query, sql.Named("PublisherID", publisherID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete publisher"})
	}

	return c.JSON(fiber.Map{"message": "Publisher marked as deleted"})
}

// Remove Publisher (Hard Delete)
func RemovePublisher(c *fiber.Ctx) error {
	publisherID := c.Params("id")
	query := `DELETE FROM tb_publisher WHERE publisher_ID = @PublisherID`
	_, err := config.DB.Exec(query, sql.Named("PublisherID", publisherID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove publisher"})
	}

	return c.JSON(fiber.Map{"message": "Publisher removed successfully"})
}
