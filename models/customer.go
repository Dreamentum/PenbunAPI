package models

type Customer struct {
	CustomerID     string   `json:"customer_id"`
	CustomerTypeID string   `json:"customer_type_id"`
	CustomerName   string   `json:"customer_name"`
	TaxID          *string  `json:"tax_id"`
	BranchName     *string  `json:"branch_name"`
	ContactPerson  *string  `json:"contact_person"`
	Phone1         *string  `json:"phone1"`
	Phone2         *string  `json:"phone2"`
	Email          *string  `json:"email"`
	LineID         *string  `json:"line_id"`
	Address        *string  `json:"address"`
	SubDistrict    *string  `json:"sub_district"`
	District       *string  `json:"district"`
	Province       *string  `json:"province"`
	ZipCode        *string  `json:"zip_code"`
	CreditLimit    *float64 `json:"credit_limit"`
	CreditTermDay  *int     `json:"credit_term_day"`
	Note           *string  `json:"note"`
	UpdateBy       *string  `json:"update_by"`
	UpdateDate     *string  `json:"update_date"`
	IsActive       *bool    `json:"is_active,omitempty"`
	IsDelete       bool     `json:"is_delete"`

	// Optional: Join field
	CustomerTypeName *string `json:"customer_type_name,omitempty"`
}
