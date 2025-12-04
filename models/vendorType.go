package models

type VendorType struct {
	VendorTypeID string  `json:"vendor_type_id"`
	Prefix       string  `json:"prefix"`
	TypeName     string  `json:"type_name"`
	Description  *string `json:"description"`
	UpdateBy     *string `json:"update_by"`
	UpdateDate   *string `json:"update_date"`
	IDStatus     bool    `json:"id_status"`
	IsDelete     bool    `json:"is_delete"`
}
