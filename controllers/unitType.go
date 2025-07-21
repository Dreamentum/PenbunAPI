package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SelectAllUnitType(c *fiber.Ctx) error {
	query := `
		SELECT unit_type_id, unit_type_name, description, update_by, update_date, id_status
		FROM tb_unit_type
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch unit types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var result []models.UnitType
	for rows.Next() {
		var ut models.UnitType
		if err := rows.Scan(&ut.UnitTypeID, &ut.UnitTypeName, &ut.Description, &ut.UpdateBy, &ut.UpdateDate, &ut.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		result = append(result, ut)
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    result,
	})
}

func SelectPageUnitType(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT unit_type_id, unit_type_name, description, update_by, update_date, id_status
		FROM tb_unit_type
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch unit types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var result []models.UnitType
	for rows.Next() {
		var ut models.UnitType
		if err := rows.Scan(&ut.UnitTypeID, &ut.UnitTypeName, &ut.Description, &ut.UpdateBy, &ut.UpdateDate, &ut.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		result = append(result, ut)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_unit_type WHERE is_delete = 0`).Scan(&total)
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
			"unitTypes": result,
		},
	})
}

func SelectUnitTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT unit_type_id, unit_type_name, description, update_by, update_date, id_status
		FROM tb_unit_type
		WHERE unit_type_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))

	var ut models.UnitType
	if err := row.Scan(&ut.UnitTypeID, &ut.UnitTypeName, &ut.Description, &ut.UpdateBy, &ut.UpdateDate, &ut.IDStatus); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Unit type not found",
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
		Data:    ut,
	})
}

func SelectUnitTypeByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT unit_type_id, unit_type_name, description, update_by, update_date, id_status
		FROM tb_unit_type
		WHERE unit_type_name LIKE '%' + @Name + '%' AND is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to fetch unit types",
			Data:    nil,
		})
	}
	defer rows.Close()

	var result []models.UnitType
	for rows.Next() {
		var ut models.UnitType
		if err := rows.Scan(&ut.UnitTypeID, &ut.UnitTypeName, &ut.Description, &ut.UpdateBy, &ut.UpdateDate, &ut.IDStatus); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status:  "error",
				Message: "Failed to read data",
				Data:    nil,
			})
		}
		result = append(result, ut)
	}

	if len(result) == 0 {
		return c.Status(404).JSON(models.ApiResponse{
			Status:  "error",
			Message: "No matching unit type found",
			Data:    nil,
		})
	}

	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "",
		Data:    result,
	})
}

func InsertUnitType(c *fiber.Ctx) error {
	var ut models.UnitType
	if err := c.BodyParser(&ut); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		INSERT INTO tb_unit_type (unit_type_name, description, update_by)
		VALUES (@Name, @Desc, @By)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("Name", ut.UnitTypeName),
				sql.Named("Desc", ut.Description),
				sql.Named("By", ut.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to insert unit type",
			Data:    nil,
		})
	}

	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Unit type added successfully",
		Data:    nil,
	})
}

func UpdateUnitTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var ut models.UnitType
	if err := c.BodyParser(&ut); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	query := `
		UPDATE tb_unit_type
		SET unit_type_name = COALESCE(NULLIF(@Name, ''), unit_type_name),
			description = @Desc,
			update_by = @By
		WHERE unit_type_id = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("Name", ut.UnitTypeName),
				sql.Named("Desc", ut.Description),
				sql.Named("By", ut.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update unit type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Unit type updated successfully",
		Data:    nil,
	})
}

func DeleteUnitTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_unit_type
		SET is_delete = 1,
			update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE unit_type_id = @ID
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
			Message: "Failed to soft delete unit type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Unit type marked as deleted",
		Data:    nil,
	})
}

func RemoveUnitTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_unit_type WHERE unit_type_id = @ID`
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
			Message: "Failed to hard delete unit type",
			Data:    nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Unit type removed successfully",
		Data:    nil,
	})
}
