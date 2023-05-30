package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
)

// func TransferMoney(c echo.Context) error{
// 	transaction:= entities.Transaction{}
// 	tx_db :=c.Get("tx_db")

// }

func RewardController(c echo.Context) error {

	reward := entities.Collection{}
	err := c.Bind(&reward)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, reward)
}
