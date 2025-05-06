package models

type DiscountType struct {
	DiscountTypeID   string  `json:"discount_type_id"`
	TypeName         string  `json:"type_name"`
	DiscountUnitType string  `json:"discount_unit_type"`
	Description      *string `json:"description"`
	UpdateBy         *string `json:"update_by"`
	UpdateDate       *string `json:"update_date"`
	IDStatus         bool    `json:"id_status"`
}
