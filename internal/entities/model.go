package entities

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        		string   			`gorm:"primary_key" json:"id"`
	CreatedAt 		string				`json:"created_at"`
	UpdatedAt 		string				`json:"updated_at"`
	DeletedAt 		gorm.DeletedAt	 	`gorm:"index" json:"deleted_at"`
}

func (model *Model) BeforeCreate(scope *gorm.DB) error {
	uuid,err := uuid.New().MarshalText()
	if err!=nil{
		return err
	}
	model.ID=string(uuid)
	model.CreatedAt = time.Now().Format(time.UnixDate)
	return nil

}

func (model *Model) BeforeUpdate(scope *gorm.DB) error{
	model.UpdatedAt = time.Now().Format(time.UnixDate)
	return nil
}

