package repositories

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math"

	"github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/gorm"
)

type MerchantRepo interface {
	CreateMerchant(entities.Merchant) (*entities.Merchant, error)
	FindMerchantById(string) (*entities.Merchant, error)
	FindMerchantByPhone(string) (*entities.Merchant, error)
	RetrivePublicKey(phone string) (string, error)
	PointCollection(userPhone string, point float64, merchantId string) (*entities.Wallet, error)

	GenerateKeyPair() (string, string, error)
	UpdateMerchant(entities.Merchant) error
	GetAllMerchants() (*[]entities.Merchant, error)
	DeleteAll() error
	CreateUser(entities.User, string) (*entities.User, error)

	FindAllUsers(from string, to string, all bool, page int64, perpage int64) (*entities.GetAllUsers, error)
	// UpdataUserPoints(string , float64) error
}

type MerchantRepoImpl struct {
	Db              *gorm.DB
	UserRepo        UserRepo
	TransactionRepo TransactionRepo
}

func NewMerchantRepo(db *gorm.DB, userRepo UserRepo, trxRepo TransactionRepo) MerchantRepo {
	return &MerchantRepoImpl{
		Db:              db,
		UserRepo:        userRepo,
		TransactionRepo: trxRepo,
	}
}

func (db *MerchantRepoImpl) AddUserToMerchant(merchantId string, userId string) (*entities.Merchant, *entities.User, error) {

	user, err := db.UserRepo.FindUserById(userId)
	if err != nil {
		return nil, nil, err
	}
	merchant, err := db.FindMerchantById(merchantId)
	if err != nil {
		return nil, nil, err
	}
	userExists := false
	for _, m := range merchant.Users {
		if m.ID == user.ID {
			userExists = true
			break
		}
	}

	if userExists {
		return merchant, user, nil
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

func (db *MerchantRepoImpl) CreateUser(user entities.User, merchantId string) (*entities.User, error) {
	User, err := db.UserRepo.FindUserByPhone(user.PhoneNumber)
	if err == nil {
		_, _, err := db.AddUserToMerchant(merchantId, User.ID)
		if err != nil {
			return nil, err
		}
		return User, nil
	}

	err = db.Db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	_, _, err = db.AddUserToMerchant(merchantId, user.ID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *MerchantRepoImpl) PointCollection(userPhone string, point float64, merchantId string) (*entities.Wallet, error) {
	user, err := db.UserRepo.FindUserByPhone(userPhone)
	if err != nil {
		return nil, err
	}
	merchant, err := db.FindMerchantById(merchantId)
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

func (db *MerchantRepoImpl) GenerateKeyPair() (string, string, error) {

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

func (db *MerchantRepoImpl) CreateMerchant(merchant entities.Merchant) (*entities.Merchant, error) {
	err := db.Db.Create(&merchant).Error
	if err != nil {
		return nil, err
	}
	storedMerchant, err := db.FindMerchantByPhone(merchant.PhoneNumber)
	if err != nil {
		return nil, err
	}

	return storedMerchant, nil

}

func (db *MerchantRepoImpl) GetAllMerchants() (*[]entities.Merchant, error) {
	allMerchats := []entities.Merchant{}
	err := db.Db.Preload("Users").Find(&allMerchats).Error
	if err != nil {
		return nil, err
	}
	return &allMerchats, nil
}

func (db *MerchantRepoImpl) FindMerchantById(id string) (*entities.Merchant, error) {
	var merchant entities.Merchant

	err := db.Db.Preload("Users").Where("id=?", id).Take(&merchant).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil

}

func (db *MerchantRepoImpl) FindMerchantByPhone(phone string) (*entities.Merchant, error) {
	merchant := entities.Merchant{}

	err := db.Db.Preload("Users").Where("phone_number=?", phone).Take(&merchant).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil

}

func (db MerchantRepoImpl) RetrivePublicKey(phone string) (string, error) {
	merchant, err := db.FindMerchantByPhone(phone)
	if err != nil {
		return "", err
	}
	publicKey := merchant.PublicKey
	return publicKey, nil
}

func (db *MerchantRepoImpl) FindAllUsers(from string, to string, all bool, page int64, perpage int64) (*entities.GetAllUsers, error) {

	var users []entities.User
	var totalPage int64
	var prev int64
	var next int64

	if page <= 1 {
		prev = 0
	} else {

		prev = page - 1
	}
	next = page + 1

	sql := "SELECT * FROM users"

	db.Db.Model(entities.User{}).Count(&totalPage)

	if all {
		db.Db.Model(entities.User{}).Order("created_at DESC").Find(&users)
		return &entities.GetAllUsers{
			Users:    users,
			Total:    totalPage,
			LastPage: 0,
			Page:     0,
		}, nil
	}

	if from != "" && to != "" {
		sql = fmt.Sprintf("SELECT * FROM users WHERE created_at BETWEEN %s AND %s ", from, to)
		countTotal := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE created_at BETWEEN %s AND %s ", from, to)
		db.Db.Raw(countTotal).Scan(&totalPage)
	}
	db.Db.Raw(sql).Scan(&users)

	lastPage := math.Ceil(float64(totalPage / perpage))
	if page == int64(lastPage) {
		next = 0
	}
	return &entities.GetAllUsers{
		Users:    users,
		Total:    totalPage,
		Page:     page,
		Prev:     prev,
		Next:     next,
		LastPage: int64(lastPage),
	}, nil

}

func (db *MerchantRepoImpl) UpdateMerchant(merchant entities.Merchant) error {
	err := db.Db.Save(&merchant).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *MerchantRepoImpl) DeleteAll() error {
	// err := db.Db.Delete(&entities.Merchant{}).Error
	err := db.Db.Unscoped().Delete(&entities.Merchant{}).Error
	if err != nil {
		return err
	}
	return nil
}
