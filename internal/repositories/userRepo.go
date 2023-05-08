package repositories

import (
	"github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/gorm"
)

type UserRepo interface {
	// BuyAirtime()
	// Donate()
	// TransferCash()
	CreateUser(entities.User, string) (*entities.User, error)
	FindUserById(string) (*entities.User, error)
	FindUserByPhone(string) (*entities.User, error)
	AddMerchant(merchantId string, userId string) (*entities.Merchant, error)
}

type UserRepoImpl struct {
	Db           *gorm.DB
	merchantRepo MerchantRepo
}

func NewUserRepo(db *gorm.DB, merchantRepo MerchantRepo) UserRepo {
	return &UserRepoImpl{
		Db:           db,
		merchantRepo: merchantRepo,
	}
}

func (db *UserRepoImpl) CreateUser(user entities.User, merchantId string) (*entities.User, error) {
	err := db.Db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	foundUser, err := db.FindUserByPhone(user.PhoneNumber)

	if err != nil {
		return nil, err
	}
	// err = db.AddMerchant(merchantId, user.ID)
	// if err != nil {
	// 	return nil, err
	// }

	return foundUser, nil

}


func (db *UserRepoImpl) FindUserByPhone(phone string) (*entities.User, error) {
	user := entities.User{}

	err := db.Db.Where("phone_number=?", phone).Preload("Merchants").Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (db *UserRepoImpl) FindUserById(id string) (*entities.User, error) {
	user := entities.User{}
	err := db.Db.Where("id=?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *UserRepoImpl) AddMerchant(merchantId string, userId string) (*entities.Merchant, error) {

	user, err := db.FindUserById(userId)
	if err != nil {
		return nil, err
	}
	merchant, err := db.merchantRepo.FindMerchantById(merchantId)
	if err != nil {
		return nil, err
	}

	if err := db.Db.Model(&merchant).Association("Users").Append(user); err != nil {
		return nil, err
	}

	if err := db.Db.Model(&user).Association("Merchants").Append(merchant); err != nil {
		return nil, err
	}
	return merchant, nil
}
