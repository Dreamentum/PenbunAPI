package models

import "time"

type Order struct {
	AutoID         int        `json:"auto_id"`
	Prefix         string     `json:"prefix"`
	OrderID        string     `json:"order_id"`
	CustomerID     string     `json:"customer_id"`
	WarehouseID    string     `json:"warehouse_id"`
	DocDate        time.Time  `json:"doc_date"`
	DocType        string     `json:"doc_type"` // CASH, CREDIT
	TotalAmount    float64    `json:"total_amount"`
	DiscountAmount float64    `json:"discount_amount"`
	NetAmount      float64    `json:"net_amount"`
	VatAmount      float64    `json:"vat_amount"`
	GrandTotal     float64    `json:"grand_total"`
	UpdateBy       string     `json:"update_by"`
	UpdateDate     time.Time  `json:"update_date"`
	IsActive       bool       `json:"is_active"`
	IsDelete       bool       `json:"is_delete"`

	// Join Fields
	CustomerName   *string    `json:"customer_name,omitempty"`
}

type OrderItem struct {
	AutoID         int        `json:"auto_id"`
	OrderID        string     `json:"order_id"`
	ProductID      string     `json:"product_id"`
	Qty            float64    `json:"qty"`
	UnitPrice      float64    `json:"unit_price"`
	DiscountAmount float64    `json:"discount_amount"`
	LineTotal      float64    `json:"line_total"`
	Remark         *string    `json:"remark"`
	UpdateDate     *time.Time `json:"update_date"`
	IsDelete       bool       `json:"is_delete"`

	// Join Fields
	ProductName    *string    `json:"product_name,omitempty"`
}
