package models

type Discount struct {
	DiscountID     string   `json:"discount_id"`
	DiscountTypeID string   `json:"discount_type_id"`
	DiscountName   string   `json:"discount_name"`
	DiscountCode   *string  `json:"discount_code"`
	Description    *string  `json:"description"`
	DiscountValue  float64  `json:"discount_value"`
	IsPercent      bool     `json:"is_percent"`
	MinOrderAmount *float64 `json:"min_order_amount"`
	StartDate      *string  `json:"start_date"`
	EndDate        *string  `json:"end_date"`
	UpdateBy       *string  `json:"update_by"`
	UpdateDate     *string  `json:"update_date"`
	IsActive       *bool    `json:"is_active,omitempty"`
	IsDelete       bool     `json:"is_delete"`

	// Optional: Join field
	DiscountTypeName *string `json:"discount_type_name,omitempty"`
}
