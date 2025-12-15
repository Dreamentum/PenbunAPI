package models

type ProductGroup struct {
	ProductGroupID    string  `json:"product_group_id"`
	ProductCategoryID string  `json:"product_category_id"`
	ProductGroupName  string  `json:"product_group_name"`
	Description       *string `json:"description"`
	UpdateBy          *string `json:"update_by"`
	UpdateDate        *string `json:"update_date"`
	IsActive          *bool   `json:"is_active,omitempty"`
	IsDelete          bool    `json:"is_delete"`

	// Optional: Join field
	CategoryName *string `json:"category_name,omitempty"`
}
