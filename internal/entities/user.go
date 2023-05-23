package entities

type User struct{
	Model
	PhoneNumber		string		`json:"phonenumber" validate:"required"`
	UserName		string		`json:"username" validate:"required"`
	Merchants		[]*Merchant		`gorm:"many2many:merchant_users;"`
	// Balance 		float64 		`json:"balance"`
	// Transaction 	Transaction			`gorm:"foreignkey:transaction.ID"`
}


type UserLogin struct {

	PhoneNumber		string 	    `json:"phonenumber"`
	
}


type Transaction struct{
	Model
	Amount 		float64 	`json:"amount"`
	Type 		string		`json:"type"`
	Action	 	string		`json:"action"`

}

// type MerchantUser struct{
// 	MerchantID		string	
// 	UserID			string
// }