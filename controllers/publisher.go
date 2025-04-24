package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

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

func SelectAllPublishers(c *fiber.Ctx) error {
	query := `
		SELECT publisher_code, publisher_type_id, publisher_name, contact_name1, contact_name2,
		       email, phone1, phone2, address, district, province, zip_code,
			   note, discount_id, update_by, update_date, id_status
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
			&publisher.Note, &publisher.DiscountID, &publisher.UpdateBy, &publisher.UpdateDate, &publisher.IDStatus,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
		}
		publishers = append(publishers, publisher)
	}
	return c.JSON(publishers)
}

func SelectPagePublishers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT publisher_code, publisher_type_id, publisher_name, contact_name1, contact_name2,
		       email, phone1, phone2, address, district, province, zip_code,
			   note, discount_id, update_by, update_date, id_status
		FROM tb_publisher
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
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
			&publisher.Note, &publisher.DiscountID, &publisher.UpdateBy, &publisher.UpdateDate, &publisher.IDStatus,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
		}
		publishers = append(publishers, publisher)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_publisher WHERE is_delete = 0`).Scan(&total)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count records"})
	}

	return c.JSON(fiber.Map{
		"page":       page,
		"limit":      limit,
		"total":      total,
		"publishers": publishers,
	})
}

func SelectPublisherByID(c *fiber.Ctx) error {
	publisherID := c.Params("id")
	query := `
		SELECT publisher_code, publisher_type_id, publisher_name, discount_id, contact_name1, contact_name2,
		       email, phone1, phone2, address, district, province, zip_code,
			   note, update_by, update_date, id_status
		FROM tb_publisher
		WHERE publisher_code = @PublisherID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("PublisherID", publisherID))
	var publisher models.Publisher
	if err := row.Scan(
		&publisher.PublisherCode, &publisher.PublisherTypeID, &publisher.PublisherName, &publisher.DiscountID,
		&publisher.ContactName1, &publisher.ContactName2,
		&publisher.Email, &publisher.Phone1, &publisher.Phone2,
		&publisher.Address, &publisher.District, &publisher.Province, &publisher.ZipCode,
		&publisher.Note, &publisher.UpdateBy, &publisher.UpdateDate, &publisher.IDStatus,
	); err != nil {
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
	query := `
		INSERT INTO tb_publisher (publisher_type_id, publisher_name, discount_id, contact_name1, contact_name2,
			email, phone1, phone2, address, district, province, zip_code, note, update_by)
		VALUES (@PublisherTypeID, @PublisherName, @DiscountID, @ContactName1, @ContactName2,
			@Email, @Phone1, @Phone2, @Address, @District, @Province, @ZipCode, @Note, @UpdateBy)
	`
	err := executeTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("PublisherTypeID", publisher.PublisherTypeID),
				sql.Named("PublisherName", publisher.PublisherName),
				sql.Named("DiscountID", publisher.DiscountID),
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

func UpdatePublisherByID(c *fiber.Ctx) error {
	publisherID := c.Params("id")
	var updated models.Publisher
	if err := c.BodyParser(&updated); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	query := `
		UPDATE tb_publisher
		SET publisher_type_id = @PublisherTypeID,
			publisher_name = COALESCE(NULLIF(@PublisherName, ''), publisher_name),
			discount_id = COALESCE(NULLIF(@DiscountID, ''), discount_id),
			contact_name1 = @ContactName1,
			contact_name2 = @ContactName2,
			email = @Email,
			phone1 = @Phone1,
			phone2 = @Phone2,
			address = @Address,
			district = @District,
			province = @Province,
			zip_code = COALESCE(NULLIF(@ZipCode, ''), zip_code),
			note = @Note,
			update_by = @UpdateBy
		WHERE publisher_code = @PublisherID AND is_delete = 0
	`
	err := executeTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("PublisherTypeID", updated.PublisherTypeID),
				sql.Named("PublisherName", updated.PublisherName),
				sql.Named("DiscountID", updated.DiscountID),
				sql.Named("ContactName1", updated.ContactName1),
				sql.Named("ContactName2", updated.ContactName2),
				sql.Named("Email", updated.Email),
				sql.Named("Phone1", updated.Phone1),
				sql.Named("Phone2", updated.Phone2),
				sql.Named("Address", updated.Address),
				sql.Named("District", updated.District),
				sql.Named("Province", updated.Province),
				sql.Named("ZipCode", updated.ZipCode),
				sql.Named("Note", updated.Note),
				sql.Named("UpdateBy", updated.UpdateBy),
				sql.Named("PublisherID", publisherID),
			)
			return err
		},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update publisher"})
	}
	return c.JSON(fiber.Map{"message": "Publisher updated successfully"})
}

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

func RemovePublisherByID(c *fiber.Ctx) error {
	publisherID := c.Params("id")
	query := `DELETE FROM tb_publisher WHERE publisher_code = @PublisherID`
	_, err := config.DB.Exec(query, sql.Named("PublisherID", publisherID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove publisher"})
	}
	return c.JSON(fiber.Map{"message": "Publisher removed successfully"})
}
