package models

import "time"

type ProductType struct {
	ProductTypeID     string     `json:"product_type_id"`
	TypeName          string     `json:"type_name"`            // เก็บรหัสกลุ่ม (อ้างอิง tb_product_type_group.group_code) จากคอลัมน์ p.type_group_name
	TypeNameGroupCode string     `json:"type_name_group_code"` // ชื่อกลุ่มจากการ JOIN กับ tb_product_type_group.group_name
	GroupName         string     `json:"group_name"`
	Description       string     `json:"description"`
	UpdateBy          string     `json:"update_by"`
	UpdateDate        *time.Time `json:"update_date,omitempty"`
	IDStatus          bool       `json:"id_status"` // BIT
	IsDelete          bool       `json:"is_delete"` // BIT
}
