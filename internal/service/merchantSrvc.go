package service

import (
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
)


type MerchantService interface {
	CreateMerchant(entities.Merchant) (*entities.Merchant, bool)
	FindMerchantById(string) (*entities.Merchant, bool)
	FindMerchantByPhone(string) (*entities.Merchant, bool)
	GetAllMerchants()(*[]entities.Merchant,bool)
	
}


type MerchantServiceImpl struct{
	merchantRepo     repositories.MerchantRepo
}

func NewMerchantSrvc(repo repositories.MerchantRepo) MerchantService{
	return &MerchantServiceImpl{
		merchantRepo:repo ,
		
	}
}


func (srvc MerchantServiceImpl) CreateMerchant(merchant entities.Merchant) (*entities.Merchant,bool){
	regMerchant,err:=srvc.merchantRepo.CreateMerchant(merchant)
	if err!=nil{
		return nil,false
	}
	return regMerchant,true
}

func (srvc MerchantServiceImpl) FindMerchantById(merchantId string) (*entities.Merchant, bool){
	merchant, err := srvc.merchantRepo.FindMerchantById(merchantId)
	if err!=nil{
		return nil, false
	}

	return merchant,true
}

func (srvc MerchantServiceImpl) FindMerchantByPhone (phone string) (*entities.Merchant, bool){
	merchant, err:= srvc.merchantRepo.FindMerchantByPhone(phone)
	if err!=nil{
		return nil,false
	}
	return merchant,true
}

func (srvc MerchantServiceImpl) GetAllMerchants()(*[]entities.Merchant,bool){
	data,err:=srvc.merchantRepo.GetAllMerchants()
	if err!=nil{
		return nil,false
	}
	return data, true
}
