package entities

type Merchant struct {
	Model
	MerchantName 	string 		` validate:" required , max=30"  json:"name" `
	Password     	string 		` validate:" required , max=30, min=6 "  json:"password" `
	PhoneNumber  	string 		`  validate:" required "  json:"phonenumber" gorm:"size:10;" `
	Token        	string 		`json:"token"`
	BusinessName 	string 		` validate:" required , max=30"  json:"businessname" `
	PrivateKey   	string
	PublicKey    	string
	TokenTable		[]*TokenTable 		`gorm:"foreignkey:ID"`
	PointConfiguration 		float64		`json:"pointConfig,omitempty" gorm:"default:1"`
	Users        []*User 	`gorm:"many2many:wallets;"`
}

type PointConfig struct{
	PointConfiguration 		float64		`json:"pointConfig,omitempty" gorm:"default:1"`

}

type CreatedMerchantResponse struct {
	MerchantName	 string 	
	Password    	 string 	
	PhoneNumber 	 string 	
	BusinessName	 string 	
	Users       	 []*UserInfo
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

type TokenTable struct{
	Model
	Token 			string		`json:"token"`
	MerchantId 		string			`json:"merchantid" gorm:"foreignkey:ID"`
}

type GetAllUsers struct {
	Users    []User
	Total    int64
	Page     int64
	Prev     int64
	Next     int64
	LastPage int64
}
