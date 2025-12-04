package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllPublisherTypes(c *fiber.Ctx) error {
	query := `
		SELECT publisher_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_publisher_type
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch publisher types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var result []models.PublisherType
	for rows.Next() {
		var pt models.PublisherType
		if err := rows.Scan(&pt.PublisherTypeID, &pt.TypeName, &pt.Description, &pt.UpdateBy, &pt.UpdateDate, &pt.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		result = append(result, pt)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    result,
	})
}

func SelectPagePublisherTypes(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT publisher_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_publisher_type
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch publisher types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var result []models.PublisherType
	for rows.Next() {
		var pt models.PublisherType
		if err := rows.Scan(&pt.PublisherTypeID, &pt.TypeName, &pt.Description, &pt.UpdateBy, &pt.UpdateDate, &pt.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		result = append(result, pt)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_publisher_type WHERE is_delete = 0`).Scan(&total)
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
			"page":          page,
			"limit":         limit,
			"total":         total,
			"publisherType": result,
		},
	})
}

func SelectPublisherTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT publisher_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_publisher_type
		WHERE publisher_type_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var pt models.PublisherType
	if err := row.Scan(&pt.PublisherTypeID, &pt.TypeName, &pt.Description, &pt.UpdateBy, &pt.UpdateDate, &pt.IDStatus); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Publisher type not found",
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
		Data:    pt,
	})
}

func SelectPublisherTypeByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT publisher_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_publisher_type
		WHERE type_name LIKE '%' + @Name + '%' AND is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch publisher types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var result []models.PublisherType
	for rows.Next() {
		var pt models.PublisherType
		if err := rows.Scan(&pt.PublisherTypeID, &pt.TypeName, &pt.Description, &pt.UpdateBy, &pt.UpdateDate, &pt.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		result = append(result, pt)
	}

	if len(result) == 0 {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "No matching publisher type found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    result,
	})
}

func InsertPublisherType(c *fiber.Ctx) error {
	var pt models.PublisherType
	if err := c.BodyParser(&pt); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		INSERT INTO tb_publisher_type (type_name, description, update_by)
		VALUES (@TypeName, @Description, @UpdateBy)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeName", pt.TypeName),
				sql.Named("Description", pt.Description),
				sql.Named("UpdateBy", pt.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to insert publisher type",
			Data:    nil,
		})
	}

	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Publisher type added successfully",
		Data:    nil,
	})
}

func UpdatePublisherTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var pt models.PublisherType
	if err := c.BodyParser(&pt); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		UPDATE tb_publisher_type
		SET type_name = COALESCE(NULLIF(@TypeName, ''), type_name),
			description = @Description,
			update_by = @UpdateBy
		WHERE publisher_type_id = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeName", pt.TypeName),
				sql.Named("Description", pt.Description),
				sql.Named("UpdateBy", pt.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update publisher type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Publisher type updated successfully",
		Data:    nil,
	})
}

func DeletePublisherTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_publisher_type
		SET is_delete = 1,
			update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE publisher_type_id = @ID
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
			Message: "Failed to soft delete publisher type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Publisher type marked as deleted",
		Data:    nil,
	})
}

func RemovePublisherTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_publisher_type WHERE publisher_type_id = @ID`
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
			Message: "Failed to hard delete publisher type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Publisher type removed successfully",
		Data:    nil,
	})
}
