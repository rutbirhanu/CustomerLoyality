package service

import (
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
)

type UserService interface {
	CreateUser(entities.User) (*entities.User, bool)
	FindUserById(string) (*entities.User, bool)
	FindUserByPhone(string) (*entities.User, bool)
	// AddMerchant(entities.Merchant, entities.User)(bool)
	UserLogin(entities.UserLogin,string)(*entities.User,bool)
	
}

type UserSrvcImpl struct {
	UserRepo 		repositories.UserRepo   
}


func NewUserSrvc (userRepo repositories.UserRepo) UserService{
	return &UserSrvcImpl{
		UserRepo: userRepo,
	}
}

func (userSrvc  *UserSrvcImpl) CreateUser(user entities.User) (*entities.User, bool){
	User, err := userSrvc.UserRepo.CreateUser(user)
	if err!= nil{
		return nil,false
	}
	
	return User,true

}


func (userSrvc  *UserSrvcImpl) FindUserById (userId string) (*entities.User, bool){
	user,err:=userSrvc.UserRepo.FindUserById(userId)
	if err!= nil{
		return nil,false
	}
	return user,true
}


func (userSrvc  *UserSrvcImpl) FindUserByPhone (phone string) (*entities.User, bool){
	user,err:=userSrvc.UserRepo.FindUserByPhone(phone)
	if err!=nil{
		return nil,false
	}
	return user,true
}


// func (userSrvc *UserSrvcImpl) AddMerchant(merchant entities.Merchant, user entities.User) (bool){
// 	err:=userSrvc.UserRepo.AddMerchant(merchant, user)
// 	if err!=nil{
// 		return false
// 	}
// 	return true
// }

func (userSrvc *UserSrvcImpl) UserLogin(user entities.UserLogin, merchantId string) (*entities.User,bool){
	merchant,err :=userSrvc.UserRepo.UserLogin(user,merchantId)

	if err!=nil{
		return nil,false
	}
	return merchant,true
}
