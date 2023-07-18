package service

import (
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
	"gorm.io/gorm"
)



type TrxService interface {
	WithTrx(trxDB *gorm.DB) TrxServiceImpl
	FindSingleWallet(userId string, merchantId string) (*entities.Wallet, bool)
	PerformTransaction(amount float64, types string, recevierId string, userWalletID string, action string) bool
	Donate(charityid string, userId string, merchantId string, amount float64) bool
	TransferPoints(amount float64, userId string, merchantId string, to string) (*entities.Wallet, bool)


	PointCollection(userPhone string, point float64, merchantId string) (*entities.Wallet, bool)
	BuyAirTime(amount float64, userId string, merchantId string) bool
	SendSMS(text string, from string, to string) (bool)
	FindSenderInfo(walletid string) (*entities.Wallet, bool)
	FindMerchantFromWallet(id string) (*entities.Merchant, bool)
	FindUserFromWallet(id string) (*entities.User, bool)
}

type TrxServiceImpl struct {
	TrxRepo repositories.TransactionRepo
}

func NewTransactionSrvc(trxRepo repositories.TransactionRepo) TrxService {
	return &TrxServiceImpl{
		TrxRepo: trxRepo,
	}
}

func (srvc TrxServiceImpl) WithTrx(trxDB *gorm.DB) TrxServiceImpl {
	srvc.TrxRepo = srvc.TrxRepo.WithTrx(trxDB)
	return srvc
}

func (srvc TrxServiceImpl) FindSingleWallet(userId string, merchantId string) (*entities.Wallet, bool) {
	wallet, err := srvc.TrxRepo.FindSingleWallet(userId, merchantId)
	if err != nil {
		return nil, false
	}
	return wallet, true
}

func (srvc TrxServiceImpl) PerformTransaction(amount float64, types string, recevierId string, userWalletID string, action string) bool {
	err := srvc.TrxRepo.PerformTransaction(amount, types, recevierId, userWalletID, action)
	if err != nil {
		return false
	}
	return true
}

func (srvc TrxServiceImpl) TransferPoints(amount float64, userId string, merchantId string, to string) (*entities.Wallet, bool) {
	wallet, err := srvc.TrxRepo.TransferPoints(amount, userId, merchantId, to)
	if err != nil {
		return nil, false
	}
	return wallet, true
}

func (srvc TrxServiceImpl) Donate(charityid string, userId string, merchantId string, amount float64) bool {
	err := srvc.TrxRepo.Donate(charityid, userId, merchantId, amount)
	if err != nil {
		return false
	}
	return true
}


func (srvc *TrxServiceImpl) PointCollection(userPhone string, point float64, merchantId string) (*entities.Wallet, bool){
	wallet,err:= srvc.TrxRepo.PointCollection(userPhone,point,merchantId)
	if err!=nil{
		return nil,false
	}
	return wallet,true
}

func (srvc *TrxServiceImpl) BuyAirTime(amount float64, userId string, merchantId string) bool{
	err:= srvc.TrxRepo.BuyAirTime(amount,userId,merchantId)
	if err!=nil{
		return false
	}
	return true
}

func (srvc *TrxServiceImpl) SendSMS(text string, from string, to string) bool{
	err := srvc.TrxRepo.SendSMS(text,from,to)
	if err!=nil{
		return false
	}
	return true
}

func (srvc *TrxServiceImpl) FindSenderInfo(walletid string) (*entities.Wallet, bool){
	wallet,err:= srvc.TrxRepo.FindSenderInfo(walletid)
	if err!=nil{
		return nil,false
	}
	return wallet,true

}

func (srvc *TrxServiceImpl) FindMerchantFromWallet(id string) (*entities.Merchant, bool){
	merchant,err:= srvc.TrxRepo.FindMerchantFromWallet(id)
	if err!=nil{
		return nil,false
	}
	return merchant,true
}

func (srvc *TrxServiceImpl) FindUserFromWallet(id string) (*entities.User, bool){
	user,err:=srvc.TrxRepo.FindUserFromWallet(id)
	if err!=nil{
		return nil,false
	}
	return user,true
}