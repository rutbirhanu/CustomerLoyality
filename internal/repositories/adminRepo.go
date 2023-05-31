package repositories

import (
	"github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/gorm"
)

type AdminRepo interface {

	// FindAllMerchant(from string, to string, perpage int, all bool, page int) (*entities.Merchant, error)
	// DeleteMerchant(string) (*entities.Merchant,error)
	FindWalletById(string) (*entities.Wallet, error)
	AddMerchant(string, string) (*entities.Merchant, *entities.User, error)
}

type AdminRepoImpl struct {
	Db           *gorm.DB
	userRepo     UserRepo
	merchantRepo MerchantRepo
}

func NewAdminRepo(db *gorm.DB, user UserRepo, merchant MerchantRepo) AdminRepo {
	return &AdminRepoImpl{
		Db:           db,
		userRepo:     user,
		merchantRepo: merchant,
	}
}

func (db *AdminRepoImpl) FindWalletById(id string) (*entities.Wallet, error) {
	userMerchant := entities.Wallet{}
	if err := db.Db.Find(&userMerchant, "id=?", id).Error; err != nil {
		return nil, err
	}

	return &userMerchant, nil
}

func (db *AdminRepoImpl) AddMerchant(merchantId string, userId string) (*entities.Merchant, *entities.User, error) {

	user, err := db.userRepo.FindUserById(userId)
	if err != nil {
		return nil, nil, err
	}

	merchant, err := db.merchantRepo.FindMerchantById(merchantId)
	if err != nil {
		return nil, nil, err
	}

	merchant.Users = append(merchant.Users, user)

	newWallet := entities.Wallet{
		UserID:     user.ID,
		MerchantID: merchant.ID,
	}
	err = db.Db.Create(&newWallet).Error

	if err != nil {
		return nil, nil, err
	}

	return merchant, user, nil

}
