package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
)

func TransferPoint(repo repositories.TransactionRepo) echo.HandlerFunc{
	return func(c echo.Context) error{

	userId := c.QueryParam("userid")
	merchantId := c.QueryParam("merchantid")

	// transaction:= entities.Transaction{}
	// tx_db :=c.Get("tx_db")

	toUser:=entities.TransferPoint{}
	err:=c.Bind(&toUser)
	if err!=nil{
		return err
	}
	err = repo.TransferPoints(toUser.Amount,userId,merchantId,toUser.Phone)
	if err!=nil{
		return err
	}


return nil
}
}

func BuyAirTime(repo repositories.TransactionRepo) echo.HandlerFunc{
	return func(c echo.Context) error{
		userId := c.QueryParam("userid")
	merchantId := c.QueryParam("merchantid")
		amount:= entities.RequestData{}
		err:= c.Bind(&amount)
		if err!=nil{
			return err
		}
	err=repo.BuyAirTime(amount.Amount,userId,merchantId)
	if err!=nil{
		return err
	}
	return nil
	}
}

func Donate(repo repositories.TransactionRepo) echo.HandlerFunc{
	return func(c echo.Context) error{

		charityId:= c.QueryParam("charityid")
		userId := c.QueryParam("userid")
	merchantId := c.QueryParam("merchantid")
	amount:= entities.RequestData{}
	err:= c.Bind(&amount)
	
		if err!=nil{
			return err
		}
		err= repo.Donate(charityId,userId,merchantId,amount.Amount)
		if err!=nil{
			return err
		}

		return nil
	}
}

