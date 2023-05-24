package repositories

import "github.com/santimpay/customer-loyality/internal/entities"

type AdminRepo interface {
	FindMerchantById(string) (*entities.Merchant, error)
	FindMerchantByPhone(string) (*entities.Merchant, error)
	UpdateMerchant(entities.Merchant) (*entities.Merchant, error)
	FindAllMerchant(from string, to string, perpage int, all bool, page int) (*entities.Merchant, error)
	DeleteMerchant(string) (*entities.Merchant,error)
	FindUserById(string) (*entities.User, error)
	FindUserByPhone(string) (*entities.User, error)
}