package models

type ProductType struct {
	ProductTypeID string  `json:"product_type_id"`
	TypeName      string  `json:"type_name"`
	TypeGroupCode string  `json:"type_group_code"`
	TypeGroupName *string `json:"type_group_name,omitempty"` // JOIN มาจาก group table
	Description   *string `json:"description,omitempty"`
	UpdateBy      *string `json:"update_by,omitempty"`
	UpdateDate    *string `json:"update_date,omitempty"`
	IDStatus      *bool   `json:"id_status,omitempty"`
	IsDelete      bool    `json:"is_delete"`
}
