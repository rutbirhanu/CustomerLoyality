package entities

import "gorm.io/gorm"

type Merchant struct {
	Model
	MerchantName string ` validate:" required , max=30"  json:"name" `
	Password     string ` validate:" required , max=30, min=6 "  json:"password" `
	PhoneNumber  string ` gorm:"primary_key" validate:" required "  json:"phonenumber" `
	Token        string `json:"token"`
	BusinessName string ` validate:" required , max=30"  json:"businessname" `
	PrivateKey   string
	PublicKey    string
	Users        []*User `gorm:"many2many:merchant_users;"`
}

type MerchantLogin struct {
	Password    string ` validate:" required , max=30 , min=6" json:"password"`
	PhoneNumber string ` validate:" required , max=30 " json:"phonenumber"`
	Token       string `json:"token"`

	// 	RefreshToken	string 		`json:"refreshToken"`
}

type GetAllUsers struct {
	Users    []User
	Total    int64
	Page     int64
	Prev     int64
	Next     int64
	LastPage int64
}

type MerchantUsers struct {
	gorm.Model
	MerchantID string ` json:"merchantid"`
	MerchantPhone 	string 	`gorm:"foreignKey"  json:"merchantphone"`
	// Merchant   Merchant `json:"merchant"`
	UserPhone  string		`gorm:"foreignKey"  json:"userphone"`
	UserID string `gorm:"foreignKey"  json:"userid"`
	// User       User     `json:"user"`
}
