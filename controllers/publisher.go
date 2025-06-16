package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllPublisher(c *fiber.Ctx) error {
	query := `
		SELECT 
			p.publisher_code, p.publisher_name, 
			p.publisher_type_id, pt.type_name AS publisher_type_name, 
			p.contact_name1, p.contact_name2, 
			p.email, p.phone1, p.phone2, 
			p.address, p.district, p.province, p.zip_code, 
			p.note, 
			p.discount_id, d.discount_name, -- ✅ JOIN จาก tb_discount
			p.update_by, p.update_date, 
			p.id_status
		FROM tb_publisher p
		LEFT JOIN tb_publisher_type pt ON p.publisher_type_id = pt.publisher_type_id
		LEFT JOIN tb_discount d ON p.discount_id = d.discount_id -- ✅ เพิ่มบรรทัดนี้
		WHERE p.is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch publishers",
			Data:    nil,
		})
	}
	defer rows.Close()

	var publishers []models.Publisher
	for rows.Next() {
		var p models.Publisher
		if err := rows.Scan(
			&p.PublisherCode, &p.PublisherName,
			&p.PublisherTypeID, &p.PublisherTypeName,
			&p.ContactName1, &p.ContactName2,
			&p.Email, &p.Phone1, &p.Phone2,
			&p.Address, &p.District, &p.Province, &p.ZipCode,
			&p.Note,
			&p.DiscountID, &p.DiscountName,
			&p.UpdateBy, &p.UpdateDate,
			&p.IDStatus,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		publishers = append(publishers, p)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    publishers,
	})
}

func SelectPagePublisher(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT 
			p.publisher_code, p.publisher_name, 
			p.publisher_type_id, pt.type_name AS publisher_type_name, 
			p.contact_name1, p.contact_name2, 
			p.email, p.phone1, p.phone2, 
			p.address, p.district, p.province, p.zip_code, 
			p.note, 
			p.discount_id, d.discount_name,
			p.update_by, p.update_date, p.id_status
		FROM tb_publisher p
		LEFT JOIN tb_publisher_type pt ON p.publisher_type_id = pt.publisher_type_id
		LEFT JOIN tb_discount d ON p.discount_id = d.discount_id
		WHERE p.is_delete = 0
		ORDER BY p.update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch publishers",
			Data:    nil,
		})
	}
	defer rows.Close()

	var publishers []models.Publisher
	for rows.Next() {
		var p models.Publisher
		if err := rows.Scan(
			&p.PublisherCode, &p.PublisherName,
			&p.PublisherTypeID, &p.PublisherTypeName,
			&p.ContactName1, &p.ContactName2,
			&p.Email, &p.Phone1, &p.Phone2,
			&p.Address, &p.District, &p.Province, &p.ZipCode,
			&p.Note,
			&p.DiscountID, &p.DiscountName,
			&p.UpdateBy, &p.UpdateDate, &p.IDStatus,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		publishers = append(publishers, p)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_publisher WHERE is_delete = 0`).Scan(&total)
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
			"page":      page,
			"limit":     limit,
			"total":     total,
			"publisher": publishers,
		},
	})
}

func SelectPublisherByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT 
			p.publisher_code, p.publisher_name, p.publisher_type_id, 
			p.discount_id, d.discount_name,
			p.contact_name1, p.contact_name2,
			p.email, p.phone1, p.phone2, 
			p.address, p.district, p.province, p.zip_code,
			p.note, 
			p.update_by, p.update_date, p.id_status
		FROM tb_publisher p
		LEFT JOIN tb_discount d ON p.discount_id = d.discount_id -- ✅ JOIN
		WHERE p.publisher_code = @ID AND p.is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var p models.Publisher
	if err := row.Scan(
		&p.PublisherCode, &p.PublisherName, &p.PublisherTypeID,
		&p.DiscountID, &p.DiscountName,
		&p.ContactName1, &p.ContactName2,
		&p.Email, &p.Phone1, &p.Phone2,
		&p.Address, &p.District, &p.Province, &p.ZipCode,
		&p.Note,
		&p.UpdateBy, &p.UpdateDate, &p.IDStatus,
	); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Publisher not found",
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
		Data:    p,
	})
}

func SelectPublisherByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT 
			p.publisher_code, 
			p.publisher_type_id, pt.type_name,
			p.publisher_name, 
			p.contact_name1, p.contact_name2,
			p.email, p.phone1, p.phone2,
			p.address, p.district, p.province, p.zip_code,
			p.note, 
			p.discount_id, d.discount_name, -- ✅ join discount_name
			p.update_by, p.update_date, p.id_status
		FROM tb_publisher p
		LEFT JOIN tb_publisher_type pt ON p.publisher_type_id = pt.publisher_type_id
		LEFT JOIN tb_discount d ON p.discount_id = d.discount_id
		WHERE p.publisher_name LIKE '%' + @Name + '%' AND p.is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch publishers",
			Data:    nil,
		})
	}
	defer rows.Close()

	var publishers []models.Publisher
	for rows.Next() {
		var p models.Publisher
		if err := rows.Scan(
			&p.PublisherCode,
			&p.PublisherTypeID, &p.PublisherTypeName, // ✅ เพิ่มรับ type_name
			&p.PublisherName,
			&p.ContactName1, &p.ContactName2,
			&p.Email, &p.Phone1, &p.Phone2,
			&p.Address, &p.District, &p.Province, &p.ZipCode,
			&p.Note,
			&p.DiscountID, &p.DiscountName, // ✅ รับ discount_name
			&p.UpdateBy, &p.UpdateDate, &p.IDStatus,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		publishers = append(publishers, p)
	}

	if len(publishers) == 0 {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "No matching publisher found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    publishers,
	})
}

func InsertPublisher(c *fiber.Ctx) error {
	var p models.Publisher
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request",
			Data:    nil,
		})
	}
	query := `
		INSERT INTO tb_publisher (publisher_type_id, publisher_name, discount_id, contact_name1, contact_name2,
			email, phone1, phone2, address, district, province, zip_code, note, update_by)
		VALUES (@TypeID, @Name, @DiscountID, @Contact1, @Contact2,
			@Email, @Phone1, @Phone2, @Address, @District, @Province, @Zip, @Note, @UpdateBy)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeID", p.PublisherTypeID),
				sql.Named("Name", p.PublisherName),
				sql.Named("DiscountID", p.DiscountID),
				sql.Named("Contact1", p.ContactName1),
				sql.Named("Contact2", p.ContactName2),
				sql.Named("Email", p.Email),
				sql.Named("Phone1", p.Phone1),
				sql.Named("Phone2", p.Phone2),
				sql.Named("Address", p.Address),
				sql.Named("District", p.District),
				sql.Named("Province", p.Province),
				sql.Named("Zip", p.ZipCode),
				sql.Named("Note", p.Note),
				sql.Named("UpdateBy", p.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to insert publisher",
			Data:    nil,
		})
	}
	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Publisher added successfully",
		Data:    nil,
	})
}

func UpdatePublisherByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var p models.Publisher
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}
	query := `
		UPDATE tb_publisher
		SET publisher_type_id = @TypeID,
			publisher_name = COALESCE(NULLIF(@Name, ''), publisher_name),
			discount_id = COALESCE(NULLIF(@DiscountID, ''), discount_id),
			contact_name1 = @Contact1,
			contact_name2 = @Contact2,
			email = @Email,
			phone1 = @Phone1,
			phone2 = @Phone2,
			address = @Address,
			district = @District,
			province = @Province,
			zip_code = COALESCE(NULLIF(@Zip, ''), zip_code),
			note = @Note,
			update_by = @UpdateBy
		WHERE publisher_code = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeID", p.PublisherTypeID),
				sql.Named("Name", p.PublisherName),
				sql.Named("DiscountID", p.DiscountID),
				sql.Named("Contact1", p.ContactName1),
				sql.Named("Contact2", p.ContactName2),
				sql.Named("Email", p.Email),
				sql.Named("Phone1", p.Phone1),
				sql.Named("Phone2", p.Phone2),
				sql.Named("Address", p.Address),
				sql.Named("District", p.District),
				sql.Named("Province", p.Province),
				sql.Named("Zip", p.ZipCode),
				sql.Named("Note", p.Note),
				sql.Named("UpdateBy", p.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update publisher",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Publisher updated successfully",
		Data:    nil,
	})
}

func DeletePublisherByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_publisher
		SET is_delete = 1,
			update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE publisher_code = @ID
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
			Message: "Failed to soft delete publisher",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Publisher marked as deleted",
		Data:    nil,
	})
}

func RemovePublisherByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_publisher WHERE publisher_code = @ID`
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
			Message: "Failed to hard delete publisher",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Publisher removed successfully",
		Data:    nil,
	})
}
