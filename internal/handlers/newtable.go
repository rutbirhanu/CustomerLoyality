package handlers

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
)

func JoinTable(repo repositories.UserMerchantRepo) echo.HandlerFunc {
	return func (c  echo.Context) error{

	new:=entities.UserMerchant{}
	err:=c.Bind(&new)

	if err!=nil{
		return c.JSON(http.StatusBadRequest,err)
	}
	
	// response, err :=repo.CreateUserMerchant(new)
	// if err!=nil{
	// 	return c.JSON(http.StatusBadRequest,err)
	// }
	// c.JSON(http.StatusAccepted,response)
	return nil
		}
}