package repositories

import (
	"github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/gorm"
)

type TransactionRepo interface {
	// CreateEntry()
	// UpdateBalance()
	// CreateTransaction()
}

type TransactionRepoImpl struct{
	Db 		*gorm.DB
}

func NewTransactionRepo (db *gorm.DB) TransactionRepo{
	return &TransactionRepoImpl{
		Db: db,
	}
}

func (db *TransactionRepoImpl) PerformTransaction(amount float64, types string, action string, to string, userMerchantId string ) error {

	switch action {
	case "debit":
		amount = -amount
	case "credit":
		amount = +amount
	}

	UserWithMerchant:= entities.MerchantUsers{}
	err:= db.Db.First(&UserWithMerchant, userMerchantId).Error
	if err!=nil{
		return err
	}
	user :=entities.User{}

	err=db.Db.First(&user, UserWithMerchant.UserID).Error
	if err!=nil{
		return err
	}
	if user.Balance < amount {
		return nil
	}

	switch types{
	case "cashOut", "donate":
		toUser:= entities.User{}
		err= db.Db.First(&toUser).Error
		if err!=nil{
			return err
		}
		
		toUser.Balance+=amount
	}

	user.Balance += amount
	return nil

}
