package controllers

import (
	"github.com/gofiber/fiber/v2"

	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
)

func SelectAllPublisherType(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT autoID, publisher_type_id, type_name, description, is_active, update_by, update_date, is_delete FROM tb_publisher_type WHERE is_delete = 0")
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	defer rows.Close()

	var items []models.PublisherType
	for rows.Next() {
		var item models.PublisherType
		if err := rows.Scan(&item.AutoID, &item.PublisherTypeID, &item.TypeName, &item.Description, &item.IsActive, &item.UpdateBy, &item.UpdateDate, &item.IsDelete); err != nil {
			return utils.ErrorResponse(c, err.Error())
		}
		items = append(items, item)
	}
	return utils.SuccessResponse(c, "Publisher type list retrieved", items)
}

func SelectPagePublisherType(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	rows, err := config.DB.Query("SELECT autoID, publisher_type_id, type_name, description, is_active, update_by, update_date, is_delete FROM tb_publisher_type WHERE is_delete = 0 ORDER BY update_date DESC OFFSET ? ROWS FETCH NEXT ? ROWS ONLY", offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	defer rows.Close()

	var items []models.PublisherType
	for rows.Next() {
		var item models.PublisherType
		if err := rows.Scan(&item.AutoID, &item.PublisherTypeID, &item.TypeName, &item.Description, &item.IsActive, &item.UpdateBy, &item.UpdateDate, &item.IsDelete); err != nil {
			return utils.ErrorResponse(c, err.Error())
		}
		items = append(items, item)
	}
	return utils.SuccessResponse(c, "Publisher type page retrieved", items)
}

func SelectPublisherTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.PublisherType
	err := config.DB.QueryRow("SELECT autoID, publisher_type_id, type_name, description, is_active, update_by, update_date, is_delete FROM tb_publisher_type WHERE publisher_type_id = ? AND is_delete = 0", id).
		Scan(&item.AutoID, &item.PublisherTypeID, &item.TypeName, &item.Description, &item.IsActive, &item.UpdateBy, &item.UpdateDate, &item.IsDelete)
	if err != nil {
		return utils.FailResponse(c, "Publisher type not found")
	}
	return utils.SuccessResponse(c, "Publisher type found", item)
}

func SelectPublisherTypeByName(c *fiber.Ctx) error {
	name := c.Params("name")
	rows, err := config.DB.Query("SELECT autoID, publisher_type_id, type_name, description, is_active, update_by, update_date, is_delete FROM tb_publisher_type WHERE type_name LIKE '%' + ? + '%' AND is_delete = 0", name)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	defer rows.Close()

	var items []models.PublisherType
	for rows.Next() {
		var item models.PublisherType
		if err := rows.Scan(&item.AutoID, &item.PublisherTypeID, &item.TypeName, &item.Description, &item.IsActive, &item.UpdateBy, &item.UpdateDate, &item.IsDelete); err != nil {
			return utils.ErrorResponse(c, err.Error())
		}
		items = append(items, item)
	}
	return utils.SuccessResponse(c, "Publisher type search results", items)
}

func InsertPublisherType(c *fiber.Ctx) error {
	var item models.PublisherType
	if err := c.BodyParser(&item); err != nil {
		return utils.FailResponse(c, "Invalid request body")
	}
	if item.TypeName == "" {
		return utils.FailResponse(c, "Type name is required")
	}

	steps := []utils.TransactionStep{
		{Name: "InsertPublisherType", Query: "INSERT INTO tb_publisher_type (type_name, description, update_by) VALUES (?, ?, ?)",
			Args: []interface{}{item.TypeName, item.Description, item.UpdateBy}},
	}
	if err := utils.ExecuteTransaction(steps); err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, "Publisher type added successfully", fiber.Map{"type_name": item.TypeName})
}

func UpdatePublisherTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.PublisherType
	if err := c.BodyParser(&item); err != nil {
		return utils.FailResponse(c, "Invalid request body")
	}

	steps := []utils.TransactionStep{
		{Name: "UpdatePublisherType", Query: "UPDATE tb_publisher_type SET type_name = COALESCE(NULLIF(?, ''), type_name), description = COALESCE(?, description), update_by = ? WHERE publisher_type_id = ? AND is_delete = 0",
			Args: []interface{}{item.TypeName, item.Description, item.UpdateBy, id}},
	}
	if err := utils.ExecuteTransaction(steps); err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, "Publisher type updated successfully", fiber.Map{"publisher_type_id": id})
}

func DeletePublisherTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Query("user", "UNKNOWN")
	steps := []utils.TransactionStep{
		{Name: "DeletePublisherType", Query: "UPDATE tb_publisher_type SET is_delete = 1, update_by = ? WHERE publisher_type_id = ?", Args: []interface{}{username, id}},
	}
	if err := utils.ExecuteTransaction(steps); err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, "Publisher type deleted successfully", nil)
}

func RemovePublisherTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	steps := []utils.TransactionStep{
		{Name: "RemovePublisherType", Query: "DELETE FROM tb_publisher_type WHERE publisher_type_id = ?", Args: []interface{}{id}},
	}
	if err := utils.ExecuteTransaction(steps); err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, "Publisher type removed permanently", nil)
}
