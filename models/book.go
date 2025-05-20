package models

type Book struct {
	BookCode         string   `json:"book_code"`
	BookName         string   `json:"book_name"`
	BookISBN         *string  `json:"book_isbn"`
	BookBarcode      *string  `json:"book_barcode"`
	BookTypeID       string   `json:"book_type_id"`
	BookFormatTypeID *string  `json:"book_format_type_id"`
	BookExternalCode *string  `json:"book_external_code"`
	PublisherCode    string   `json:"publisher_code"`
	BookPrice        *float64 `json:"book_price"`
	BookDiscount     *float64 `json:"book_discount"`
	Description      *string  `json:"description"`
	Note             *string  `json:"note"`
	UpdateBy         string   `json:"update_by"`
	UpdateDate       string   `json:"update_date"`
	IDStatus         bool     `json:"id_status"`
	IsDelete         bool     `json:"is_delete"`
}
