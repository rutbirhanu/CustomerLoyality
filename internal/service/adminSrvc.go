package service

import (
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
	"gorm.io/gorm"
)

type AdminService interface {
	WithTrx(*gorm.DB) AdminServiceImpl
	FindWalletById(string) (*entities.Wallet, bool)
	AddUserToMerchant(string, string) (*entities.Merchant, *entities.User, bool)
	// PointCollection(userPhone string, point float64, merchantId string) (*entities.Wallet, bool)
}

type AdminServiceImpl struct {
	AdminRepo repositories.AdminRepo
}

func NewAdminService(adminRepo repositories.AdminRepo) AdminService {
	return &AdminServiceImpl{
		AdminRepo: adminRepo,
	}
}

func (srvc AdminServiceImpl) WithTrx(trxHandle *gorm.DB) AdminServiceImpl {
	srvc.AdminRepo = srvc.AdminRepo.WithTrx(trxHandle)
	return srvc
}

func (srvc AdminServiceImpl) FindWalletById(id string) (*entities.Wallet, bool) {
	wallet, err := srvc.AdminRepo.FindWalletById(id)
	if err != nil {
		return nil, false
	}
	return wallet, true
}

func (srvc AdminServiceImpl) AddUserToMerchant(merchantid string, userid string) (*entities.Merchant, *entities.User, bool) {
	merchant, user, err := srvc.AdminRepo.AddUserToMerchant(merchantid, userid)
	if err != nil {
		return nil, nil, false
	}
	return merchant, user, true
}


// func (srvc AdminServiceImpl) PointCollection(userPhone string, point float64, merchantId string) (*entities.Wallet, bool) {
// 	wallet, err := srvc.AdminRepo.PointCollection(userPhone, point, merchantId)
// 	if err != nil {
// 		return nil, false
// 	}
// 	return wallet, true
// }
