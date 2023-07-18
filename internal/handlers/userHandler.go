package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
	"github.com/santimpay/customer-loyality/internal/service"
	"github.com/santimpay/customer-loyality/internal/util"
)

func UserRegistration(repo repositories.UserRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := entities.User{}
		err := c.Bind(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "can not parse data")
		}
		createdUser, err := repo.CreateUser(user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "can not create")
		}

		// test this here so that others can be improved
		// response:= entities.CreatedUserResponse{
		// 	PhoneNumber: createdUser.PhoneNumber,
		// 	UserName: createdUser.UserName,
		// 	Merchants: []*entities.UsersMerchantResponse{},
		// }
		c.JSON(http.StatusCreated, createdUser)
		return nil
	}
}

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

func GetWalletById(repo repositories.AdminRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("Walletid")
		Wallet, err := repo.FindWalletById(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusAccepted, Wallet)
	}
}

func Login(repo repositories.AdminRepo, srvc service.UserService, merRepo repositories.MerchantRepo) echo.HandlerFunc {
	return func(c echo.Context) error {

		merchantId := c.Param("merchantid")
		userId := c.Param("userid")

		mer, uss, err := repo.AddUserToMerchant(merchantId, userId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusAccepted, mer)
		c.JSON(http.StatusAccepted, uss)

		return nil
	}
}

func Loginn(repo repositories.UserRepo, srvc service.UserService, adminRepo repositories.AdminRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := entities.UserLogin{}
		err := c.Bind(&data)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "can not parse data")
		}
		merchantId := c.Param("merchantid")

		user, err := repo.FindUserByPhone(data.PhoneNumber)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "incorrect input ")
		}

		// private, _, err := repo.GenerateKeyPair()
		// if err != nil {
		// 	c.JSON(http.StatusBadGateway, err)
		// }
		private,err:= repo.RetrivePrivateKey(user.PhoneNumber)
		if err != nil {
			return c.JSON(http.StatusConflict, "can not retrive private key")
		}
		token, err := util.GenerateToken(user.PhoneNumber, user.ID, []byte(private), util.User, merchantId)
		if err != nil {
			return c.JSON(http.StatusConflict, "can not create token")
		}
		user.Token = token
		data.Token = token

		cookie := &http.Cookie{
			Name:  "auth-token",
			Value: token,
		}
		cookie.SameSite = http.SameSiteLaxMode
		cookie.HttpOnly = true
		c.SetCookie(cookie)

		c.JSON(http.StatusAccepted, private)
		c.JSON(http.StatusAccepted, data)
		return nil
	}
}
