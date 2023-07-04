package service

import (
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
)

type UserService interface {
	FindUserById(string) (*entities.User, bool)
	FindUserByPhone(string) (*entities.User, bool)
	// UserLogin(entities.UserLogin,string)(*entities.User,bool)
	GenerateKeyPair() (string, string, bool)
	RetrivePublicKey(phone string) (string, bool)
	RetrivePrivateKey(phone string) (string, bool)
	// CheckBalance()

}

type UserSrvcImpl struct {
	UserRepo repositories.UserRepo
}

func NewUserSrvc(userRepo repositories.UserRepo) UserService {
	return &UserSrvcImpl{
		UserRepo: userRepo,
	}
}

func (userSrvc *UserSrvcImpl) FindUserById(userId string) (*entities.User, bool) {
	user, err := userSrvc.UserRepo.FindUserById(userId)
	if err != nil {
		return nil, false
	}
	return user, true
}

func (userSrvc *UserSrvcImpl) FindUserByPhone(phone string) (*entities.User, bool) {
	user, err := userSrvc.UserRepo.FindUserByPhone(phone)
	if err != nil {
		return nil, false
	}
	return user, true
}

// func (userSrvc *UserSrvcImpl) UserLogin(user entities.UserLogin, merchantId string) (*entities.User,bool){
// 	merchant,err :=userSrvc.UserRepo.UserLogin(user,merchantId)

// 	if err!=nil{
// 		return nil,false
// 	}
// 	return merchant,true
// }

func (userSrvc *UserSrvcImpl) RetrivePublicKey(phone string) (string, bool) {
	key, err := userSrvc.UserRepo.RetrivePublicKey(phone)
	if err != nil {
		return "", false
	}
	return key, true
}

func (userSrvc *UserSrvcImpl) GenerateKeyPair() (string, string, bool) {
	private, public, err := userSrvc.UserRepo.GenerateKeyPair()
	if err != nil {
		return "", "", false
	}
	return private, public, true
}

func (userSrvc *UserSrvcImpl) RetrivePrivateKey(phone string) (string, bool) {
	privateKey, err := userSrvc.UserRepo.RetrivePrivateKey(phone)
	if err != nil {
		return "", false
	}
	return privateKey, true
}
