package structs

type TendersGuruResponse struct {
	Data []TendersGuruData `json:"data"`
}

type TendersGuruData struct {
	ID              string    `json:"id"`
	Date            string    `json:"date"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Place           string    `json:"place"`
	AwardedValue    string    `json:"awarded_value"`
	AwardedCurrency string    `json:"awarded_currency"`
	Purchaser       Purchaser `json:"purchaser"`
	Awarded         []Awarded `json:"awarded"`
}

type Purchaser struct {
	Name string `json:"name"`
}

type Awarded struct {
	Suppliers []Suppliers `json:"suppliers"`
}

type Suppliers struct {
	Name string `json:"name"`
}
