package models

// Reference represents the structure of the tb_reference table
type Reference struct {
	RowID      int    `json:"row_id"`
	RefID      string `json:"ref_id"`
	RefInt     int    `json:"ref_int"`
	RefText    string `json:"ref_text"`
	UpdateBy   string `json:"update_by"`
	UpdateDate string `json:"update_date"`
}