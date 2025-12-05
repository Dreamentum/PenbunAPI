package models

import "time"

type Product struct {
	AutoID         int       `json:"auto_id"`
	Prefix         string    `json:"prefix"`
	ProductID      string    `json:"product_id"`
	ProductNameTH  string    `json:"product_name_th"`
	ProductNameEN  *string   `json:"product_name_en"`
	ProductTypeID  string    `json:"product_type_id"`
	FormatTypeID   *string   `json:"format_type_id"`
	VendorID       *string   `json:"vendor_id"`
	UnitTypeID     *string   `json:"unit_type_id"`
	ISBN           *string   `json:"isbn"`
	AuthorName     *string   `json:"author_name"`
	PublisherDate  *string   `json:"publisher_date"` // Keeping as string to match existing patterns (or time.Time if prefered, but book uses *string usually? Let's check book.go)
	EditionNumber  *int      `json:"edition_number"`
	Price          float64   `json:"price"`
	Cost           float64   `json:"cost"`
	Description    *string   `json:"description"`
	Note           *string   `json:"note"`
	IsService      bool      `json:"is_service"`
	IDStatus       bool      `json:"id_status"`
	IsDelete       bool      `json:"is_delete"`
	UpdateBy       *string   `json:"update_by"`
	UpdateDate     *time.Time `json:"update_date"` // book.go uses string? Let's re-verify book.go
}
