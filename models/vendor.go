package models

type Vendor struct {
	VendorID     *string `json:"vendor_code"`
	VendorTypeID string  `json:"vendor_type_id"`
	TypeName     *string `json:"type_name"` // 👈 เพิ่มสำหรับ JOIN กับ tb_vendor_type
	DiscountID   *string `json:"discount_id"`
	DiscountName *string `json:"discount_name"`
	VendorName   string  `json:"vendor_name"`
	ContactName1 *string `json:"contact_name1"`
	ContactName2 *string `json:"contact_name2"`
	Email        *string `json:"email"`
	Phone1       *string `json:"phone1"`
	Phone2       *string `json:"phone2"`
	Address      *string `json:"address"`
	District     *string `json:"district"`
	Province     *string `json:"province"`
	ZipCode      *string `json:"zip_code"`
	Note         *string `json:"note"`
	UpdateBy     *string `json:"update_by"`
	UpdateDate   *string `json:"update_date"`
	IDStatus     bool    `json:"id_status"`
}
