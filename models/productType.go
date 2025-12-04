package models

import "time"

type ProductType struct {
	ProductTypeID      string     `json:"product_type_id"`
	TypeName           string     `json:"type_name"`
	ProductTypeGroupID string     `json:"product_type_group_id"` // FK (NVARCHAR) → tb_product_type_group.group_id
	GroupId            string     `json:"group_id"`              // Redundant but kept for compatibility (same as ProductTypeGroupID)
	GroupName          string     `json:"group_name"`            // JOIN from tb_product_type_group.group_name
	Description        string     `json:"description"`
	UpdateBy           string     `json:"update_by"`
	UpdateDate         *time.Time `json:"update_date,omitempty"`
	IDStatus           bool       `json:"id_status"` // BIT
	IsDelete           bool       `json:"is_delete"` // BIT

	// Optional input (alternative to ProductTypeGroupID)
	TypeNameGroupCode string `json:"type_name_group_code,omitempty"` // e.g., "PTG000123"
}
