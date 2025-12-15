package models

import "time"

type Product struct {
	AutoID         int       `json:"auto_id"`
	Prefix         string    `json:"prefix"`
	ProductID      string    `json:"product_id"`
	ProductNameTH  string    `json:"product_name_th"`
	ProductNameEN  *string   `json:"product_name_en"`
	ProductTypeID  string    `json:"product_type_id"` // This is technically product_group_id in v2.2 SQL. I should rename it to match SQL `product_group_id`.
	FormatTypeID   *string   `json:"format_type_id"` // SQL: product_format_type_id. Go: format_type_id. Close enough but inconsistency.
	VendorID       *string   `json:"vendor_id"`
	UnitTypeID     *string   `json:"unit_type_id"`
	ISBN           *string   `json:"isbn"`
	AuthorName     *string   `json:"author_name"`
	PublisherDate  *string   `json:"publisher_date"`
	EditionNumber  *int      `json:"edition_number"`
	Price          float64   `json:"price"` // SQL: sell_price. Go value: price.
	Cost           float64   `json:"cost"`  // SQL: cost_price. Go value: cost.
	Description    *string   `json:"description"`
	Note           *string   `json:"note"`
	CountStock     bool      `json:"count_stock"` // SQL: count_stock (1=stock, 0=service).
	IsActive       bool      `json:"is_active"`
	IsDelete       bool      `json:"is_delete"`
	UpdateBy       *string   `json:"update_by"`
	UpdateDate     *time.Time `json:"update_date"`
}
