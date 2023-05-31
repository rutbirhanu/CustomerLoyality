package handlers

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
	"github.com/santimpay/customer-loyality/internal/service"
	"github.com/santimpay/customer-loyality/internal/util"
)

func RegisterUser(userSrvc service.UserService, repo repositories.UserRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := entities.User{}
		err := c.Bind(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "can not parse data")
		}

		privateKey,publicKey, err := repo.GenerateKeyPair()
		
		if err != nil {
			return err
		}

		userData := entities.User{
			PhoneNumber: user.PhoneNumber,
			UserName: user.UserName,
			PrivateKey: privateKey,
			PublicKey: publicKey,
		}

		User, err := repo.CreateUser(userData)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		c.JSON(http.StatusCreated, User)

		return nil

	}
}


func Login(repo repositories.UserRepo, srvc service.UserService, merRepo repositories.MerchantRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		// user := entities.User{}
		// err := c.Bind(&user)
		// if err != nil {
		// 	return c.JSON(http.StatusBadRequest, "can not parse")
		// }
		merchantId := c.Param("merchantid")
		userId := c.Param("userid")

		// Merchant, err := merRepo.FindMerchantById(merchantId)
		// if err != nil {
		// 	return c.JSON(http.StatusBadRequest, "merchant not found")
		// }

		// _, found := srvc.FindUserById(userId)
		// if found {
		// 	return c.JSON(http.StatusBadRequest, "user with phone number already exist")
		// }

		mer, uss, err := repo.AddMerchant(merchantId, userId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusAccepted, mer)
		c.JSON(http.StatusAccepted, uss)

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


func GetWalletById(repo repositories.WalletRepo) echo.HandlerFunc{
	return func(c echo.Context) error{
		id:= c.Param("Walletid")
		Wallet, err:= repo.FindWalletById(id)
		if err!=nil{
			return c.JSON(http.StatusBadRequest, err)
		} 
		return c.JSON(http.StatusAccepted, Wallet)
	}
}


func Loginn(srv service.MerchantService, repo repositories.MerchantRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := entities.MerchantLogin{}
		err := c.Bind(&data)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "can not parse data")
		}
		user, exist := srv.FindMerchantByPhone(data.PhoneNumber)
		if !exist {
			return c.JSON(http.StatusBadRequest, "incorrect input ")
		}
		passCheck := util.VerifyPassword(data.Password, user.Password)
		if !passCheck {
			return c.JSON(http.StatusConflict, "incorrect password")
		}
		private , _ ,err:=repo.GenerateKeyPair()
		if err!=nil{
			c.JSON(http.StatusBadGateway,err)
		}
		token, err := util.GenerateToken(user.PhoneNumber, user.ID, user.MerchantName, []byte(private), string(util.Merchant))
		if err != nil {
			return c.JSON(http.StatusConflict, "can not create token")
		}
		user.Token = token
		data.Token = token

		cookie:= &http.Cookie{
			Name: "auth-token",
			Value: token,
		}
		cookie.SameSite=http.SameSiteLaxMode
		cookie.HttpOnly=true
		c.SetCookie(cookie)

		
		// c.JSON(http.StatusAccepted, privateKey)
		c.JSON(http.StatusAccepted,data)
		return nil
	}
}