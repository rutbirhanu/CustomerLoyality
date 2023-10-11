package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/entities"
	"github.com/santimpay/customer-loyality/internal/repositories"
	"gorm.io/gorm"
)

func GenerateNewToken(repo repositories.ApiRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		merchantId := c.Param("merchantid")

		token, err := repo.GenerateNewToken(merchantId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusAccepted, token)
		return nil
	}
}

func RemoveToken(repo repositories.ApiRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Param("token")
		err := repo.RemoveToken(token)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, "successfully removed")
	}
}

func GiveWallet(repo repositories.ApiRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		bearerToken := c.Request().Header.Get("Authorization")
		userPhone := c.QueryParam("userPhone")
		userName := c.QueryParam("userName")
		if bearerToken == "" || !strings.HasPrefix(bearerToken, "Bearer ") {
			// Handle missing or invalid token
			return c.JSON(http.StatusUnauthorized, "Invalid or missing token")
		}
		token := strings.TrimPrefix(bearerToken, "Bearer ")
		txHandler := c.Get("db_tx")
		trxRepo := repo.WithTrx(txHandler.(*gorm.DB))
		merchantId, err := trxRepo.FindMerchantFromToken(token)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		user, err := trxRepo.GiveWallet(userPhone, userName, merchantId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, user)
	}
}

func PointConfig(repo repositories.ApiRepo) echo.HandlerFunc {
	return func(c echo.Context) error {

		merchantid := c.Param("merchantid")
		config := entities.PointConfig{}
		err := c.Bind(&config)

		if err != nil {
			return c.JSON(http.StatusBadRequest, "error while parsing ")
		}
		err = repo.PointConfiguration(config.PointConfiguration, merchantid)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "can not set the point configuration")

		}
		return c.JSON(http.StatusOK, "point configuration is successfully set")

	}
}

// GivePoint(usersPhone string, amount float64, merchantId string) error

func GivePoint(repo repositories.ApiRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		userPhone := c.QueryParam("userPhone")
		amountStr := c.QueryParam("amount")
		txHandler := c.Get("db_tx")
		trxRepo := repo.WithTrx(txHandler.(*gorm.DB))
		merchantId, err := trxRepo.FindMerchantFromToken(token)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		err = trxRepo.GivePoint(userPhone, amount, merchantId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, "point given to user")

	}
}

// BuyGiftCard(merchantid string, amount float64, recipentPhone string, purchaserPhone string) error
func BuyGiftCard(repo repositories.ApiRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		amountStr := c.QueryParam("amount")
		recipentPhone := c.QueryParam("recipentPhone")
		purchaserPhone := c.QueryParam("purchaserPhone")
		txHandler := c.Get("db_tx")
		trxRepo := repo.WithTrx(txHandler.(*gorm.DB))
		bearerToken := c.Request().Header.Get("Authorization")
		if bearerToken == "" || !strings.HasPrefix(bearerToken, "Bearer ") {
			// Handle missing or invalid token
			return c.JSON(http.StatusUnauthorized, "Invalid or missing token")
		}
		token := strings.TrimPrefix(bearerToken, "Bearer ")
		merchantId, err := trxRepo.FindMerchantFromToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err)
		}
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		err = trxRepo.BuyGiftCard(merchantId, amount, recipentPhone, purchaserPhone)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusAccepted, "successfull gift card purchase")
	}
}

// redeemGiftCard(merchantId string, giftcardCode string, totalPrice float64) (float64,error)
