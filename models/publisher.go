package models

type Publisher struct {
	AutoID      int     `json:"auto_id"`
	PublisherID string  `json:"publisher_id"`
	PublisherName string `json:"publisher_name"`
	Address     *string `json:"address,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	IsActive bool `json:"is_active"`
	UpdateBy    string  `json:"update_by"`
	UpdateDate  string  `json:"update_date"`
	IsDelete    bool    `json:"is_delete"`
}
