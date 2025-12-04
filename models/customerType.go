package models

type CustomerType struct {
	CustomerTypeID string  `json:"customer_type_id"`
	TypeName       string  `json:"type_name"`
	Description    *string `json:"description"`
	UpdateBy       string  `json:"update_by"`
	UpdateDate     string  `json:"update_date"`
	IDStatus       bool    `json:"id_status"`
	IsDelete       bool    `json:"is_delete"`
}
