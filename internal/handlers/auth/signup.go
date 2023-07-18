package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/service"
	"github.com/santimpay/customer-loyality/internal/repositories"
	"github.com/santimpay/customer-loyality/internal/util"
)

func Signup(srvc service.MerchantService, repo repositories.MerchantRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		merchant := entities.Merchant{}
		err := c.Bind(&merchant)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		_, exist := srvc.FindMerchantByPhone(merchant.PhoneNumber)
		if exist {
			return c.JSON(http.StatusBadRequest, "phone number already exist")

		}
		hashedPass := util.HashPassword(merchant.Password)
		privateKey,publicKey, err := repo.GenerateKeyPair()
		
		if err != nil {
			return err
		}


		merchantData := entities.Merchant{
			MerchantName: merchant.MerchantName,
			Password:     hashedPass,
			PhoneNumber:  merchant.PhoneNumber,
			BusinessName: merchant.BusinessName,
			PrivateKey: privateKey,
			PublicKey: publicKey,
		}
		data, stored := srvc.CreateMerchant(merchantData)
		if !stored {
			return c.JSON(http.StatusBadRequest, err)
		}
		
		response:= entities.CreatedMerchantResponse{
			MerchantName: data.MerchantName,
			Password:     hashedPass,
			PhoneNumber:  data.PhoneNumber,
			BusinessName: data.BusinessName,
		}
		c.JSON(http.StatusAccepted, response)
		return nil
	}
}
