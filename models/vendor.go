package models

type Vendor struct {
	VendorID      string  `json:"vendor_id"`
	VendorTypeID  string  `json:"vendor_type_id"`
	VendorName    string  `json:"vendor_name"`
	TaxID         *string `json:"tax_id"`
	BranchName    *string `json:"branch_name"`
	ContactPerson *string `json:"contact_person"`
	Phone1        *string `json:"phone1"`
	Phone2        *string `json:"phone2"`
	Email         *string `json:"email"`
	Website       *string `json:"website"`
	Address       *string `json:"address"`
	SubDistrict   *string `json:"sub_district"`
	District      *string `json:"district"`
	Province      *string `json:"province"`
	ZipCode       *string `json:"zip_code"`
	CreditTermDay *int    `json:"credit_term_day"`
	Currency      *string `json:"currency"`
	Note          *string `json:"note"`
	UpdateBy      *string `json:"update_by"`
	UpdateDate    *string `json:"update_date"`
	IsActive      *bool   `json:"is_active,omitempty"`
	IsDelete      bool    `json:"is_delete"`

	// Optional: Join field
	VendorTypeName *string `json:"vendor_type_name,omitempty"`
}
