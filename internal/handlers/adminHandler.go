package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
	// "gorm.io/gorm"
)

func PointCollection(repo repositories.AdminRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		merchantId := c.Param("merchantid")
		data := entities.Collection{}
		err := c.Bind(&data)
		// tx_handler:= c.Get("db_tx")
		if err != nil {
			return c.JSON(http.StatusBadRequest, "can not parse data")
		}
		// trxRepo:= repo.WithTrx(tx_handler.(*gorm.DB))
		wallet,err:= repo.PointCollection(data.UserPhone,data.Points,merchantId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusAccepted, wallet)

	}
}
