package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllBookTypes(c *fiber.Ctx) error {
	query := `
		SELECT book_type_id, type_name, description, update_by, update_date, id_status, is_delete
		FROM tb_book_type
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch book types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.BookType
	for rows.Next() {
		var bt models.BookType
		if err := rows.Scan(&bt.BookTypeID, &bt.TypeName, &bt.Description, &bt.UpdateBy, &bt.UpdateDate, &bt.IDStatus, &bt.IsDelete); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		list = append(list, bt)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    list,
	})
}

func SelectPageBookTypes(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT book_type_id, type_name, description, update_by, update_date, id_status, is_delete
		FROM tb_book_type
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch book types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.BookType
	for rows.Next() {
		var bt models.BookType
		if err := rows.Scan(&bt.BookTypeID, &bt.TypeName, &bt.Description, &bt.UpdateBy, &bt.UpdateDate, &bt.IDStatus, &bt.IsDelete); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		list = append(list, bt)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_book_type WHERE is_delete = 0`).Scan(&total)
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
			"page":     page,
			"limit":    limit,
			"total":    total,
			"bookType": list,
		},
	})
}

func SelectBookTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT book_type_id, type_name, description, update_by, update_date, id_status, is_delete
		FROM tb_book_type
		WHERE book_type_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))
	var bt models.BookType
	if err := row.Scan(&bt.BookTypeID, &bt.TypeName, &bt.Description, &bt.UpdateBy, &bt.UpdateDate, &bt.IDStatus, &bt.IsDelete); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Book type not found",
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
		Data:    bt,
	})
}

func SelectBookTypeByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT book_type_id, type_name, description, update_by, update_date, id_status, is_delete
		FROM tb_book_type
		WHERE type_name LIKE '%' + @Name + '%' AND is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch book types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.BookType
	for rows.Next() {
		var bt models.BookType
		if err := rows.Scan(&bt.BookTypeID, &bt.TypeName, &bt.Description, &bt.UpdateBy, &bt.UpdateDate, &bt.IDStatus, &bt.IsDelete); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		list = append(list, bt)
	}

	if len(list) == 0 {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "No matching book type found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    list,
	})
}

func InsertBookType(c *fiber.Ctx) error {
	var bt models.BookType
	if err := c.BodyParser(&bt); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}
	query := `
		INSERT INTO tb_book_type (type_name, description, update_by)
		VALUES (@TypeName, @Description, @UpdateBy)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeName", bt.TypeName),
				sql.Named("Description", bt.Description),
				sql.Named("UpdateBy", bt.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to insert book type",
			Data:    nil,
		})
	}
	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Book type added successfully",
		Data:    nil,
	})
}

func UpdateBookTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var bt models.BookType
	if err := c.BodyParser(&bt); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}
	query := `
		UPDATE tb_book_type
		SET type_name = COALESCE(NULLIF(@TypeName, ''), type_name),
			description = @Description,
			update_by = @UpdateBy
		WHERE book_type_id = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeName", bt.TypeName),
				sql.Named("Description", bt.Description),
				sql.Named("UpdateBy", bt.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update book type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Book type updated successfully",
		Data:    nil,
	})
}

func DeleteBookTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_book_type
		SET is_delete = 1,
		    update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE book_type_id = @ID
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
			Message: "Failed to soft delete book type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Book type marked as deleted",
		Data:    nil,
	})
}

func RemoveBookTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_book_type WHERE book_type_id = @ID`
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
			Message: "Failed to hard delete book type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Book type removed successfully",
		Data:    nil,
	})
}
