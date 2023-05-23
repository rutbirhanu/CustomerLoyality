package db_utils

import (
	"fmt"

	"github.com/santimpay/customer-loyality/configs"
	"github.com/santimpay/customer-loyality/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(config *configs.DbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s password=%s dbname=%s user=%s sslmode=disable",
		config.Host, config.Port, config.Password, config.DbName, config.Username)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&entities.Merchant{},
		&entities.User{},
		
	); err != nil {
		return nil, err
	}

	return db, nil
}
