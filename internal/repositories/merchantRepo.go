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
	// GetMerchant(entities.Merchant) *entities.Merchant
	GenerateKeyPair() (string, string, error)
	UpdateMerchant(entities.Merchant) (error)
	GetAllMerchants()(*[]entities.Merchant,error)
	FindAllUsers(from string, to string, all bool, page int64, perpage int64) (*entities.GetAllUsers, error)
}

type MerchantRepoImpl struct {
	Db 				*gorm.DB
	userRepo		 UserRepo
}

func NewMerchantRepo(db *gorm.DB, UserRepo UserRepo) MerchantRepo {
	return &MerchantRepoImpl{
		Db: db,
		userRepo: UserRepo,
	}
}

// func (db MerchantRepoImpl) GetMerchant(merchant entities.Merchant) *entities.Merchant {
// 	return &entities.Merchant{
// 		PhoneNumber:  merchant.PhoneNumber,
// 		Password:     merchant.Password,
// 		MerchantName: merchant.MerchantName,
// 		BusinessName: merchant.BusinessName,
// 	}

// }

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

func (db *MerchantRepoImpl) GetAllMerchants()(*[]entities.Merchant,error){
	allMerchats:=[]entities.Merchant{}
	err:=db.Db.Find(&allMerchats).Error
	if err!=nil{
		return nil,err
	}
	return &allMerchats,nil
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


func (db *MerchantRepoImpl) AddUser(userId string, merchantId string) error {
	merchant, err := db.FindMerchantById(merchantId)
	if err != nil {
		return err
	}
	user,err:= db.userRepo.FindUserById(userId)

	if err != nil {
		return err
	}
	if err := db.Db.Model(&merchant).Association("Users").Append(user); err != nil {
		return err
	}
	if err := db.Db.Model(&user).Association("Users").Append(merchant); err != nil {
		return err
	}
	return nil
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


func (db *MerchantRepoImpl) UpdateMerchant(merchant entities.Merchant)(error){
	err := db.Db.Save(&merchant).Error
	if err!=nil{
		return err
	}
	return nil
}