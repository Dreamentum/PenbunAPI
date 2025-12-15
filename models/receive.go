package models

import "time"

type ReceiveNote struct {
	AutoID        int        `json:"auto_id"`
	Prefix        string     `json:"prefix"`
	ReceiveNoteID string     `json:"receive_note_id"`
	VendorID      string     `json:"vendor_id"`
	WarehouseID   string     `json:"warehouse_id"`
	DocDate       time.Time  `json:"doc_date"`
	RefInvoiceNo  *string    `json:"ref_invoice_no"`
	ReceiveType   string     `json:"receive_type"` // PO, GIFT, RETURN
	TotalAmount   float64    `json:"total_amount"`
	Note          *string    `json:"note"`
	UpdateBy      string     `json:"update_by"`
	UpdateDate    time.Time  `json:"update_date"`
	IsActive      bool       `json:"is_active"`
	IsDelete      bool       `json:"is_delete"`
	
	// Join Fields (Optional)
	VendorName    *string    `json:"vendor_name,omitempty"`
}

type ReceiveItem struct {
	AutoID        int        `json:"auto_id"`
	ReceiveNoteID string     `json:"receive_note_id"`
	ProductID     string     `json:"product_id"`
	Qty           float64    `json:"qty"`
	UnitCost      float64    `json:"unit_cost"`
	LineTotal     float64    `json:"line_total"`
	Remark        *string    `json:"remark"`
	UpdateDate    *time.Time `json:"update_date"`
	IsDelete      bool       `json:"is_delete"`

	// Join Fields
	ProductName   *string    `json:"product_name,omitempty"`
}
