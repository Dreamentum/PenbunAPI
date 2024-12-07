package models

type Publisher struct {
	PublisherID   string `json:"publisher_ID"`
	PublisherName string `json:"publisher_name"`
	Note          string `json:"note"`
	UpdateDate    string `json:"update_date"`
}
