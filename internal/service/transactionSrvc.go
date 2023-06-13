package service

import (
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
	"gorm.io/gorm"
)

type TrxService interface {
	WithTrx(trxDB *gorm.DB) TrxServiceImpl
	// FindUserrr(userid string, merchantid string) (string, error)
	FindSingleWallet(userId string, merchantId string) (*entities.Wallet, bool)
	PerformTransaction(amount float64, types string, recevierId string, userWalletID string, action string) bool
	Donate(charityid string, userId string, merchantId string, amount float64) bool
	// BuyAirTime(amount float64, userId string, merchantId string) bool
	TransferPoints(amount float64, userId string, merchantId string, to string) (*entities.Wallet, bool)
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
