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
	
	GenerateKeyPair() (string, string, error)
	CreateUser(entities.User) (*entities.User, error)
	FindUserById(string) (*entities.User, error)
	FindUserByPhone(string) (*entities.User, error)
	// UserLogin(user entities.UserLogin, merchantId string) (*entities.User, error)
	UpdateUser(user *entities.User) error
	RetrivePublicKey(phone string) (string, error)
	RetrivePrivateKey(phone string) (string, error)
	// CheckBalance()
}

type UserRepoImpl struct {
	Db           *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &UserRepoImpl{
		Db:           db,
	}
}


func (db *UserRepoImpl) CreateUser(User entities.User) (*entities.User, error){
	err := db.Db.Create(&User).Error
	if err != nil {
		return nil, err
	}
	return &User,nil

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


func (db *UserRepoImpl) UpdateUser(user *entities.User) error {
	err := db.Db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
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

func (db *UserRepoImpl) RetrivePublicKey(phone string) (string, error){
	user,err:= db.FindUserByPhone(phone)
	if err!=nil{
		return "", err
	}
	return user.PublicKey, nil
}

func (db UserRepoImpl) RetrivePrivateKey(phone string) (string, error) {
	user, err := db.FindUserByPhone(phone)
	if err != nil {
		return "", err
	}
	privateKey := user.PrivateKey
	return privateKey, nil
}
