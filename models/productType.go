package models

import "time"

type ProductType struct {
	ProductTypeID string     `json:"product_type_id"`
	TypeName      string     `json:"type_name"`
	TypeGroupCode string     `json:"type_group_code"`
	TypeGroupName *string    `json:"type_group_name,omitempty"` // JOIN จาก group table
	Description   *string    `json:"description,omitempty"`
	UpdateBy      *string    `json:"update_by,omitempty"`
	UpdateDate    *time.Time `json:"update_date,omitempty"` // DATETIME → *time.Time
	IDStatus      bool       `json:"id_status"`             // BIT → bool
	IsDelete      bool       `json:"is_delete"`
}
