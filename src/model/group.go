package model

import "time"

type GROUP struct {
	ID          *string   `json:"id"`
	Name        *string   `json:"name"`
	Level       *string   `json:"level"`
	Description *string   `json:"description"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"createDate" sql:"createAt"`
	UpdatedAt   time.Time `json:"updateDate" sql:"updateAt"`
	CreateBy    *string   `json:"createBy"`
	UpdateBy    *string   `json:"updateBy"`
	IsDelete    bool      `json:"isdelete"`
}

// Group array of GROUP type
type GROUPS []GROUP
