package models

type ProductFormatType struct {
	AutoID             int     `json:"auto_id"`
	ProductFormatTypeID string `json:"product_format_type_id"`
	TypeName           string  `json:"type_name"`
	Description        *string `json:"description,omitempty"`
	IsActive bool `json:"is_active"`
	UpdateBy           string  `json:"update_by"`
	UpdateDate         string  `json:"update_date"`
	IsDelete           bool    `json:"is_delete"`
}
