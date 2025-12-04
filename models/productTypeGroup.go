package models

import "time"

type ProductTypeGroup struct {
	GroupId     string     `json:"group_id"`
	GroupName   string     `json:"group_name"`
	Description *string    `json:"description,omitempty"`
	UpdateBy    *string    `json:"update_by,omitempty"`
	UpdateDate  *time.Time `json:"update_date,omitempty"` // DATETIME → *time.Time
	IDStatus    bool       `json:"id_status"`             // BIT → bool
	IsDelete    bool       `json:"is_delete"`
}
