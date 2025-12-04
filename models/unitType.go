package models

type UnitType struct {
	UnitTypeID   string `json:"unit_type_id"`
	UnitTypeName string `json:"unit_type_name"`
	Description  string `json:"description"`
	UpdateBy     string `json:"update_by"`
	UpdateDate   string `json:"update_date"`
	IDStatus     bool   `json:"id_status"`
}
