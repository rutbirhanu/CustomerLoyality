package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/service"

	// "github.com/santimpay/customer-loyality/internal/service"
	"github.com/santimpay/customer-loyality/internal/util"
)

func Signup(srvc service.MerchantService) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := entities.Merchant{}
		err := c.Bind(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		_, exist := srvc.FindMerchantByPhone(user.PhoneNumber)
		if exist {
			return c.JSON(http.StatusBadRequest, "phone number already exist")

		}
		hashedPass := util.HashPassword(user.Password)

		userData := entities.Merchant{
			MerchantName: user.MerchantName,
			Password:     hashedPass,
			PhoneNumber:  user.PhoneNumber,
			BusinessName: user.BusinessName,
		}
		data, stored := srvc.CreateMerchant(userData)
		if !stored {
			return c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusAccepted, data)
		return nil
	}
}
