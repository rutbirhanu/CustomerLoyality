package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/service"
	"github.com/santimpay/customer-loyality/internal/util"
)

func Login(srv service.MerchantService) echo.HandlerFunc{
	return func (c echo.Context) error{
		data:= entities.MerchantLogin{}
		err:=c.Bind(&data)
		if err!=nil{
			return c.JSON(http.StatusBadRequest, "can not parse data")
		}
		user,exist:=srv.FindMerchantByPhone(data.PhoneNumber)
		if !exist{
			return c.JSON(http.StatusBadRequest, "incorrect input ")
		}
		passCheck:= util.VerifyPassword(data.Password,user.Password)
		if !passCheck{
			return c.JSON(http.StatusConflict, "incorrect password")
		}
		token,err:=util.GenerateToken(user.PhoneNumber , user.ID , user.MerchantName)
		if err!=nil{
			return c.JSON(http.StatusConflict,"can not create token")
		}
		user.Token=token
		data.Token=token
		c.JSON(http.StatusAccepted,data)
		return nil
	}
}


