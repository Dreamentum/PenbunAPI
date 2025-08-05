package models

type ProductFormatType struct {
	ProductFormatTypeID string  `json:"product_format_type_id"`
	FormatName          string  `json:"format_name"`
	Description         *string `json:"description,omitempty"`
	UpdateBy            *string `json:"update_by,omitempty"`
	UpdateDate          *string `json:"update_date,omitempty"`
	IDStatus            bool    `json:"id_status"`
	IsDelete            bool    `json:"is_delete"`
}
