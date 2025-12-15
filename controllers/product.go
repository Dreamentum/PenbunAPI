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
func SelectAllProducts(c *fiber.Ctx) error {
	query := `
		SELECT autoID, prefix, product_id, product_name_th, product_name_en,
		       product_type_id, format_type_id, vendor_id, unit_type_id,
		       isbn, author_name, publisher_date, edition_number,
		       price, cost, description, note,
		       count_stock, is_active, is_delete, update_by, update_date
		FROM tb_product
		WHERE is_delete = 0
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to fetch products", Data: nil,
		})
	}
	defer rows.Close()

	var list []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.AutoID, &p.Prefix, &p.ProductID, &p.ProductNameTH, &p.ProductNameEN,
			&p.ProductTypeID, &p.FormatTypeID, &p.VendorID, &p.UnitTypeID,
			&p.ISBN, &p.AuthorName, &p.PublisherDate, &p.EditionNumber,
			&p.Price, &p.Cost, &p.Description, &p.Note,
			&p.CountStock, &p.IsActive, &p.IsDelete, &p.UpdateBy, &p.UpdateDate,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status: "error", Message: "Failed to read data", Data: nil,
			})
		}
		list = append(list, p)
	}
	return c.JSON(models.ApiResponse{
		Status: "success", Message: "", Data: list,
	})
}

// 2. Select Paging
func SelectPageProducts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := `
		SELECT autoID, prefix, product_id, product_name_th, product_name_en,
		       product_type_id, format_type_id, vendor_id, unit_type_id,
		       isbn, author_name, publisher_date, edition_number,
		       price, cost, description, note,
		       count_stock, is_active, is_delete, update_by, update_date
		FROM tb_product
		WHERE is_delete = 0
		ORDER BY update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`
	rows, err := config.DB.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to fetch products", Data: nil,
		})
	}
	defer rows.Close()

	var list []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.AutoID, &p.Prefix, &p.ProductID, &p.ProductNameTH, &p.ProductNameEN,
			&p.ProductTypeID, &p.FormatTypeID, &p.VendorID, &p.UnitTypeID,
			&p.ISBN, &p.AuthorName, &p.PublisherDate, &p.EditionNumber,
			&p.Price, &p.Cost, &p.Description, &p.Note,
			&p.CountStock, &p.IsActive, &p.IsDelete, &p.UpdateBy, &p.UpdateDate,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status: "error", Message: "Failed to read data", Data: nil,
			})
		}
		list = append(list, p)
	}

	var total int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM tb_product WHERE is_delete = 0`).Scan(&total)
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
			"products": list,
		},
	})
}

// 3. Select By ID
func SelectProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		SELECT autoID, prefix, product_id, product_name_th, product_name_en,
		       product_type_id, format_type_id, vendor_id, unit_type_id,
		       isbn, author_name, publisher_date, edition_number,
		       price, cost, description, note,
		       count_stock, is_active, is_delete, update_by, update_date
		FROM tb_product
		WHERE product_id = @ID AND is_delete = 0
	`
	row := config.DB.QueryRow(query, sql.Named("ID", id))
	var p models.Product
	if err := row.Scan(
		&p.AutoID, &p.Prefix, &p.ProductID, &p.ProductNameTH, &p.ProductNameEN,
		&p.ProductTypeID, &p.FormatTypeID, &p.VendorID, &p.UnitTypeID,
		&p.ISBN, &p.AuthorName, &p.PublisherDate, &p.EditionNumber,
		&p.Price, &p.Cost, &p.Description, &p.Note,
		&p.CountStock, &p.IsActive, &p.IsDelete, &p.UpdateBy, &p.UpdateDate,
	); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{
				Status: "error", Message: "Product not found", Data: nil,
			})
		}
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to read product", Data: nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status: "success", Message: "", Data: p,
	})
}

// 4. Select By Name (LIKE)
func SelectProductByName(c *fiber.Ctx) error {
	name := c.Params("name")
	query := `
		SELECT autoID, prefix, product_id, product_name_th, product_name_en,
		       product_type_id, format_type_id, vendor_id, unit_type_id,
		       isbn, author_name, publisher_date, edition_number,
		       price, cost, description, note,
		       count_stock, is_active, is_delete, update_by, update_date
		FROM tb_product
		WHERE (product_name_th LIKE '%' + @Name + '%' OR product_name_en LIKE '%' + @Name + '%') AND is_delete = 0
	`
	rows, err := config.DB.Query(query, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to search products", Data: nil,
		})
	}
	defer rows.Close()

	var list []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.AutoID, &p.Prefix, &p.ProductID, &p.ProductNameTH, &p.ProductNameEN,
			&p.ProductTypeID, &p.FormatTypeID, &p.VendorID, &p.UnitTypeID,
			&p.ISBN, &p.AuthorName, &p.PublisherDate, &p.EditionNumber,
			&p.Price, &p.Cost, &p.Description, &p.Note,
			&p.CountStock, &p.IsActive, &p.IsDelete, &p.UpdateBy, &p.UpdateDate,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{
				Status: "error", Message: "Failed to read data", Data: nil,
			})
		}
		list = append(list, p)
	}
	return c.JSON(models.ApiResponse{
		Status: "success", Message: "", Data: list,
	})
}

// 5. Insert
func InsertProduct(c *fiber.Ctx) error {
	var p models.Product
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status: "error", Message: "Invalid request body", Data: nil,
		})
	}

	// 🚩 DUMMY ID for TRIGGER mechanism (product_id is NOT NULL)
	dummyID := "TEMP"

	query := `
		INSERT INTO tb_product (
			product_id, product_name_th, product_name_en,
			product_type_id, format_type_id, vendor_id, unit_type_id,
			isbn, author_name, publisher_date, edition_number,
			price, cost, description, note,
			count_stock, update_by
		)
		VALUES (
			@ProductID, @NameTH, @NameEN,
			@TypeID, @FormatID, @VendorID, @UnitID,
			@ISBN, @Author, @PubDate, @Edition,
			@Price, @Cost, @Desc, @Note,
			@CountStock, @UpdateBy
		)
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("ProductID", dummyID), // Passed strictly to satisfy NOT NULL constraint
				sql.Named("NameTH", p.ProductNameTH),
				sql.Named("NameEN", p.ProductNameEN),
				sql.Named("TypeID", p.ProductTypeID),
				sql.Named("FormatID", p.FormatTypeID),
				sql.Named("VendorID", p.VendorID),
				sql.Named("UnitID", p.UnitTypeID),
				sql.Named("ISBN", p.ISBN),
				sql.Named("Author", p.AuthorName),
				sql.Named("PubDate", p.PublisherDate),
				sql.Named("Edition", p.EditionNumber),
				sql.Named("Price", p.Price),
				sql.Named("Cost", p.Cost),
				sql.Named("Desc", p.Description),
				sql.Named("Note", p.Note),
				sql.Named("CountStock", p.CountStock),
				sql.Named("UpdateBy", p.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to insert product", Data: nil,
		})
	}
	return c.Status(201).JSON(models.ApiResponse{
		Status: "success", Message: "Product added successfully", Data: nil,
	})
}

// 6. Update
func UpdateProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var p models.Product
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(models.ApiResponse{
			Status: "error", Message: "Invalid request body", Data: nil,
		})
	}

	query := `
		UPDATE tb_product
		SET product_name_th = COALESCE(NULLIF(@NameTH, ''), product_name_th),
			product_name_en = @NameEN,
			product_type_id = @TypeID,
			format_type_id = @FormatID,
			vendor_id = @VendorID,
			unit_type_id = @UnitID,
			isbn = @ISBN,
			author_name = @Author,
			publisher_date = @PubDate,
			edition_number = @Edition,
			price = @Price,
			cost = @Cost,
			description = @Desc,
			note = @Note,
			count_stock = @CountStock,
			update_by = @UpdateBy
		WHERE product_id = @ID AND is_delete = 0
	`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query,
				sql.Named("NameTH", p.ProductNameTH),
				sql.Named("NameEN", p.ProductNameEN),
				sql.Named("TypeID", p.ProductTypeID),
				sql.Named("FormatID", p.FormatTypeID),
				sql.Named("VendorID", p.VendorID),
				sql.Named("UnitID", p.UnitTypeID),
				sql.Named("ISBN", p.ISBN),
				sql.Named("Author", p.AuthorName),
				sql.Named("PubDate", p.PublisherDate),
				sql.Named("Edition", p.EditionNumber),
				sql.Named("Price", p.Price),
				sql.Named("Cost", p.Cost),
				sql.Named("Desc", p.Description),
				sql.Named("Note", p.Note),
				sql.Named("CountStock", p.CountStock),
				sql.Named("UpdateBy", p.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to update product", Data: nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status: "success", Message: "Product updated successfully", Data: nil,
	})
}

// 7. Delete (Soft)
func DeleteProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
		UPDATE tb_product
		SET is_delete = 1,
			update_date = CAST(SYSDATETIMEOFFSET() AT TIME ZONE 'SE Asia Standard Time' AS DATETIME)
		WHERE product_id = @ID
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
			Status: "error", Message: "Failed to delete product", Data: nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status: "success", Message: "Product deleted successfully", Data: nil,
	})
}

// 8. Remove (Hard)
func RemoveProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `DELETE FROM tb_product WHERE product_id = @ID`
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(query, sql.Named("ID", id))
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status: "error", Message: "Failed to remove product", Data: nil,
		})
	}
	return c.JSON(models.ApiResponse{
		Status: "success", Message: "Product removed successfully", Data: nil,
	})
}
