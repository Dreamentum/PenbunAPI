package models

type ProductType struct {
	ProductTypeID *string `json:"product_type_id"`
	TypeName      *string `json:"type_name"`
	TypeGroupName *string `json:"type_group_name"`
	Description   *string `json:"description"`
	UpdateBy      *string `json:"update_by"`
	UpdateDate    *string `json:"update_date"`
	IDStatus      *bool   `json:"id_status"`
	IsDelete      *bool   `json:"is_delete"`
}
