package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

// Helper: Execute Database Transactions
func executeTransaction(db *sql.DB, queries []func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	for _, query := range queries {
		if err := query(tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// Select All Publishers
func SelectAllPublishers(c *fiber.Ctx) error {
	query := `
		SELECT publisher_code, publisher_type_id, publisher_name, contact_name1, contact_name2, 
		       email, phone1, phone2, address, district, province, zip_code, 
			   note, update_by, update_date, id_status
		FROM tb_publisher
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch publishers"})
	}
	defer rows.Close()

	var publishers []models.Publisher

	for rows.Next() {
		var publisher models.Publisher
		if err := rows.Scan(
			&publisher.PublisherCode, &publisher.PublisherTypeID, &publisher.PublisherName,
			&publisher.ContactName1, &publisher.ContactName2,
			&publisher.Email, &publisher.Phone1, &publisher.Phone2,
			&publisher.Address, &publisher.District, &publisher.Province, &publisher.ZipCode,
			&publisher.Note, &publisher.UpdateBy, &publisher.UpdateDate, &publisher.IDStatus,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
		}
		publishers = append(publishers, publisher)
	}

	return c.JSON(publishers)
}

// Select Publisher By ID
func SelectPublisherByID(c *fiber.Ctx) error {
	publisherID := c.Params("id")
	query := `
		SELECT publisher_code, publisher_type_id, publisher_name, contact_name1, contact_name2, 
		       email, phone1, phone2, address, district, province, zip_code, 
			   note, update_by, update_date, id_status
		FROM tb_publisher
		WHERE publisher_code = @PublisherID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("PublisherID", publisherID))

	var publisher models.Publisher
	if err := row.Scan(&publisher.PublisherCode, &publisher.PublisherTypeID, &publisher.PublisherName,
		&publisher.ContactName1, &publisher.ContactName2,
		&publisher.Email, &publisher.Phone1, &publisher.Phone2,
		&publisher.Address, &publisher.District, &publisher.Province, &publisher.ZipCode,
		&publisher.Note, &publisher.UpdateBy, &publisher.UpdateDate, &publisher.IDStatus); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Publisher not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
	}

	return c.JSON(publisher)
}

// Insert Publisher
func InsertPublisher(c *fiber.Ctx) error {
	var publisher models.Publisher
	if err := c.BodyParser(&publisher); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	query := `
		INSERT INTO tb_publisher (publisher_type_id, publisher_name, contact_name1, contact_name2, email, phone1, phone2, address, district, province, zip_code, note, update_by)
		VALUES (@PublisherTypeID, @PublisherName, @ContactName1, @ContactName2, @Email, @Phone1, @Phone2, @Address, @District, @Province, @ZipCode, @Note, @UpdateBy)
	`
	err := executeTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("PublisherTypeID", publisher.PublisherTypeID),
				sql.Named("PublisherName", publisher.PublisherName),
				sql.Named("ContactName1", publisher.ContactName1),
				sql.Named("ContactName2", publisher.ContactName2),
				sql.Named("Email", publisher.Email),
				sql.Named("Phone1", publisher.Phone1),
				sql.Named("Phone2", publisher.Phone2),
				sql.Named("Address", publisher.Address),
				sql.Named("District", publisher.District),
				sql.Named("Province", publisher.Province),
				sql.Named("ZipCode", publisher.ZipCode),
				sql.Named("Note", publisher.Note),
				sql.Named("UpdateBy", publisher.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert publisher"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Publisher added successfully"})
}

// Update Publisher By ID
// UpdatePublisherByID updates a publisher's details by its ID
func UpdatePublisherByID(c *fiber.Ctx) error {
	publisherID := c.Params("id") // Get the publisher ID from the route parameters

	var updatedPublisher models.Publisher
	if err := c.BodyParser(&updatedPublisher); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Begin a transaction
	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("[ERROR] Failed to begin transaction:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database transaction error"})
	}

	// Defer rollback in case of failure
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	query := `
		UPDATE tb_publisher
		SET publisher_type_id = @PublisherTypeID,
		    publisher_name = @PublisherName,
		    contact_name1 = @ContactName1,
		    contact_name2 = @ContactName2,
		    email = @Email,
		    phone1 = @Phone1,
		    phone2 = @Phone2,
		    address = @Address,
		    district = @District,
		    province = @Province,
		    zip_code = @ZipCode,
		    note = @Note,
		    update_by = @UpdateBy
		WHERE publisher_code = @PublisherID AND is_delete = 0
	`

	_, err = tx.Exec(
		query,
		sql.Named("PublisherTypeID", updatedPublisher.PublisherTypeID),
		sql.Named("PublisherName", updatedPublisher.PublisherName),
		sql.Named("ContactName1", updatedPublisher.ContactName1),
		sql.Named("ContactName2", updatedPublisher.ContactName2),
		sql.Named("Email", updatedPublisher.Email),
		sql.Named("Phone1", updatedPublisher.Phone1),
		sql.Named("Phone2", updatedPublisher.Phone2),
		sql.Named("Address", updatedPublisher.Address),
		sql.Named("District", updatedPublisher.District),
		sql.Named("Province", updatedPublisher.Province),
		sql.Named("ZipCode", updatedPublisher.ZipCode),
		sql.Named("Note", updatedPublisher.Note),
		sql.Named("UpdateBy", updatedPublisher.UpdateBy),
		sql.Named("PublisherID", publisherID),
	)

	if err != nil {
		log.Println("[ERROR] Failed to execute update query:", err)
		tx.Rollback()
		return c.JSON(fiber.Map{"error": "Failed to update publisher"})
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Println("[ERROR] Failed to commit transaction:", err)
		return c.JSON(fiber.Map{"error": "Transaction commit failed"})
	}

	return c.JSON(fiber.Map{
		"message": "Publisher updated successfully",
		"data": map[string]interface{}{
			"publisher_id":     publisherID,
			"publisher_name":   updatedPublisher.PublisherName,
			"updated_by":       updatedPublisher.UpdateBy,
			"update_timestamp": updatedPublisher.UpdateDate,
		},
	})
}

// Delete Publisher (Update is_delete)
func DeletePublisherByID(c *fiber.Ctx) error {
	publisherID := c.Params("id")
	query := `
		UPDATE tb_publisher
		SET is_delete = 1, update_date = GETDATE()
		WHERE publisher_code = @PublisherID
	`
	_, err := config.DB.Exec(query, sql.Named("PublisherID", publisherID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete publisher"})
	}

	return c.JSON(fiber.Map{"message": "Publisher marked as deleted"})
}

// Remove Publisher (Hard Delete)
func RemovePublisherByID(c *fiber.Ctx) error {
	publisherID := c.Params("id")
	query := `
		DELETE FROM tb_publisher WHERE publisher_code = @PublisherID
	`
	_, err := config.DB.Exec(query, sql.Named("PublisherID", publisherID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove publisher"})
	}

	return c.JSON(fiber.Map{"message": "Publisher removed successfully"})
}
