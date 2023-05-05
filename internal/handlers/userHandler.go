package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
	"github.com/santimpay/customer-loyality/internal/service"
)

func RegisterUser(userSrvc service.UserService, repo repositories.UserRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		merchantID := c.Param("merchantid")
		user := entities.User{}
		err := c.Bind(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "can not parse data")
		}
		if merchantID == "" {
			return c.JSON(http.StatusBadRequest, "merchant ID is not provided")
		}

		userData, created := userSrvc.CreateUser(user, merchantID)
		if !created {
			return c.JSON(http.StatusBadRequest, "can not create user")
		}

		merch,err := repo.AddMerchant(merchantID,user.ID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusCreated, userData)
		c.JSON(http.StatusCreated,merch)
		return nil

	}
}
