package repositories

import (
	// "github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/gorm"
)

type UserMerchantRepo interface {

	// CreateUserMerchant(entities.UserMerchant) (*entities.UserMerchant, error)
	// AddMerchant( string, string) (*entities.Merchant, *entities.User, error)
}

type UserMerchantRepoImpl struct {
	Db           *gorm.DB
	merchantRepo MerchantRepo
	userRepo     UserRepo
}

func NewUserMerchantRepo(db *gorm.DB, mer MerchantRepo, uss UserRepo) UserMerchantRepo {
	return &UserMerchantRepoImpl{
		Db:           db,
		merchantRepo: mer,
		userRepo:     uss,
	}
}

// func (db *UserMerchantRepoImpl) AddMerchant(merchantId string,  user entities.User) (*entities.Merchant, *entities.User, error) {
// 	merchant,err:= db.merchantRepo.FindMerchantById(merchantId)
// 	if err!=nil{
// 		return nil,nil,err
// 	}

// 	merchant.Users = append(merchant.Users, &user)
// 	db.Db.Save(&merchant)

// 	return merchant, &user, nil
// }

