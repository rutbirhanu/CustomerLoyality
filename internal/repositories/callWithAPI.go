package repositories

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt"
	"github.com/matoous/go-nanoid/v2"

	// "github.com/google/uuid"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/util"
	"gorm.io/gorm"
)

type ApiRepo interface {
	WithTrx(trxDb *gorm.DB) ApiRepo
	GenerateNewToken(merchantid string) (string, error)
	RemoveToken(token string) error
	GiveWallet(phone string, username string, merchantId string) (*entities.User, error)
	PointConfiguration(ratio float64, merchantid string) error
	GivePoint(usersPhone string, amount float64, merchantId string) error
	BuyGiftCard(merchantid string, amount float64, recipentPhone string, purchaserPhone string) error
	FindGiftCardByCode( giftcardCode string)(*entities.GiftCard,error)
	redeemGiftCard(merchantId string, giftcardCode string, totalPrice float64) (float64,error)
}

type ApiRepoImpl struct {
	Db           *gorm.DB
	MerchantRepo MerchantRepo
	UserRepo     UserRepo
	TrxRepo      TransactionRepo
}

type ClaimData struct {
	Phone      string
	Role       util.Role
	MerchantId string
	jwt.StandardClaims
}

func NewApiRepo(db *gorm.DB, merchantRepo MerchantRepo, userRepo UserRepo, transactionRepo TransactionRepo) ApiRepo {
	return &ApiRepoImpl{
		Db:           db,
		MerchantRepo: merchantRepo,
		UserRepo:     userRepo,
		TrxRepo:      transactionRepo,
	}
}

func (repo *ApiRepoImpl) WithTrx(trxDb *gorm.DB) ApiRepo {
	if trxDb == nil {
		log.Print("trx db not found")
	}
	repo.Db = trxDb
	return repo
}

//token generation when they want to
func (repo ApiRepoImpl) GenerateNewToken(merchantid string)(string,error){

	nanoid, err := gonanoid.New()
	fmt.Print(nanoid)
	if err != nil {
		return "",err
	}

	tokenTable := &entities.TokenTable{
		Token:      nanoid,
		MerchantId: merchantid,
	}

	err = repo.Db.Create(tokenTable).Error
	if err != nil {
		return "", err
	}

	return nanoid,nil
}

// func (repo ApiRepoImpl) GenerateNewToken(merchantid string) (string, error) {
// 	merchant,err := repo.MerchantRepo.FindMerchantById(merchantid)
// 	if err!=nil{
// 		return "",err
// 	}

// 	claims := &ClaimData{
// 		Phone:      merchant.PhoneNumber,
// 		Role:       util.Merchant,
// 		MerchantId: merchantid,
// 		StandardClaims: jwt.StandardClaims{
// 			IssuedAt:  time.Now().Unix(),
// 			ExpiresAt: time.Now().Add(time.Hour * time.Duration(2)).Unix(),
// 		},
// 	}

// 	key, err := repo.MerchantRepo.RetrivePrivateKey(merchant.PhoneNumber)
// 	if err != nil {
// 		return "", err
// 	}
// 	// privateKey,err:= jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
// 	// if err!=nil{
// 	// 	return "",err
// 	// }
// 	privateKey, err := x509.ParseECPrivateKey([]byte(key))
// 	if err != nil {
// 		return "", err
// 	}

// 	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(privateKey)
// 	if err != nil {
// 		return "", err
// 	}

// 	tokenTable := &entities.TokenTable{
// 		Token:      token,
// 		MerchantId: merchantid,
// 	}

// 	err = repo.Db.Create(tokenTable).Error
// 	if err != nil {
// 		return "", err
// 	}

// 	return token, nil
// }

// able to delete the token

func (repo ApiRepoImpl) RemoveToken(token string) error {
	Token := &entities.TokenTable{}

	err := repo.Db.Model(&entities.TokenTable{}).Where("token=?", token).Find(Token).Error
	if err != nil {
		return err
	}

	err = repo.Db.Delete(Token.ID).Error
	if err != nil {
		return err
	}
	return nil
}

//autenticate with token

// automate the registering of users
//create end point /user/userId/giveWallet

func (repo ApiRepoImpl) GiveWallet(phone string, username string, merchantId string) (*entities.User, error) {
	user, err := repo.UserRepo.FindUserByPhone(phone)

	if err == nil {
		User, err := repo.MerchantRepo.CreateUser(*user, merchantId)
		if err != nil {
			return nil, err
		}
		return User, nil
	} else {
		err = repo.Db.Create(&user).Error
		if err != nil {
			return nil, err
		}
		_, _, err = repo.MerchantRepo.AddUserToMerchant(merchantId, user.ID)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

}

//automate giving points

func (repo ApiRepoImpl) GivePoint(usersPhone string, amount float64, merchantId string) error {
	merchant, err := repo.MerchantRepo.FindMerchantById(merchantId)
	if err != nil {
		return err
	}
	user, err := repo.UserRepo.FindUserByPhone(usersPhone)
	if err != nil {
		return err
	}

	userExists := false
	for _, m := range merchant.Users {
		if m.ID == user.ID {
			userExists = true
			break
		}
	}

	if !userExists {
		return errors.New("recipent user not found")
	}
	wallet := entities.Wallet{}
	err = repo.Db.Model(&entities.Wallet{}).Where("merchant_id=? AND  user_id=?").First(wallet).Error
	if err != nil {
		return err
	}
	wallet.Balance += (amount * merchant.PointConfiguration)
	return nil
}

// point configuration
func (repo ApiRepoImpl) PointConfiguration(ratio float64, merchantid string) error {
	merchant, err := repo.MerchantRepo.FindMerchantById(merchantid)
	if err != nil {
		return err
	}
	merchant.PointConfiguration = ratio
	return nil
}

//buy gift card

func (repo ApiRepoImpl) BuyGiftCard(merchantid string, amount float64, recipentPhone string, purchaserPhone string) error {
var reciverId string

	purchaseWallet,err:= repo.TrxRepo.FindSingleWallet(purchaserPhone,merchantid)
	if err!=nil{
		return err
	}
	if purchaseWallet.Balance < amount{
return errors.New("your balance is insufficient")
	}

	nanoid, err := gonanoid.New()
	// uuid, err := uuid.New().MarshalText()
	fmt.Print(nanoid)
	if err != nil {
		return err
	}
purchaseWallet.Balance-=amount
recipent,err:= repo.UserRepo.FindUserByPhone(recipentPhone)
if err!=nil{
reciverId=""
}else{
	reciverId=recipent.ID
}

err = repo.TrxRepo.PerformTransaction(amount,"debit",reciverId,purchaseWallet.ID,"buy gift card")
if err!=nil{
	return errors.New("failed to record transaction ")
}
	record := entities.GiftCard{
		Amount:        amount,
		RecipentPhone: recipentPhone,
		MerchantId:    merchantid,
		GiftCardCode:  nanoid,
		PurchaserPhone:purchaserPhone,
	}

	err = repo.Db.Create(&record).Error
	if err != nil {
		return err
	}
	return nil
}

// send gift card
// in the handler

//redeemGiftcard
func (repo ApiRepoImpl) FindGiftCardByCode( giftcardCode string)(*entities.GiftCard,error){
	giftCard := entities.GiftCard{}
	err:= repo.Db.Model(&entities.GiftCard{}).Where("gift_card_code=?",giftcardCode).First(giftCard).Error
	if err!=nil{
		return nil,err
	}

return &giftCard,nil
}

func (repo ApiRepoImpl) redeemGiftCard(merchantId string, giftcardCode string, totalPrice float64) (float64,error){
	merchant,err := repo.MerchantRepo.FindMerchantById(merchantId)
	if err!=nil{
		return -1,err
	}
	giftCard,err:= repo.FindGiftCardByCode(giftcardCode)
	if err!=nil{
		return -1,err
	}
	giftCardAmount := giftCard.Amount * merchant.PointConfiguration
	if totalPrice > giftCardAmount{
		totalPrice-=giftCardAmount
		giftCard.Amount=0
		return totalPrice,nil

	}else{
		giftCardAmount-=totalPrice
		return 0,nil
	}

}
