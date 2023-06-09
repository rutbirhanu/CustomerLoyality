package repositories

import (
	"errors"
	"fmt"
	"log"

	"github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/gorm"
)

type TransactionRepo interface {
	WithTrx(trxDB *gorm.DB) TransactionRepo
	FindUserrr(userid string, merchantid string) (string, error)
	FindSingleWallet(userId string, merchantId string) (*entities.Wallet, error)
	PerformTransaction(amount float64, types string, recevierId string, userWalletID string) error
	Donate(charityid string, userId string, merchantId string, amount float64) error
	BuyAirTime(amount float64, userId string, merchantId string) error
	TransferPoints(amount float64, userId string, merchantId string, to string) (*entities.Wallet, error)
}

type TransactionRepoImpl struct {
	Db       *gorm.DB
	UserRepo UserRepo
}

func NewTransactionRepo(db *gorm.DB, userRepo UserRepo) TransactionRepo {
	return &TransactionRepoImpl{
		Db:       db,
		UserRepo: userRepo,
	}
}

func (db *TransactionRepoImpl) WithTrx(trxDB *gorm.DB) TransactionRepo {
	if trxDB == nil {
		log.Print("txn db not founc")
	}
	db.Db = trxDB
	return db

}

func (db *TransactionRepoImpl) FindSingleWallet(userId string, merchantId string) (*entities.Wallet, error) {
	userWallet := entities.Wallet{}
	result := db.Db.Table("wallets").Where("user_id=? AND merchant_id=?", userId, merchantId).First(&userWallet)

	if result.Error != nil {
		return nil, result.Error
	}
	return &userWallet, nil

}

func (db *TransactionRepoImpl) PerformTransaction(amount float64, types string, recevierId string, userWalletId string) error {

	transaction := entities.Transaction{
		Amount:         amount,
		Type:           types,
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
	err = db.PerformTransaction(amount, "donation", id, userWallet.ID)
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

	err = db.PerformTransaction(amount, "air time", "", userWallet.ID)
	if err != nil {
		return err
	}

	return nil
}

func (db *TransactionRepoImpl) TransferPoints(amount float64, userId string, merchantId string, toPhoneNumber string) (*entities.Wallet, error) {
	var recevierId string
	userWallet, err := db.FindSingleWallet(userId, merchantId)
	if err != nil {
		return userWallet, err
	}
	if userWallet.Balance < amount {
		return userWallet, errors.New("amount not sufficient")
	}
	merchant := entities.Merchant{}
	err = db.Db.Model(&entities.Merchant{}).Preload("Users", "phone_number=?", toPhoneNumber).Where("id=?", merchantId).First(&merchant).Error
	if err != nil {
		return userWallet, errors.New("receiver user not found ")
	}

	for _, user := range merchant.Users {
		recevierId = user.ID
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
	err = db.PerformTransaction(amount, "point transfer", receiverWallet.ID, userWallet.ID)
	if err != nil {
		return userWallet, err
	}

	return userWallet, nil
}

func (db *TransactionRepoImpl) FindUserrr(userid string, merchantid string) (string, error) {
	var useridd string
	merchant := entities.Merchant{}
	err := db.Db.Model(&entities.Merchant{}).Preload("Users", "id=?", userid).Where("id=?", merchantid).First(&merchant).Error
	if err != nil {
		return "", err
	}
	for _, users := range merchant.Users {
		// Access the fields of the user model here
		useridd = users.ID
		fmt.Println(users)
	}
	return useridd, nil
}
