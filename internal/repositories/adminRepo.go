package repositories

import (
	"log"

	"github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/gorm"
)

type AdminRepo interface {

	// FindAllMerchant(from string, to string, perpage int, all bool, page int) (*entities.Merchant, error)
	// DeleteMerchant(string) (*entities.Merchant,error)
	WithTrx(*gorm.DB) AdminRepo
	FindWalletById(string) (*entities.Wallet, error)
	AddMerchant(string, string) (*entities.Merchant, *entities.User, error)
	PointCollection(userPhone string, point float64, merchantId string) (*entities.Wallet, error)
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

func (repo *AdminRepoImpl) WithTrx(trxDB *gorm.DB) AdminRepo {
	if trxDB == nil {
		log.Printf("transaction db not found")
		return repo
	}
	repo.Db = trxDB
	return repo
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

	// for _, m := range user.Merchants {
	// 	if m.ID == merchant.ID {
	// 		merchantExists = true
	// 		break
	// 	}
	// }

	// if !merchantExists {
	// 	user.Merchants = append(user.Merchants, merchant)
	// 	err = db.Db.Model(&user).Association("Merchants").Append(&merchant)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

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

func (db *AdminRepoImpl) PointCollection(userPhone string, point float64, merchantId string) (*entities.Wallet, error) {
	user, err := db.userRepo.FindUserByPhone(userPhone)
	if err != nil {
		return nil, err
	}
	merchant, err := db.merchantRepo.FindMerchantById(merchantId)
	if err != nil {
		return nil, err
	}

	userWallet := entities.Wallet{}

	result := db.Db.Table("wallets").
		Where("user_id = ? AND merchant_id = ?", user.ID, merchant.ID).
		Find(&userWallet)

	if result.Error != nil {
		return nil, result.Error
	}

	userWallet.Balance += point
	db.Db.Save(userWallet)
	return &userWallet, nil

}
