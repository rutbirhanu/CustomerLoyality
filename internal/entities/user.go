package entities

type User struct{
	Model
	// Password		string		`json:"password" validate:"required, max=30, min=6"`
	PhoneNumber		string		`json:"phonenumber" validate:"required, max=30, min=6"`
	UserName		string		`json:"username" validate:"required, max=30, min=6"`
	// Merchants		[]*Merchant		`gorm:"many2many:merchant_users;"`
}


type MerchantUser struct{
	MerchantID		string	
	UserID			string
}