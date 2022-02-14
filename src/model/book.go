package model

import "time"

type BOOK struct {
	ID          *string  `json:"id"`
	Isbn        *string  `json:"isbn"`
	Name        *string  `json:"name"`
	PublisherID *string  `json:"publisherId"`
	UpdatedDate *string  `json:"updatedDate"`
	Price       *float64 `json:"price"`
}

// BOOKS array of BOOK type
type BOOKS []BOOK

type BOOK_CATALOG struct {
	ID          *string   `json:"id"`
	Name        *string   `json:"name"`
	OwenerName  *string   `json:"owenername"`
	Description *string   `json:"description"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"createDate" sql:"createAt"`
	UpdatedAt   time.Time `json:"updateDate" sql:"updateAt"`
	CreateBy    *string   `json:"createBy"`
	UpdateBy    *string   `json:"updateBy"`
	IsDelete    bool      `json:"isdelete"`
}

// BOOK_CATALOGS array of BOOK_CATALOG (catalog)
type BOOK_CATALOGS []BOOK_CATALOG

type BOOK_GROUP struct {
	ID           *string  `json:"id"`
	Name         *string  `json:"name"`
	MasterID     *string  `json:"masterid"`
	MasterName   *string  `json:"mastername"`
	MasterVolumn *float64 `json:"mastervolumn"`
}

// BOOK_GROUPS array of BOOK_GROUP type
type BOOK_GROUPS []BOOK_GROUP
