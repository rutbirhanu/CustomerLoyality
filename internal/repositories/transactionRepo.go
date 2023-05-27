package repositories

import (
	"log"

	"github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/gorm"
)

type TransactionRepo interface {
	// CreateEntry()
	// UpdateBalance()
	// CreateTransaction() error
}

type TransactionRepoImpl struct {
	Db *gorm.DB
	UserRepo 	UserRepo
}

func NewTransactionRepo(db *gorm.DB, userRepo UserRepo) TransactionRepo {
	return &TransactionRepoImpl{
		Db: db,
		UserRepo: userRepo,
	}
}

func (db *TransactionRepoImpl) WithTrx(trxDB *gorm.DB) TransactionRepo {
	if trxDB == nil {
		log.Print("txn db not founc")
	}
	db.Db = trxDB
	return db
	
}

func (db *TransactionRepoImpl) PerformTransaction(amount float64, types string, action string, to string, userMerchantId string) error {

	switch action {
	case "debit":
		amount = -amount
	case "credit":
		amount = +amount
	}

	UserWithMerchant := entities.MerchantUser{}
	err := db.Db.First(&UserWithMerchant, userMerchantId).Error
	if err != nil {
		return err
	}
	user := entities.User{}

	err = db.Db.First(&user, UserWithMerchant.UserID).Error
	if err != nil {
		return err
	}
	if user.Balance < amount {
		return nil
	}

	switch types {
	case "cashOut", "donate":
		toUser := entities.User{}
		err = db.Db.First(&toUser).Error
		if err != nil {
			return err
		}

		toUser.Balance -= amount
	}

	user.Balance += amount

	return nil

}

// func (db *TransactionRepoImpl) CreateTransaction() error {
// 	transaction := entities.Transaction{}
// 	err := db.PerformTransaction(transaction.Amount, transaction.Action, transaction.Type, transaction.To, transaction.UserMerchantID)
// 	if err != nil {
// 		return nil
// 	}
// 	err = db.Db.Create(&transaction).Error
// 	if err != nil {
// 		return nil
// 	}
// 	return nil
// }


func (db * TransactionRepoImpl) GiveReward(amount float64, reason string, userMerchantId string) (*entities.Reward,error){

	point :=  (amount/10)
	
	reward := entities.Reward{
		Points: point,
		Reason: reason,
		UserMerchantId: userMerchantId,
	}
	
	UserWithMerchant := entities.MerchantUser{}
	err := db.Db.First(&UserWithMerchant, userMerchantId).Error
	if err != nil {
		return nil,err
	}

	user := entities.User{}

	err = db.Db.First(&user, UserWithMerchant.UserID).Error
	if err != nil {
		return nil,err
	}
	user.Balance+= reward.Points
	err = db.UserRepo.UpdateUser(&user)
	if err!=nil{
		return nil,err
	}
return &reward,nil

}
