package entities

type User struct {
	Model
	PhoneNumber string      ` gorm:"primary_key" json:"phonenumber" validate:"required"`
	UserName    string      `json:"username" validate:"required"`
	Merchants   []*Merchant `gorm:"many2many:merchant_users;"`
	Balance     float64     `json:"balance"`
	// Transaction 	[]*Transaction
}

type UserLogin struct {
	PhoneNumber string `json:"phonenumber"`
}


// type MerchantUser struct{
// 	MerchantID		string
// 	UserID			string
// }