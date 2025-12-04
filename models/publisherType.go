package models

type PublisherType struct {
	PublisherTypeID string  `json:"publisher_type_id"`
	TypeName        string  `json:"type_name"`
	Description     *string `json:"description"`
	UpdateBy        *string `json:"update_by"`
	UpdateDate      *string `json:"update_date"`
	IDStatus        bool    `json:"id_status"`
	IsDelete        bool    `json:"is_delete"`
}
