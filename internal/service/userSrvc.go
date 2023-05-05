package service

import (
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
)

type UserService interface {
	CreateUser(entities.User, string) (*entities.User, bool)
	FindUserById(string) (*entities.User, bool)
	FindUserByPhone(string) (*entities.User, bool)
}

type UserSrvcImpl struct {
	UserRepo 		repositories.UserRepo   
}


func NewUserSrvc (userRepo repositories.UserRepo) UserService{
	return &UserSrvcImpl{
		UserRepo: userRepo,
	}
}

func (userSrvc  *UserSrvcImpl) CreateUser(user entities.User, merchantId string) (*entities.User, bool){
	User, err := userSrvc.UserRepo.CreateUser(user,merchantId)
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

