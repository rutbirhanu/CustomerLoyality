package repositories

import (
	"github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/gorm"
)

type MerchantUserRepo interface {
	FindMerchantUserById(string)  (*entities.MerchantUser, error)
	// CreateUserMerchant(entities.UserMerchant) (*entities.UserMerchant, error)
	// AddMerchant(string, string) (*entities.Merchant, *entities.User, error)
}

type MerchantUserRepoImpl struct {
	Db           *gorm.DB
}

func NewMerchantUserRepo(db *gorm.DB) MerchantUserRepo {
	return &MerchantUserRepoImpl{
		Db: db,
	}
}


func (db *MerchantUserRepoImpl) FindMerchantUserById(id string) (*entities.MerchantUser, error){
userMerchant := entities.MerchantUser{}
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
