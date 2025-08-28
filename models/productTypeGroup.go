package models

import "time"

type ProductTypeGroup struct {
	GroupCode   string     `json:"group_code"`
	Prefix      string     `json:"prefix,omitempty"`
	GroupName   string     `json:"group_name"`
	Description *string    `json:"description,omitempty"`
	UpdateBy    *string    `json:"update_by,omitempty"`
	UpdateDate  *time.Time `json:"update_date,omitempty"`
	IDStatus    bool       `json:"id_status"` // ✅ BIT → bool
	IsDelete    bool       `json:"is_delete"`
}
