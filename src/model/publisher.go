package model

import "time"

type PUBLISHER struct {
	ID           *string   `json:"id"`
	Name         *string   `json:"name"`
	ContactName  *string   `json:"contactname"`
	ContactPhone *string   `json:"contactphone"`
	Address      *string   `json:"address"`
	District     *string   `json:"district"`
	ZipCode      *string   `json:"zipcode"`
	Description  *string   `json:"description"`
	Province     *string   `json:"province"`
	Status       bool      `json:"status"`
	CreatedAt    time.Time `json:"createDate" sql:"createAt"`
	UpdatedAt    time.Time `json:"updateDate" sql:"updateAt"`
	CreateBy     *string   `json:"createBy"`
	UpdateBy     *string   `json:"updateBy"`
	IsDelete     bool      `json:"isdelete"`
}

type PUBLISHERS []PUBLISHER
// Group array of PUBLISHER type
// type PUBLISHER []PUBLISHERS
