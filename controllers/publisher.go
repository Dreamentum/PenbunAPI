package controllers

import (
	"github.com/gofiber/fiber/v2"

	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
)

func SelectAllPublisher(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT autoID, publisher_id, publisher_name, address, phone, is_active, update_by, update_date, is_delete FROM tb_publisher WHERE is_delete = 0")
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	defer rows.Close()

	var items []models.Publisher
	for rows.Next() {
		var item models.Publisher
		if err := rows.Scan(&item.AutoID, &item.PublisherID, &item.PublisherName, &item.Address, &item.Phone, &item.IsActive, &item.UpdateBy, &item.UpdateDate, &item.IsDelete); err != nil {
			return utils.ErrorResponse(c, err.Error())
		}
		items = append(items, item)
	}
	return utils.SuccessResponse(c, "Publisher list retrieved", items)
}

func SelectPagePublisher(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	rows, err := config.DB.Query("SELECT autoID, publisher_id, publisher_name, address, phone, is_active, update_by, update_date, is_delete FROM tb_publisher WHERE is_delete = 0 ORDER BY update_date DESC OFFSET ? ROWS FETCH NEXT ? ROWS ONLY", offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	defer rows.Close()

	var items []models.Publisher
	for rows.Next() {
		var item models.Publisher
		if err := rows.Scan(&item.AutoID, &item.PublisherID, &item.PublisherName, &item.Address, &item.Phone, &item.IsActive, &item.UpdateBy, &item.UpdateDate, &item.IsDelete); err != nil {
			return utils.ErrorResponse(c, err.Error())
		}
		items = append(items, item)
	}
	return utils.SuccessResponse(c, "Publisher page retrieved", items)
}

func SelectPublisherByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.Publisher
	err := config.DB.QueryRow("SELECT autoID, publisher_id, publisher_name, address, phone, is_active, update_by, update_date, is_delete FROM tb_publisher WHERE publisher_id = ? AND is_delete = 0", id).
		Scan(&item.AutoID, &item.PublisherID, &item.PublisherName, &item.Address, &item.Phone, &item.IsActive, &item.UpdateBy, &item.UpdateDate, &item.IsDelete)
	if err != nil {
		return utils.FailResponse(c, "Publisher not found")
	}
	return utils.SuccessResponse(c, "Publisher found", item)
}

func SelectPublisherByName(c *fiber.Ctx) error {
	name := c.Params("name")
	rows, err := config.DB.Query("SELECT autoID, publisher_id, publisher_name, address, phone, is_active, update_by, update_date, is_delete FROM tb_publisher WHERE publisher_name LIKE '%' + ? + '%' AND is_delete = 0", name)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	defer rows.Close()

	var items []models.Publisher
	for rows.Next() {
		var item models.Publisher
		if err := rows.Scan(&item.AutoID, &item.PublisherID, &item.PublisherName, &item.Address, &item.Phone, &item.IsActive, &item.UpdateBy, &item.UpdateDate, &item.IsDelete); err != nil {
			return utils.ErrorResponse(c, err.Error())
		}
		items = append(items, item)
	}
	return utils.SuccessResponse(c, "Publisher search results", items)
}

func InsertPublisher(c *fiber.Ctx) error {
	var item models.Publisher
	if err := c.BodyParser(&item); err != nil {
		return utils.FailResponse(c, "Invalid request body")
	}

	if item.PublisherName == "" {
		return utils.FailResponse(c, "Publisher name is required")
	}

	steps := []utils.TransactionStep{
		{Name: "InsertPublisher", Query: "INSERT INTO tb_publisher (publisher_name, address, phone, update_by) VALUES (?, ?, ?, ?)",
			Args: []interface{}{item.PublisherName, item.Address, item.Phone, item.UpdateBy}},
	}
	if err := utils.ExecuteTransaction(steps); err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, "Publisher added successfully", fiber.Map{"publisher_name": item.PublisherName})
}

func UpdatePublisherByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.Publisher
	if err := c.BodyParser(&item); err != nil {
		return utils.FailResponse(c, "Invalid request body")
	}

	steps := []utils.TransactionStep{
		{Name: "UpdatePublisher", Query: "UPDATE tb_publisher SET publisher_name = COALESCE(NULLIF(?, ''), publisher_name), address = COALESCE(?, address), phone = COALESCE(?, phone), update_by = ? WHERE publisher_id = ? AND is_delete = 0",
			Args: []interface{}{item.PublisherName, item.Address, item.Phone, item.UpdateBy, id}},
	}
	if err := utils.ExecuteTransaction(steps); err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, "Publisher updated successfully", fiber.Map{"publisher_id": id})
}

func DeletePublisherByID(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Query("user", "UNKNOWN")

	steps := []utils.TransactionStep{
		{Name: "DeletePublisher", Query: "UPDATE tb_publisher SET is_delete = 1, update_by = ? WHERE publisher_id = ?", Args: []interface{}{username, id}},
	}
	if err := utils.ExecuteTransaction(steps); err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, "Publisher deleted successfully", nil)
}

func RemovePublisherByID(c *fiber.Ctx) error {
	id := c.Params("id")

	steps := []utils.TransactionStep{
		{Name: "RemovePublisher", Query: "DELETE FROM tb_publisher WHERE publisher_id = ?", Args: []interface{}{id}},
	}
	if err := utils.ExecuteTransaction(steps); err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, "Publisher removed permanently", nil)
}
