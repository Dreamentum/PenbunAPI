package models

type Customer struct {
	CustomerCode   string  `json:"customer_code"`
	CustomerName   string  `json:"customer_name"`
	BizID          string  `json:"biz_id"`
	CustomerTypeID string  `json:"customer_type_id"`
	ContractID     *string `json:"contract_id"`
	DiscountID     *string `json:"discount_id"`
	FirstName      string  `json:"first_name"`
	LastName       *string `json:"last_name"`
	Email          *string `json:"email"`
	Phone1         *string `json:"phone1"`
	Phone2         *string `json:"phone2"`
	TaxID          *string `json:"tax_id"`
	Address        *string `json:"address"`
	District       *string `json:"district"`
	Province       *string `json:"province"`
	ZipCode        *string `json:"zip_code"`
	Note           *string `json:"note"`
	Reference1     *string `json:"reference1"`
	Reference2     *string `json:"reference2"`
	UpdateBy       string  `json:"update_by"`
	UpdateDate     string  `json:"update_date"`
	IDStatus       bool    `json:"id_status"`
	IsDelete       bool    `json:"is_delete"`
}
