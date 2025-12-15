package models

type Warehouse struct {
	WarehouseID        string  `json:"warehouse_id"`
	WarehouseCode      string  `json:"warehouse_code"`
	WarehouseName      string  `json:"warehouse_name"`
	Description        *string `json:"description"`
	IsMainDC           *bool   `json:"is_main_dc"`
	AllowNegativeStock *bool   `json:"allow_negative_stock"`
	UpdateBy           *string `json:"update_by"`
	UpdateDate         *string `json:"update_date"`
	IsActive           *bool   `json:"is_active,omitempty"`
	IsDelete           bool    `json:"is_delete"`
}
