package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
)

// func TransferMoney(c echo.Context) error{
// 	transaction:= entities.Transaction{}
// 	tx_db :=c.Get("tx_db")

// }


func RewardController(c echo.Context) error{

	reward := entities.Reward{}
	fmt.Print(reward)
	return nil
}