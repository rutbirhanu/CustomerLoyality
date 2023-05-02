package repositories

import (
	"fmt"
	"math"

	"github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/gorm"
)

type MerchantRepo interface {
	CreateMerchant(entities.Merchant) (*entities.Merchant, error)
	FindMerchantById(string) (*entities.Merchant, error)
	FindMerchantByPhone(string) (*entities.Merchant, error)
	GetMerchant(entities.Merchant) *entities.Merchant


	CreateUser(entities.User)(*entities.User, error)
	FindUserById(string)(*entities.User,error)
	FindUserByPhone(string)(*entities.User, error)
	// UpdateUser(entities.Merchant) (*entities.User, error)
	FindAllUsers(from string, to string, all bool, page int, perpage int) (*entities.GetAllUsers,error)
}

type MerchantRepoImpl struct {
	Db *gorm.DB
}

func NewMerchantRepo(db *gorm.DB) MerchantRepo {
	return &MerchantRepoImpl{
		Db: db,
	}
}

func (db MerchantRepoImpl) GetMerchant(merchant entities.Merchant) *entities.Merchant {
	return &entities.Merchant{
		PhoneNumber:  merchant.PhoneNumber,
		Password:     merchant.Password,
		MerchantName: merchant.MerchantName,
		BusinessName: merchant.BusinessName,
	}

}

func (db *MerchantRepoImpl) CreateMerchant(merchant entities.Merchant) (*entities.Merchant, error) {
	err := db.Db.Create(&merchant).Error
	if err != nil {
		return nil, err
	}
	storedMerchant,err := db.FindMerchantByPhone(merchant.PhoneNumber)
	if err!=nil{
		return nil,err
	}

	regMerchant := db.GetMerchant(*storedMerchant)
	return regMerchant, nil

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


func (db *MerchantRepoImpl) CreateUser(user entities.User) (*entities.User,error){

	err:= db.Db.Create(&user).Error
	if err!=nil{
		return nil,err
	}

	foundUser, err:= db.FindUserByPhone(user.PhoneNumber)

	if err!=nil{
		return nil,err
	}
	return foundUser,nil

}


func (db *MerchantRepoImpl) FindUserByPhone (phone string) (*entities.User, error){
	user :=entities.User{}

	err:= db.Db.Where("phone_number=?",phone).Take(&user).Error
	if err!=nil{
		return nil,err
	}
	return &user,nil

}

func (db *MerchantRepoImpl) FindUserById(id string) (*entities.User, error){
	user:= entities.User{}
	err:=db.Db.Where("id=?",id).Take(&user).Error
	if err!=nil{
		return nil,err
	}
	return &user,nil
}

func (db *MerchantRepoImpl) FindAllUsers(from string, to string, all bool, page int64, perpage int64) (*entities.GetAllUsers, error){

	var users []entities.User
	var totalPage int64
	var prev int64
	var next int64

	if page<=1{
		prev=0
	}else{

		prev= page-1
	}
	next = page+1

	sql:= "SELECT * FROM users"

	db.Db.Model(entities.User{}).Count(&totalPage)

	if all{
		db.Db.Model(entities.User{}).Order("created_at DESC").Find(&users)
		return &users, nil
	}
	
	if from!="" && to!="" {
		sql= fmt.Sprintf("SELECT * FROM users WHERE created_at BETWEEN %s AND %s ", from,to)
		countTotal:= fmt.Sprintf("SELECT COUNT(*) FROM users WHERE created_at BETWEEN %s AND %s ", from,to)
		db.Db.Raw(countTotal).Scan(&totalPage)
	}
	
	lastPage := int64(math.Ceil(float64(totalPage/perpage)))
	if page == lastPage{
		next=0
	}
	return &entities.GetAllUsers{
		Users:users,
		Total:totalPage,
		PerPage:perpage,
		Prev:prev,
		Next:next,
		LastPage:lastPage,
	},nil

}
