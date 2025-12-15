package models

type CustomerType struct {
	CustomerTypeID string  `json:"customer_type_id"`
	CustomerTypeName string  `json:"customer_type_name"`
	BaseCreditDay    *int    `json:"base_credit_day"`
	Description      *string `json:"description"`
	UpdateBy       *string `json:"update_by"`
	UpdateDate     *string `json:"update_date"`
	IsActive       *bool   `json:"is_active,omitempty"`
	IsDelete       bool    `json:"is_delete"`
}
