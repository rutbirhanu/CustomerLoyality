package entities

import "regexp"

type Merchant struct {
	Model
	MerchantName       string ` validate:"required,max=30"  json:"name" `
	Password           string ` validate:"required,max=30,min=6"  json:"password" `
	PhoneNumber        string `  validate:"required"  json:"phonenumber" gorm:"size:10;" `
	Token              string `json:"token"`
	BusinessName       string ` validate:"required,max=30"  json:"businessname" `
	PrivateKey         string
	PublicKey          string
	ApiToken           []*TokenTable
	PointConfiguration float64 `json:"pointConfig,omitempty" gorm:"default:1"`
	Users              []*User `gorm:"many2many:wallets;"`
}

func (m *Merchant) ValidatePhoneNumber() bool {
	// Define a regex pattern for 10 digits starting with "219"
	regexPattern := `^215\d{5}$`

	// Compile the regex pattern
	regex := regexp.MustCompile(regexPattern)

	// Check if the phone number matches the pattern
	return regex.MatchString(m.PhoneNumber)
}

type PointConfig struct {
	PointConfiguration float64 `json:"pointConfig,omitempty" gorm:"default:1"`
}

type CreatedMerchantResponse struct {
	MerchantName       string
	Password           string
	PhoneNumber        string
	BusinessName       string
	PointConfiguration float64
	Users              []*UserInfo
}

type MerchantInfo struct {
	BusinessName string ` validate:" required , max=30"  json:"businessname" `
}

type MerchantLogin struct {
	Password    string ` validate:" required , max=30 , min=6" json:"password"`
	PhoneNumber string ` validate:" required , max=30 " json:"phonenumber"`
	Token       string `json:"token"`

	// 	RefreshToken	string 		`json:"refreshToken"`
}

type TokenTable struct {
	Model
	Token      string `json:"token"`
	MerchantId string `json:"merchantid" gorm:"foreignkey:ID"`
}

type GetAllUsers struct {
	Users    []User
	Total    int64
	Page     int64
	Prev     int64
	Next     int64
	LastPage int64
}
