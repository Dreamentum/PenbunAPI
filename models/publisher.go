package models

type Publisher struct {
	PublisherCode     string  `json:"publisher_code"`
	PublisherTypeID   string  `json:"publisher_type_id"`
	PublisherTypeName *string `json:"publisher_type_name"` // 👈 New field for JOIN
	PublisherName     string  `json:"publisher_name"`
	ContactName1      *string `json:"contact_name1"`
	ContactName2      *string `json:"contact_name2"`
	Email             *string `json:"email"`
	Phone1            *string `json:"phone1"`
	Phone2            *string `json:"phone2"`
	Address           *string `json:"address"`
	District          *string `json:"district"`
	Province          *string `json:"province"`
	ZipCode           *string `json:"zip_code"`
	Note              *string `json:"note"`
	DiscountID        *string `json:"discount_id"`
	UpdateBy          *string `json:"update_by"`
	UpdateDate        *string `json:"update_date"`
	IDStatus          bool    `json:"id_status"`
	IsDelete          bool    `json:"is_delete"`
}
