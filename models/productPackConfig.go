package models

type ProductPackConfig struct {
	AutoID              int    `json:"auto_id"`
	ProductPackConfigID string `json:"product_pack_config_id"`
	ProductID           string `json:"product_id"`
	BundleQty           int    `json:"bundle_qty"`
	UnitTypeID          string `json:"unit_type_id"`
	Note                string `json:"note"`
	UpdateBy            string `json:"update_by"`
	UpdateDate          string `json:"update_date"`
	IDStatus            bool   `json:"id_status"`
	IsDelete            bool   `json:"is_delete"`
}
