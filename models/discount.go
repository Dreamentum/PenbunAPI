package models

import "time"

type Discount struct {
	DiscountID    string     `json:"discount_id"`
	DiscountName  string     `json:"discount_name"`
	DiscountType  string     `json:"discount_type"`
	DiscountValue float64    `json:"discount_value"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	Note          *string    `json:"note"` // <-- nullable
	UpdateBy      string     `json:"update_by"`
	UpdateDate    *time.Time `json:"update_date"` // <-- nullable
	IDStatus      bool       `json:"id_status"`
}
