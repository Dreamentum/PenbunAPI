package models

import "time"

type ProductTypeGroup struct {
	GroupCode   string     `json:"group_code"`       // from trigger
	Prefix      string     `json:"prefix,omitempty"` // ใช้ตอน insert เพื่อให้ trigger สร้างรหัส
	GroupName   string     `json:"group_name"`
	Description *string    `json:"description,omitempty"`
	UpdateBy    *string    `json:"update_by,omitempty"`
	UpdateDate  *time.Time `json:"update_date,omitempty"` // DATETIME → time.Time
	IDStatus    string     `json:"id_status"`             // NVARCHAR(40) → string
	IsDelete    bool       `json:"is_delete"`
}
