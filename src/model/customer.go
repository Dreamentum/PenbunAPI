package model

import "time"

type CUSTOMER struct {
	ID        *string   `json:"id"`
	Name      *string   `json:"name"`
	Mobile    *string   `json:"mobile"`
	Address   *string   `json:"address"`
	Discount  *float64  `json:"discount"`
	CreatedAt time.Time `json:"createDate" sql:"createAt"`
	UpdatedAt time.Time `json:"updateDate" sql:"updateAt"`
}

// BOOKS array of BOOK type
type CUSTOMERS []CUSTOMER
