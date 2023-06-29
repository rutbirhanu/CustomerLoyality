package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
	"gorm.io/gorm"
)


func PointCollection(repo repositories.TransactionRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		merchantId := c.Param("merchantid")
		data := entities.Collection{}
		err := c.Bind(&data)
		tx_handler:= c.Get("db_tx")
		if err != nil {
			return c.JSON(http.StatusBadRequest, "can not parse data")
		}
		trxRepo:= repo.WithTrx(tx_handler.(*gorm.DB))
		wallet,err:= trxRepo.PointCollection(data.UserPhone,data.Points,merchantId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusAccepted, wallet)

	}
}


func TransferPoint(repo repositories.TransactionRepo) echo.HandlerFunc {
	return func(c echo.Context) error {

		userId := c.Param("userid")
		merchantId := c.Param("merchantid")
		tx_handler := c.Get("db_tx")
		txRepo := repo.WithTrx(tx_handler.(*gorm.DB))

		toUser := entities.TransferPoint{}
		err := c.Bind(&toUser)
		if err != nil {
			return err
		}
		wall, err := txRepo.TransferPoints(toUser.Amount, userId, merchantId, toUser.Phone)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		c.JSON(http.StatusAccepted, wall)
		return nil
	}
}

func BuyAirTime(repo repositories.TransactionRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userid")
		merchantId := c.Param("merchantid")
		amount := entities.RequestData{}
		tx_handler := c.Get("db_tx")
		txRepo := repo.WithTrx(tx_handler.(*gorm.DB))
		err := c.Bind(&amount)
		if err != nil {
			return err
		}
		err = txRepo.BuyAirTime(amount.Amount, userId, merchantId)
		if err != nil {
			return err
		}
		return nil
	}
}

func Donate(repo repositories.TransactionRepo) echo.HandlerFunc {
	return func(c echo.Context) error {

		charityId := c.QueryParam("charityid")
		userId := c.Param("userid")
		merchantId := c.Param("merchantid")
		tx_handler := c.Get("db_tx")
		txRepo := repo.WithTrx(tx_handler.(*gorm.DB))
		amount := entities.RequestData{}
		err := c.Bind(&amount)

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = txRepo.Donate(charityId, userId, merchantId, amount.Amount)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())

		}
	

		return c.JSON(http.StatusAccepted,"success")
	}
}

func FindUserMer(repo repositories.TransactionRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userid")
		merchantId := c.Param("merchantid")
		merchant, err := repo.FindSingleWallet(userId, merchantId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "user nto foudn")
		}
		return c.JSON(http.StatusAccepted, merchant)
	}
}
