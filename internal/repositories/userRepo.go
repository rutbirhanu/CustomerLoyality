package repositories

import (
	"github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/gorm"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type UserRepo interface {
	// BuyAirtime()
	// Donate()
	// TransferCash()
	GenerateKeyPair() (string, string, error)
	CreateUser(entities.User) (*entities.User, error)
	FindUserById(string) (*entities.User, error)
	FindUserByPhone(string) (*entities.User, error)
	AddMerchant(string, string) (*entities.Merchant, *entities.User, error)
	UserLogin(user entities.UserLogin, merchantId string) (*entities.User, error)
	UpdateUser(user *entities.User) error
	// CheckBalance()
	// DeleteAll()(error)
}

type UserRepoImpl struct {
	Db           *gorm.DB
	merchantRepo MerchantRepo
}

func NewUserRepo(db *gorm.DB, merchRepo MerchantRepo) UserRepo {
	return &UserRepoImpl{
		Db:           db,
		merchantRepo: merchRepo,
	}
}

func (db *UserRepoImpl) CreateUser(user entities.User) (*entities.User, error) {

	err := db.Db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *UserRepoImpl) FindUserByPhone(phone string) (*entities.User, error) {
	user := entities.User{}

	err := db.Db.Preload("Merchants").Where("phone_number=?", phone).Take(&user).Error
	// err := db.Db.Where("phone_number=?", phone).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (db *UserRepoImpl) FindUserById(id string) (*entities.User, error) {
	user := entities.User{}
	err := db.Db.Preload("Merchants").First(&user, "id=?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}


func (db *UserRepoImpl) UserLogin(user entities.UserLogin, merchantId string) (*entities.User, error) {
	User, err := db.FindUserByPhone(user.PhoneNumber)
	if err != nil {
		return nil, err
	}
	merchant, err := db.merchantRepo.FindMerchantById(merchantId)
	if err != nil {
		return nil, err
	}

	// err = db.AddMerchant(merchantId, User.ID)
	// if err := db.Db.Model(&User).Association("Merchants").Append(&merchant); err != nil {
	// 	return nil,err
	// }
	// if err := db.Db.Model(&merchant).Association("Merchants").Append(&User); err != nil {
	// 	return nil,err
	// }
	// merchant.Users=append(merchant.Users, User)
	User.Merchants = append(User.Merchants, merchant)
	// if err := db.Db.Model(&User).Association("Merchants").Append(&merchant); err != nil {
	// 	return nil, err
	// }
	// err = db.Db.Model(&entities.User{}).Where("phone_number=?", user.PhoneNumber).Updates(&User).Error
	// db.Db.Save(&User)

	err = db.UpdateUser(User)

	if err != nil {
		return nil, err
	}
	// err = db.Db.Save(&User).Error
	// if err != nil {
	// 	return nil, err
	// }

	// merchant.Users=append(merchant.Users, User)

	return User, nil
}

func (db *UserRepoImpl) UpdateUser(user *entities.User) error {
	err := db.Db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *UserRepoImpl) AddMerchant(merchantId string, userId string) (*entities.Merchant, *entities.User, error) {

	user, err := db.FindUserById(userId)
	if err != nil {
		return nil, nil, err
	}

	merchant, err := db.merchantRepo.FindMerchantById(merchantId)
	if err != nil {
		return nil, nil, err
	}

	merchantExists := false
	for _, m := range merchant.Users {
		if m.ID == user.ID {
			merchantExists = true
			break
		}
	}

	if !merchantExists {
		merchant.Users = append(merchant.Users, user)
		// user.Merchants=append(user.Merchants, merchant)
		// err = db.Db.Model(&merchant).Association("Users").Append(&user)
		// if err != nil {
		// 	return nil, nil, err
		// }
		newWallet := entities.Wallet{
			UserID: user.ID,
			MerchantID: merchant.ID,
		}
		err= db.Db.Create(&newWallet).Error

		// err = db.Db.Model(&user).Association("Merchants").Append(&merchant)
		if err != nil {
			return nil, nil, err
		}
		// err := db.Db.Save(&merchant).Error
		// if err != nil {
		// 	return nil, nil, err
		// }
	}

	return merchant, user, nil

}


func (db *UserRepoImpl) GenerateKeyPair() (string, string, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	privateKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", string(privateKeyBytes), err
	}

	publicKeyBytes = pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(privateKeyBytes), string(publicKeyBytes), nil

}

