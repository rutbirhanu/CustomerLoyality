package entities

type Merchant struct {
	Model
	MerchantName string ` validate:" required , max=30"  json:"name" `
	Password     string ` validate:" required , max=30, min=6 "  json:"password" `
	PhoneNumber  string `  validate:" required "  json:"phonenumber" gorm:"size:10;" `
	Token        string `json:"token"`
	BusinessName string ` validate:" required , max=30"  json:"businessname" `
	PrivateKey   string
	PublicKey    string
	Users        []*User `gorm:"many2many:wallets;"`
}

type CreatedMerchantResponse struct {
	MerchantName string ` validate:" required , max=30"  json:"name" `
	Password     string ` validate:" required , max=30, min=6 "  json:"password" `
	PhoneNumber  string `  validate:" required "  json:"phonenumber" `
	BusinessName string ` validate:" required , max=30"  json:"businessname" `
	Users        []*MerchantsUserResponse
}

type UsersMerchantResponse struct {
	Model
	MerchantName string ` validate:" required , max=30"  json:"name" `
	BusinessName string ` validate:" required , max=30"  json:"businessname" `
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
