package entities


type Transaction struct {
	Model
	Amount         float64        `json:"amount"`
	Type           string         `json:"type"`
	Action         string         `json:"action"`
	To             string         `json:"to"`
	UserMerchantID string         `gorm:"foreignkey"`
}


type Reward struct{
	Model
	Points 				float64 		`json:"amount"`
	Reason 				string 			`json:"reason"`
	UserMerchantId 		string			`gorm:"foreignkey"`
}

