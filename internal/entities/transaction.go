package entities

type Transaction struct {
	Model
	Amount         float64 `json:"amount"`
	Type           string  `json:"type"`
	Action         string  `json:"action"`
	To             string  `json:"to"`
	UserMerchantID string  `gorm:"foreignkey"`
}

type Collection struct {
	Model
	Points         float64 `json:"point"`
	Reason         string  `json:"reason"`
	UserMerchantId string  `gorm:"foreignkey"`
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
