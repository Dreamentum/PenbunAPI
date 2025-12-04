package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllCustomerTypes(c *fiber.Ctx) error {
	query := `
		SELECT customer_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_customer_type
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch customer types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.CustomerType
	for rows.Next() {
		var item models.CustomerType
		if err := rows.Scan(&item.CustomerTypeID, &item.TypeName, &item.Description, &item.UpdateBy, &item.UpdateDate, &item.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		list = append(list, item)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    list,
	})
}

func SelectPageCustomerTypes(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT customer_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_customer_type
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch customer types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var list []models.CustomerType
	for rows.Next() {
		var item models.CustomerType
		if err := rows.Scan(&item.CustomerTypeID, &item.TypeName, &item.Description, &item.UpdateBy, &item.UpdateDate, &item.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		list = append(list, item)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_customer_type WHERE is_delete = 0`).Scan(&total)
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
			"page":         page,
			"limit":        limit,
			"total":        total,
			"customerType": list,
		},
	})
}

func SelectCustomerTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT customer_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_customer_type
		WHERE customer_type_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var item models.CustomerType
	if err := row.Scan(&item.CustomerTypeID, &item.TypeName, &item.Description, &item.UpdateBy, &item.UpdateDate, &item.IDStatus); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Customer type not found",
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
		Data:    item,
	})
}

func SelectCustomerTypeByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT customer_type_id, type_name, description, update_by, update_date, id_status
		FROM tb_customer_type
		WHERE type_name LIKE '%' + @Name + '%' AND is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch customer types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var results []models.CustomerType
	for rows.Next() {
		var item models.CustomerType
		if err := rows.Scan(&item.CustomerTypeID, &item.TypeName, &item.Description, &item.UpdateBy, &item.UpdateDate, &item.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		results = append(results, item)
	}

	if len(results) == 0 {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "No matching customer type found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    results,
	})
}

func InsertCustomerType(c *fiber.Ctx) error {
	var item models.CustomerType
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		INSERT INTO tb_customer_type (type_name, description, update_by)
		VALUES (@TypeName, @Description, @UpdateBy)
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeName", item.TypeName),
				sql.Named("Description", item.Description),
				sql.Named("UpdateBy", item.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to insert customer type",
			Data:    nil,
		})
	}

	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Customer type added successfully",
		Data:    nil,
	})
}

func UpdateCustomerTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.CustomerType
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		UPDATE tb_customer_type
		SET type_name = COALESCE(NULLIF(@TypeName, ''), type_name),
			description = @Description,
			update_by = @UpdateBy
		WHERE customer_type_id = @ID AND is_delete = 0
	`

	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("TypeName", item.TypeName),
				sql.Named("Description", item.Description),
				sql.Named("UpdateBy", item.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update customer type",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Customer type updated successfully",
		Data:    nil,
	})
}

func DeleteCustomerTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_customer_type
		SET is_delete = 1,
		    id_status = 0,
		    update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE customer_type_id = @ID
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
			Message: "Failed to soft delete customer type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Customer type marked as deleted",
		Data:    nil,
	})
}

func RemoveCustomerTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_customer_type WHERE customer_type_id = @ID`
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
			Message: "Failed to hard delete customer type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Customer type removed successfully",
		Data:    nil,
	})
}
