package entities


type Merchant struct {
	Model
	MerchantName    string  	` validate:" required , max=30"  json:"name" `
	Password     	string  	` validate:" required , max=30, min=6 "  json:"password" `
	PhoneNumber  	string  	` validate:" required "  json:"phonenumber" `
	Token 			string       `json:"token"`
	BusinessName 	string  	` validate:" required , max=30"  json:"businessname" `
	// Users			[]*User			`gorm:"many2many:user_roles;"`
}


type MerchantLogin struct {
	Password 		string		` validate:" required , max=30 , min=6" json:"password"`
	PhoneNumber		string		` validate:" required , max=30 " json:"phonenumber"`
	Token 			string       `json:"token"`

// 	RefreshToken	string 		`json:"refreshToken"`
}

