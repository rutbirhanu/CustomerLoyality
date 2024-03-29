package entities


import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Model
	PhoneNumber string      `  json:"phonenumber" validate:"required"`
	UserName    string      `json:"username" validate:"required"`
	Merchants   []*Merchant `gorm:"many2many:wallets;"`
	Token		string 			`json:"token"`
	PrivateKey   string
	PublicKey    string
}

type UserInfo struct{
	PhoneNumber string      `  json:"phonenumber" validate:"required"`
	UserName    string      `json:"username" validate:"required"`
}

type UserLogin struct {
	PhoneNumber string 		`json:"phonenumber"`
	UserName 	string 		`json:"username"`
	Token		string 			`json:"token"`
}

type CreatedUserResponse struct{
	PhoneNumber	 string     	 `json:"phonenumber" validate:"required"`
	UserName   	 string     	 `json:"username" validate:"required"`
	Merchants  	 []*MerchantInfo  
}


type Wallet struct {
	ID         		string		 `gorm:"primary_key;"`
	MerchantID 		string		 `json:"merchantid" gorm:"foreignkey"`
	UserID 			string 		`json:"userid" gorm:"foreignkey"`
	Balance   		float64 	`json:"balance"`
}

func (um *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	uuid, err := uuid.New().MarshalText()
	if err != nil {
		return err
	}
	um.ID = string(uuid)
	return
}

