package models

type ProductTypeGroup struct {
	GroupCode   string  `json:"group_code"`
	GroupName   string  `json:"group_name"`
	Description *string `json:"description,omitempty"`
	UpdateBy    *string `json:"update_by,omitempty"`
	UpdateDate  *string `json:"update_date,omitempty"`
	IDStatus    bool    `json:"id_status"`
	IsDelete    bool    `json:"is_delete"`
}
