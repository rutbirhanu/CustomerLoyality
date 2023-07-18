package repositories

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/gorm"
)

type TransactionRepo interface {
	WithTrx(trxDB *gorm.DB) TransactionRepo
	// FindUserrr(userid string, merchantid string) (string, error)
	PointCollection(userPhone string, point float64, merchantId string) (*entities.Wallet, error)
	FindSingleWallet(userId string, merchantId string) (*entities.Wallet, error)
	PerformTransaction(amount float64, types string, recevierId string, userWalletID string, action string) error
	Donate(charityid string, userId string, merchantId string, amount float64) error
	BuyAirTime(amount float64, userId string, merchantId string) error
	SendSMS(text string, from string, to string) error
	FindSenderInfo(walletid string) (*entities.Wallet, error)
	FindMerchantFromWallet(id string) (*entities.Merchant, error)
	FindUserFromWallet(id string) (*entities.User, error)
	TransferPoints(amount float64, userId string, merchantId string, to string) (*entities.Wallet, error)
}

type TransactionRepoImpl struct {
	Db           *gorm.DB
	UserRepo     UserRepo
	MerchantRepo MerchantRepo
}

func NewTransactionRepo(db *gorm.DB, userRepo UserRepo, merchantRepo MerchantRepo) TransactionRepo {
	return &TransactionRepoImpl{
		Db:           db,
		UserRepo:     userRepo,
		MerchantRepo: merchantRepo,
	}
}

func (db *TransactionRepoImpl) WithTrx(trxDB *gorm.DB) TransactionRepo {
	if trxDB == nil {
		log.Print("txn db is not found")
	}
	db.Db = trxDB
	return db

}

func (db *TransactionRepoImpl) SendSMS(text string, from string, to string) (error) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	baseURL := os.Getenv("SMS_BASEURL")
	params := url.Values{}
	params.Set("username", os.Getenv("SMS_USERNAME"))
	params.Set("password", os.Getenv("SMS_PASSWORD"))
	params.Set("text", text)
	params.Set("from", from)
	params.Set("to", to)

	// Create an HTTP client
	client := &http.Client{}

	// Create the request with the parameters
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return err
	}
	req.URL.RawQuery = params.Encode()

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to send SMS. Status code: %d", resp.StatusCode)
	}
	
	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Print the response body
	fmt.Println(string(body))

	return nil
}

func (db *TransactionRepoImpl) FindMerchantFromWallet(id string) (*entities.Merchant, error) {
	merchant, err := db.MerchantRepo.FindMerchantById(id)
	if err != nil {
		return nil, err
	}
	return merchant, nil
}

func (db *TransactionRepoImpl) FindUserFromWallet(id string) (*entities.User, error) {
	user, err := db.UserRepo.FindUserById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (db *TransactionRepoImpl) FindSenderInfo(walletid string) (*entities.Wallet, error) {
	wallet := entities.Wallet{}
	err := db.Db.Model(&entities.Wallet{}).Where("id=?", walletid).Find(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (db *TransactionRepoImpl) PointCollection(userPhone string, point float64, merchantId string) (*entities.Wallet, error) {
	user, err := db.UserRepo.FindUserByPhone(userPhone)
	if err != nil {
		return nil, err
	}
	merchant, err := db.MerchantRepo.FindMerchantById(merchantId)
	if err != nil {
		return nil, err
	}
	ratio:=merchant.PointConfiguration

	userWallet := entities.Wallet{}

	result := db.Db.Table("wallets").
		Where("user_id = ? AND merchant_id = ?", user.ID, merchant.ID).
		Find(&userWallet)

	if result.Error != nil {
		return nil, result.Error
	}

	userWallet.Balance += (point*ratio)
	db.Db.Save(userWallet)

	err = db.PerformTransaction(point, "point collection", "", userWallet.ID, "credit")
	if err != nil {
		return &userWallet, err
	}
	return &userWallet, nil

}

func (db *TransactionRepoImpl) FindSingleWallet(userId string, merchantId string) (*entities.Wallet, error) {
	userWallet := entities.Wallet{}
	result := db.Db.Table("wallets").Where("user_id=? AND merchant_id=?", userId, merchantId).First(&userWallet)

	if result.Error != nil {
		return nil, result.Error
	}
	return &userWallet, nil

}

func (db *TransactionRepoImpl) PerformTransaction(amount float64, types string, recevierId string, userWalletId string, action string) error {

	transaction := entities.Transaction{
		Amount:         amount,
		Type:           types,
		Action:         action,
		UserMerchantID: userWalletId,
		ReceiverID:     recevierId,
	}

	err := db.Db.Create(&transaction).Error
	if err != nil {
		return errors.New("error while creating transaction")
	}

	return nil

}

func (db *TransactionRepoImpl) Donate(charityid string, userId string, merchantId string, amount float64) error {
	var account string
	var id string

	for _, value := range entities.Organizations {
		if charityid == value["id"] {
			account = value["account"]
			id = value["id"]
			break
		}
	}

	userWallet, err := db.FindSingleWallet(userId, merchantId)
	if err != nil {
		return err
	}

	if userWallet.Balance < amount {
		return errors.New("not enough balance")
	}
	fmt.Print(account)

	///and increase the charity balance

	userWallet.Balance -= amount
	db.Db.Save(&userWallet)
	err = db.PerformTransaction(amount, "donation", id, userWallet.ID, "debit")
	if err != nil {
		return err
	}

	return nil

}

func (db *TransactionRepoImpl) BuyAirTime(amount float64, userId string, merchantId string) error {
	userWallet, err := db.FindSingleWallet(userId, merchantId)
	if err != nil {
		return nil
	}
	if userWallet.Balance < amount {
		return errors.New("amount not sufficient")
	}
	userWallet.Balance -= amount

	// call the air time api

	err = db.PerformTransaction(amount, "air time", "", userWallet.ID, "debit")
	if err != nil {
		return err
	}

	return nil
}

func (db *TransactionRepoImpl) TransferPoints(amount float64, userId string, merchantId string, toPhoneNumber string) (*entities.Wallet, error) {
	var recevierId string
	userWallet, err := db.FindSingleWallet(userId, merchantId)
	if err != nil {
		return userWallet, errors.New("wallet not found")
	}
	if userWallet.Balance < amount {
		return userWallet, errors.New("amount not sufficient")
	}
	merchant := entities.Merchant{}
	err = db.Db.Model(&entities.Merchant{}).Preload("Users", "phone_number=?", toPhoneNumber).Where("id=?", merchantId).First(&merchant).Error
	if err != nil {
		return userWallet, err
	}

	for _, user := range merchant.Users {
		recevierId = user.ID
		if recevierId == userId {
			return nil, errors.New("can not transfer to current user")
		}
	}

	receiverWallet := entities.Wallet{}
	err = db.Db.Model(&entities.Wallet{}).Where("user_id=? AND merchant_id=?", recevierId, merchantId).First(&receiverWallet).Error
	if err != nil {
		return userWallet, err
	}
	userWallet.Balance -= amount
	receiverWallet.Balance += amount
	db.Db.Save(&userWallet)
	db.Db.Save(&receiverWallet)
	err = db.PerformTransaction(amount, "point transfer", receiverWallet.ID, userWallet.ID, "debit")
	if err != nil {
		return userWallet, err
	}

	return userWallet, nil
}

// func (db *TransactionRepoImpl) FindUserrr(userid string, merchantid string) (string, error) {
// 	var useridd string
// 	merchant := entities.Merchant{}
// 	err := db.Db.Model(&entities.Merchant{}).Preload("Users", "id=?", userid).Where("id=?", merchantid).First(&merchant).Error
// 	if err != nil {
// 		return "", err
// 	}
// 	for _, users := range merchant.Users {
// 		// Access the fields of the user model here
// 		useridd = users.ID
// 		fmt.Println(users)
// 	}
// 	return useridd, nil
// }
