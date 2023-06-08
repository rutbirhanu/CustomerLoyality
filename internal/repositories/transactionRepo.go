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
	FindSingleWallet(userId string, merchantId string) (*entities.Wallet, error)
	PerformTransaction(amount float64, types string, recevierId string, userWalletID string) error
	Donate(charityid string, userId string, merchantId string, amount float64) error
	BuyAirTime(amount float64, userId string, merchantId string) error
	TransferPoints(amount float64, userId string, merchantId string, to string) error
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
		ReceiverID: recevierId,
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

	// have to optimize this operation

	for _, value := range entities.Organizations {
		if charityid == value["id"] {
			account = value["account"]
			id=value["id"]
		} else {
			return errors.New("charity account not found")
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
	err = db.PerformTransaction(amount,"donation",id,userWallet.ID)
	if err!=nil{
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

	err = db.PerformTransaction(amount,"air time","",userWallet.ID)
	if err!=nil{
		return err
	}

	return nil
}

func (db *TransactionRepoImpl) TransferPoints(amount float64, userId string, merchantId string, toPhoneNumber string) error {
	userWallet, err := db.FindSingleWallet(userId, merchantId)
	if err != nil {
		return nil
	}
	if userWallet.Balance < amount {
		return errors.New("amount not sufficient")
	}

	toUser := &entities.User{}
	err = db.Db.Model(&entities.Merchant{}).Where("id = ?", merchantId).
		Preload("Users", "phone_number = ?", toPhoneNumber).
		First(&toUser).
		Error
	if err != nil {
		return errors.New("receiver user not found ")
	}

	receiverWallet := entities.Wallet{}
	err = db.Db.Model(&entities.Wallet{}).Where("user_id=? AND merchant_id=?", userId, merchantId).First(&receiverWallet).Error
	if err != nil {
		return err
	}
	userWallet.Balance-=amount
	receiverWallet.Balance+=amount

	err = db.PerformTransaction(amount,"point transfer",receiverWallet.ID,userWallet.ID)
	if err!=nil{
		return err
	}

	return nil
}
