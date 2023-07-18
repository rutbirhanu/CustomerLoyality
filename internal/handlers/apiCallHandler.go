package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
)

// RemoveToken(token string) error
// GiveWallet(phone string, username string, merchantId string) (*entities.User, error)
// GivePoint(usersPhone string, amount float64, merchantId string) error
// BuyGiftCard(merchantid string, amount float64, recipentPhone string, purchaserPhone string) error
// FindGiftCardByCode( giftcardCode string)(*entities.GiftCard,error)
// redeemGiftCard(merchantId string, giftcardCode string, totalPrice float64) (float64,error)

func GenerateNewToken(repo repositories.ApiRepo)echo.HandlerFunc{
	return func(c echo.Context)error{
		merchantId :=c.Param("merchantid")

		token,err:= repo.GenerateNewToken(merchantId)
		if err!=nil{
			return c.JSON(http.StatusBadRequest,"can not generate token")
		}
		c.JSON(http.StatusAccepted,token)
		return nil
	}
}


func RemoveToken(repo repositories.ApiRepo)echo.HandlerFunc{
	return func(c echo.Context) error{
		token:= c.QueryParam("token")
err:=repo.RemoveToken(token)
if err!=nil{
	return c.JSON(http.StatusBadRequest,err)
}
return c.JSON(http.StatusOK,"successfully removed")
	}
}


func GiveWallet(repo repositories.ApiRepo)echo.HandlerFunc{
	return func(c echo.Context)error{

	}
}

func PointConfig(repo repositories.ApiRepo) echo.HandlerFunc{
	return func(c echo.Context)error{

		merchantid:= c.Param("merchantid")
		config:=entities.PointConfig{}
		err:= c.Bind(&config)

		if err!=nil{
			return c.JSON(http.StatusBadRequest,"error while parsing ")
		}
		err=repo.PointConfiguration(config.PointConfiguration,merchantid)
		if err!=nil{
			return c.JSON(http.StatusBadRequest,"can not set the point configuration")

		}

	}
}
