package entities


import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Model
	PhoneNumber string      `  json:"phonenumber" validate:"required"`
	UserName    string      `json:"username" validate:"required"`
	Merchants   []*Merchant `gorm:"many2many:merchant_users;"`
	Balance     float64     `json:"balance"`
	// Transaction 	[]*Transaction
}

type UserLogin struct {
	PhoneNumber string `json:"phonenumber"`
}


type MerchantUser struct {
	ID         string `gorm:"primary_key;"`
	MerchantID string `  json:"merchantid"`
	// Merchant   Merchant `json:"merchant"`
	UserID string `  json:"userid"`
	// User       User     `json:"user"`
}

func (um *MerchantUser) BeforeCreate(tx *gorm.DB) (err error) {
	uuid, err := uuid.New().MarshalText()
	if err != nil {
		return err
	}
	um.ID = string(uuid)
	return
}



// type MerchantUser struct{
// 	MerchantID		string
// 	UserID			string
// }