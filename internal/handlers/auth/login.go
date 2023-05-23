package auth

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
	"github.com/santimpay/customer-loyality/internal/service"
	"github.com/santimpay/customer-loyality/internal/util"
)

func Login(srv service.MerchantService, repo repositories.MerchantRepo) echo.HandlerFunc {
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
		token, err := util.GenerateToken(user.PhoneNumber, user.ID, user.MerchantName, []byte(private))
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
