package handlers

import (
	"net/http"

	// "github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
	"github.com/santimpay/customer-loyality/internal/service"
)

func RegisterUser(userSrvc service.UserService, repo repositories.UserRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := entities.User{}
		err := c.Bind(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "can not parse data")
		}

		registered, _ := userSrvc.FindUserByPhone(user.PhoneNumber)
		if registered != nil {
			return c.JSON(http.StatusBadRequest, "already created phone")

		}
		// if !found {
		// 	return c.JSON(http.StatusBadRequest, "not found")
		// }

		userData, created := userSrvc.CreateUser(user)
		if !created {
			return c.JSON(http.StatusBadRequest, "can not create user")
		}

		c.JSON(http.StatusCreated, userData)

		return nil

	}
}

func Login(repo repositories.UserRepo, srvc service.UserService, merRepo repositories.MerchantRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		user:= entities.User{}
		err := c.Bind(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "can not parse")
		}
		merchantId := c.Param("merchantid")

		// Merchant, err := merRepo.FindMerchantById(merchantId)
		// if err != nil {
		// 	return c.JSON(http.StatusBadRequest, "merchant not found")
		// }

		// _, found := srvc.FindUserById(userId)
		// if found {
		// 	return c.JSON(http.StatusBadRequest, "user with phone number already exist")
		// }

		mer, uss, err := repo.AddMerchant(merchantId, user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusAccepted, mer)
		c.JSON(http.StatusAccepted, uss)

		return nil
	}
}

// func Login(repo repositories.UserMerchantRepo, usrepo repositories.UserRepo) echo.HandlerFunc {
// 	return func(c echo.Context) error {

// 		user := entities.UserLogin{}
// 		err:=c.Bind(&user)

// 		// validate := validator.New()
// 		// err = validate.Struct(user)
// 		// validationErrors := err.(validator.ValidationErrors)
// 		if err != nil {
// 			return c.JSON(http.StatusBadRequest, err)
// 		}

// 		merchantId := c.Param("merchantid")

// 		mer, uss, err := repo.AddMerchant(merchantId, user.PhoneNumber)
// 		if err != nil {
// 			return c.JSON(http.StatusBadRequest, err)
// 		}
// 		c.JSON(http.StatusAccepted, mer)
// 		c.JSON(http.StatusAccepted, uss)

// 		return nil
// 	}
// }

func GetUserById(srvc service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("userid")
		user, found := srvc.FindUserById(id)
		if !found {
			return c.JSON(http.StatusNotFound, " user not found")
		}
		c.JSON(http.StatusAccepted, user)
		return nil
	}
}
