package handlers

import (
	"fmt"
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
		tx_handler := c.Get("db_tx")
		if err != nil {
			return c.JSON(http.StatusBadRequest, "can not parse data")
		}
		trxRepo := repo.WithTrx(tx_handler.(*gorm.DB))
		wallet, err := trxRepo.PointCollection(data.UserPhone, data.Points, merchantId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		merchant, err := repo.FindMerchantFromWallet(wallet.MerchantID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err, statusCode := repo.SendSMS(fmt.Sprintf("you have earned %.1f points from %s ,\n thanks for choosing us", data.Points, merchant.BusinessName), "####", "0968######")
		if err != nil {
			 c.JSON(http.StatusBadRequest, err.Error())
		}
		if statusCode == 202 {
			c.JSON(http.StatusAccepted, "SMS succesfully sent")
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

		data := entities.TransferPoint{}
		err := c.Bind(&data)
		if err != nil {
			return err
		}
		wallet, err := txRepo.TransferPoints(data.Amount, userId, merchantId, data.Phone)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		merchant, err := repo.FindMerchantFromWallet(wallet.MerchantID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		user, err := repo.FindUserFromWallet(wallet.UserID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err, statusCode := repo.SendSMS(fmt.Sprintf("you have earned %.1f points from %s ,\n into your %s wallet", data.Amount, user.PhoneNumber, merchant.BusinessName), "####", "0968######")
		if err != nil {
			 c.JSON(http.StatusBadRequest, err.Error())
		}
		if statusCode == 202 {
			c.JSON(http.StatusAccepted, "SMS succesfully sent")
		}
		return c.JSON(http.StatusAccepted, wallet)
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

		return c.JSON(http.StatusAccepted, "success")
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

// func Send(repo repositories.TransactionRepo) echo.HandlerFunc{
// 	return func(c echo.Context) error{
// 		err:= repo.SendSMS("u have successfully implemented the api","9360","0968581847")
// 		if err!=nil{
// 			return c.JSON(http.StatusBadRequest,err)
// 		}
// 		return c.JSON(http.StatusAccepted,"sent")
// 	}
// }
