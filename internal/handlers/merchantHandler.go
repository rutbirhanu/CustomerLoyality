package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
	"github.com/santimpay/customer-loyality/internal/service"
)

func GetAll(srvc service.MerchantService) echo.HandlerFunc {
	return func(c echo.Context) error {
		data, found := srvc.GetAllMerchants()
		if !found {
			return c.JSON(http.StatusNotFound, "cannot find ")
		}
		c.JSON(http.StatusOK, data)
		return nil

	}
}

func RegisterUser(repo repositories.MerchantRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		merchantID:=c.Param("merchantid")
		user := entities.User{}
		err := c.Bind(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "can not parse data")
		}
// check if the user is in the user table before generating the keys

		privateKey, publicKey, err := repo.GenerateKeyPair()

		if err != nil {
			return err
		}
		userData := entities.User{
			PhoneNumber: user.PhoneNumber,
			UserName:    user.UserName,
			PrivateKey:  privateKey,
			PublicKey:   publicKey,
		}
		
		User, err := repo.CreateUser(userData,merchantID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		c.JSON(http.StatusCreated, User)

		return nil

	}
}

func FindMerchantById(srvc service.MerchantService) echo.HandlerFunc {
	return func(c echo.Context) error {
		merchantId := c.Param("merchantid")
		data, found := srvc.FindMerchantById(merchantId)
		if !found {
			return c.JSON(http.StatusBadRequest, "not found")
		}
		c.JSON(http.StatusAccepted, data)
		return nil
	}
}
