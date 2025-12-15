package models

type ProductCategory struct {
	ProductCategoryID string  `json:"product_category_id"`
	CategoryName      string  `json:"category_name"`
	CategoryCode      string  `json:"category_code"`
	Description       *string `json:"description,omitempty"`
	UpdateBy          *string `json:"update_by,omitempty"`
	UpdateDate        *string `json:"update_date,omitempty"`
	IsActive          *bool   `json:"is_active,omitempty"`
	IsDelete          bool    `json:"is_delete"`
}
