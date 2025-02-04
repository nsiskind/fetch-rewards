package models

type Item struct {
	ShortDescription string `json:"shortDescription" validate:"description"`
	Price            string `json:"price" validate:"price"`
}

type Receipt struct {
	Retailer     string  `json:"retailer" validate:"retailer"`
	PurchaseDate string  `json:"purchaseDate" validate:"required,datetime=2006-01-02"`
	PurchaseTime string  `json:"purchaseTime" validate:"required,datetime=15:04"`
	Items        []*Item `json:"items" validate:"required,min=1"`
	Total        string  `json:"total" validate:"total"`
}

type RecieptResponse struct {
	Id string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}
