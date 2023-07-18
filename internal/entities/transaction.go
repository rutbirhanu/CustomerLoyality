package entities


type Transaction struct {
	Model
	Amount         float64 `json:"amount"`
	Type           string  `json:"type"`
	Action         string  `json:"action"`
	ReceiverID     string  `json:"receiver,omitempty"`
	UserMerchantID string  `gorm:"foreignkey"`
}

type TransferPoint struct {
	Amount float64 `json:"amount"`
	Phone  string  `json:"phone"`
}

type Collection struct {
	Points    float64 `json:"point"`
	UserPhone string  `json:"phone"`
	// UserMerchantId string  		`gorm:"foreignkey"`
}

type RequestData struct {
	Amount float64 `json:"amount"`
}

type GiftCard struct {
	Model
	Amount        float64 `json:"amount"`
	RecipentPhone string  `json:"recipientPhone"`
	MerchantId    string  `json:"merchantId"`
	GiftCardCode  string  `json:"giftCardCode"`
	PurchaserPhone		string 		`json:"purchasePhone"`
	Active			bool		`json:"active,omitempty"`
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
