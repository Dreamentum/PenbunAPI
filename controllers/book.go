package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

// 1. Select All
func SelectAllBooks(c *fiber.Ctx) error {
	query := `
		SELECT book_code, book_name, book_isbn, book_barcode,
		       book_type_id, book_format_type_id, book_external_code,
		       publisher_code, book_price, book_discount, description, note,
		       update_by, update_date, id_status
		FROM tb_book
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to fetch books", Data: nil,
		})
	}
	defer rows.Close()

	var list []models.Book
	for rows.Next() {
		var item models.Book
		if err := rows.Scan(
			&item.BookCode, &item.BookName, &item.BookISBN, &item.BookBarcode,
			&item.BookTypeID, &item.BookFormatTypeID, &item.BookExternalCode,
			&item.PublisherCode, &item.BookPrice, &item.BookDiscount,
			&item.Description, &item.Note, &item.UpdateBy, &item.UpdateDate, &item.IDStatus,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status: "error", Message: "Failed to read data", Data: nil,
			})
		}
		list = append(list, item)
	}
	return c.JSON(models.ApiResponse{
		Status: "success", Message: "", Data: list,
	})
}

// 2. Select Paging
func SelectPageBooks(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT book_code, book_name, book_isbn, book_barcode,
		       book_type_id, book_format_type_id, book_external_code,
		       publisher_code, book_price, book_discount, description, note,
		       update_by, update_date, id_status
		FROM tb_book
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to fetch books", Data: nil,
		})
	}
	defer rows.Close()

	var list []models.Book
	for rows.Next() {
		var item models.Book
		if err := rows.Scan(
			&item.BookCode, &item.BookName, &item.BookISBN, &item.BookBarcode,
			&item.BookTypeID, &item.BookFormatTypeID, &item.BookExternalCode,
			&item.PublisherCode, &item.BookPrice, &item.BookDiscount,
			&item.Description, &item.Note, &item.UpdateBy, &item.UpdateDate, &item.IDStatus,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status: "error", Message: "Failed to read data", Data: nil,
			})
		}
		list = append(list, item)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_book WHERE is_delete = 0`).Scan(&total)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to count records", Data: nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status: "success",
		Data: fiber.Map{
			"page": page, "limit": limit, "total": total,
			"books": list,
		},
	})
}

// 3. Select By ID
func SelectBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT book_code, book_name, book_isbn, book_barcode,
		       book_type_id, book_format_type_id, book_external_code,
		       publisher_code, book_price, book_discount, description, note,
		       update_by, update_date, id_status
		FROM tb_book
		WHERE book_code = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))
	var b models.Book
	if err := row.Scan(
		&b.BookCode, &b.BookName, &b.BookISBN, &b.BookBarcode,
		&b.BookTypeID, &b.BookFormatTypeID, &b.BookExternalCode,
		&b.PublisherCode, &b.BookPrice, &b.BookDiscount,
		&b.Description, &b.Note, &b.UpdateBy, &b.UpdateDate, &b.IDStatus,
	); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status: "error", Message: "Book not found", Data: nil,
			})
		}
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to read book", Data: nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status: "success", Message: "", Data: b,
	})
}

// 4. Select By Name (LIKE)
func SelectBookByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT book_code, book_name, book_isbn, book_barcode,
		       book_type_id, book_format_type_id, book_external_code,
		       publisher_code, book_price, book_discount, description, note,
		       update_by, update_date, id_status
		FROM tb_book
		WHERE book_name LIKE '%' + @Name + '%' AND is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to search books", Data: nil,
		})
	}
	defer rows.Close()

	var list []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(
			&b.BookCode, &b.BookName, &b.BookISBN, &b.BookBarcode,
			&b.BookTypeID, &b.BookFormatTypeID, &b.BookExternalCode,
			&b.PublisherCode, &b.BookPrice, &b.BookDiscount,
			&b.Description, &b.Note, &b.UpdateBy, &b.UpdateDate, &b.IDStatus,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status: "error", Message: "Failed to read data", Data: nil,
			})
		}
		list = append(list, b)
	}
	return c.JSON(models.ApiResponse{
		Status: "success", Message: "", Data: list,
	})
}

// 5. Insert
func InsertBook(c *fiber.Ctx) error {
	var b models.Book
	if err := c.BodyParser(&b); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status: "error", Message: "Invalid request body", Data: nil,
		})
	}

	query := `
		INSERT INTO tb_book (
			book_name, book_isbn, book_barcode, book_type_id, book_format_type_id,
			book_external_code, publisher_code, book_price, book_discount,
			description, note, update_by
		)
		VALUES (@Name, @ISBN, @Barcode, @TypeID, @FormatTypeID,
				@ExternalCode, @Publisher, @Price, @Discount,
				@Desc, @Note, @UpdateBy)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("Name", b.BookName),
				sql.Named("ISBN", b.BookISBN),
				sql.Named("Barcode", b.BookBarcode),
				sql.Named("TypeID", b.BookTypeID),
				sql.Named("FormatTypeID", b.BookFormatTypeID),
				sql.Named("ExternalCode", b.BookExternalCode),
				sql.Named("Publisher", b.PublisherCode),
				sql.Named("Price", b.BookPrice),
				sql.Named("Discount", b.BookDiscount),
				sql.Named("Desc", b.Description),
				sql.Named("Note", b.Note),
				sql.Named("UpdateBy", b.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to insert book", Data: nil,
		})
	}
	return c.Status(201).JSON(models.ApiResponse{
		Status: "success", Message: "Book added successfully", Data: nil,
	})
}

// 6. Update
func UpdateBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var b models.Book
	if err := c.BodyParser(&b); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status: "error", Message: "Invalid request body", Data: nil,
		})
	}

	query := `
		UPDATE tb_book
		SET book_name = COALESCE(NULLIF(@Name, ''), book_name),
			book_isbn = @ISBN,
			book_barcode = @Barcode,
			book_type_id = @TypeID,
			book_format_type_id = @FormatTypeID,
			book_external_code = @ExternalCode,
			publisher_code = @Publisher,
			book_price = @Price,
			book_discount = @Discount,
			description = @Desc,
			note = @Note,
			update_by = @UpdateBy
		WHERE book_code = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("Name", b.BookName),
				sql.Named("ISBN", b.BookISBN),
				sql.Named("Barcode", b.BookBarcode),
				sql.Named("TypeID", b.BookTypeID),
				sql.Named("FormatTypeID", b.BookFormatTypeID),
				sql.Named("ExternalCode", b.BookExternalCode),
				sql.Named("Publisher", b.PublisherCode),
				sql.Named("Price", b.BookPrice),
				sql.Named("Discount", b.BookDiscount),
				sql.Named("Desc", b.Description),
				sql.Named("Note", b.Note),
				sql.Named("UpdateBy", b.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to update book", Data: nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status: "success", Message: "Book updated successfully", Data: nil,
	})
}

// 7. Delete (Soft)
func DeleteBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_book
		SET is_delete = 1,
			update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE book_code = @ID
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
			Status: "error", Message: "Failed to delete book", Data: nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status: "success", Message: "Book deleted successfully", Data: nil,
	})
}

// 8. Remove (Hard)
func RemoveBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_book WHERE book_code = @ID`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, sql.Named("ID", id))
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to remove book", Data: nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status: "success", Message: "Book removed successfully", Data: nil,
	})
}
