package middleware

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/repositories"
	"github.com/santimpay/customer-loyality/internal/util"
)

func ApiCallAuth(repo repositories.MerchantRepo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenParam :=c.Request().Header.Get("Authorization")
			var phone string
			token, err := jwt.Parse(tokenParam, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodECDSA)
				if !ok {
					return "", c.JSON(http.StatusBadRequest, "signing method do not match ")
				}
				if phoneNum, ok := token.Claims.(jwt.MapClaims)["Phone"].(string); ok {
					phone = phoneNum
				} else {
					return "", c.JSON(http.StatusBadRequest, "can not obtain user info")
				}

				publicKey, err := repo.RetrivePublicKey(phone)
				if err != nil {
					return "", c.JSON(http.StatusBadRequest, "can not retriev key")
				}

				block, _ := pem.Decode([]byte(publicKey))
				if block == nil {
					return "", c.JSON(http.StatusBadRequest, "key block not found")
				}
				key, err := x509.ParsePKIXPublicKey(block.Bytes)
				if err != nil {
					return "", c.JSON(http.StatusBadRequest, "can not parse key")
				}
				PubKey, ok := key.(*ecdsa.PublicKey)
				if !ok {
					return "", c.JSON(http.StatusBadRequest, "can not parse key")
				}
				return PubKey, nil

			})
			if err != nil {
				return c.JSON(http.StatusUnauthorized, err)
			}
			if !token.Valid {
				return c.JSON(http.StatusUnauthorized, "invalid token")
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusNotFound, "can not resolve the claims")
			}
			role, ok := claims["Role"]
			if !ok {
				return c.JSON(http.StatusNotFound, "can not resolve the role")
			}
			if role != util.Merchant {
				return c.JSON(http.StatusUnauthorized, "you are not authorized")
			}

			next(c)
			return nil
		}
	}
}
