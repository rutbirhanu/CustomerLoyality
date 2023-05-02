package repositories

import (
	"github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/gorm"
)

type MerchantRepo interface {
	CreateMerchant(entities.Merchant) (*entities.Merchant, error)
	FindMerchantById(string) (*entities.Merchant, error)
	FindMerchantByPhone(string) (*entities.Merchant, error)
	GetMerchant(entities.Merchant) *entities.Merchant


	// CreateUser(entities.User)(*entities.User, error)
	// FindUserById(string)(*entities.User,error)
	// FindUserByPhone(string)(*entities.User, error)
	// UpdateUser(entities.Merchant) (*entities.User, error)
	// FindAllUsers(from string, to string, all bool, page int, perpage int) (*entities.User,error)
}

type MerchantRepoImpl struct {
	Db *gorm.DB
}

func NewMerchantRepo(db *gorm.DB) MerchantRepo {
	return &MerchantRepoImpl{
		Db: db,
	}
}

func (db MerchantRepoImpl) GetMerchant(merchant entities.Merchant) *entities.Merchant {
	return &entities.Merchant{
		PhoneNumber:  merchant.PhoneNumber,
		Password:     merchant.Password,
		MerchantName: merchant.MerchantName,
		BusinessName: merchant.BusinessName,
	}

}

func (db *MerchantRepoImpl) CreateMerchant(merchant entities.Merchant) (*entities.Merchant, error) {
	err := db.Db.Create(&merchant).Error
	if err != nil {
		return nil, err
	}
	storedMerchant,err := db.FindMerchantByPhone(merchant.PhoneNumber)
	if err!=nil{
		return nil,err
	}

	regMerchant := db.GetMerchant(*storedMerchant)
	return regMerchant, nil

}

func (db *MerchantRepoImpl) FindMerchantById(id string) (*entities.Merchant, error) {
	var merchant entities.Merchant

	err := db.Db.Where("id=?", id).Take(&merchant).Error

	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (db *MerchantRepoImpl) FindMerchantByPhone(phone string) (*entities.Merchant, error) {
	merchant := entities.Merchant{}

	err := db.Db.Where("phone_number=?", phone).Take(&merchant).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil

}


