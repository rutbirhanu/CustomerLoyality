package repositories

import (
	"github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/gorm"
)

type WalletRepo interface {
	FindWalletById(string)  (*entities.Wallet, error)
	// CreateUserMerchant(entities.UserMerchant) (*entities.UserMerchant, error)
	// AddMerchant(string, string) (*entities.Merchant, *entities.User, error)
}

type WalletRepoImpl struct {
	Db           *gorm.DB
}

func NewWalletRepo(db *gorm.DB) WalletRepo {
	return &WalletRepoImpl{
		Db: db,
	}
}


func (db *WalletRepoImpl) FindWalletById(id string) (*entities.Wallet, error){
userMerchant := entities.Wallet{}
if err:= db.Db.Find(&userMerchant, "id=?", id).Error; err!=nil{
	return nil,err
}

return &userMerchant, nil
}



// func (db *UserMerchantRepoImpl) AddMerchant(merchantId string, userId string) (*entities.Merchant, *entities.User, error) {
// 	merchant, err := db.merchantRepo.FindMerchantByPhone(merchantId)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	user, err := db.userRepo.FindUserByPhone(userId)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	db.Db.Omit("Users").Updates(&merchant)
// 	merchant.Users = append(merchant.Users, user)

// 	db.Db.Session(&gorm.Session{FullSaveAssociations: true}).Model(&merchant).Association("User").Replace(&merchant.Users)

// 	return merchant, user, nil
// }
