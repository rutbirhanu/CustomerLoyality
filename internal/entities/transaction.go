package entities

type Transaction struct {
	Model
	Amount         float64 `json:"amount"`
	Type           string  `json:"type"`
	ReceiverID      string  `json:"receiver,omitempty"`
	UserMerchantID string  `gorm:"foreignkey"`
}

type TransferPoint struct{
	Amount 		float64		`json:"amount"`
	Phone		string 		`json:"phone"`
}

type Collection struct {
	Model
	Points         float64 		`json:"point"`
	UserPhone 		string 		`json:"phone"`
	// UserMerchantId string  		`gorm:"foreignkey"`
}

type RequestData struct{
	Amount 		float64 	`json:"amount"`
}

var Organizations = []map[string]string{
	{
		"id":      "001",
		"name":    "mekedonya",
		"account": "989987654",
	},
	{
		"id":      "002",
		"name":    "cancer association",
		"account": "6784784789",
	},
	{
		"id":      "003",
		"name":    "kidney failure charity",
		"account": "2736868648",
	},
}
