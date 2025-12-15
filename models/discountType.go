package models

type DiscountType struct {
	DiscountTypeID   string  `json:"discount_type_id"`
	DiscountTypeName string  `json:"discount_type_name"`
	Description      *string `json:"description"`
	UpdateBy         *string `json:"update_by"`
	UpdateDate       *string `json:"update_date"`
	IsActive         *bool   `json:"is_active,omitempty"`
	IsDelete         bool    `json:"is_delete"`
}
