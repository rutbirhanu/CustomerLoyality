package entities

type User struct{
	Model
	PhoneNumber		string		`json:"phonenumber" validate:"required, max=30, min=6"`
	UserName		string		`json:"username" validate:"required, max=30, min=6"`
	// Amount 			string		`json:"amount"`
	Merchants		[]*Merchant		`gorm:"many2many:merchant_users;"`
}


// type MerchantUser struct{
// 	MerchantID		string	
// 	UserID			string
// }