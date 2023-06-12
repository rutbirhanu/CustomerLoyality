package middleware

import (
	"crypto/rsa"
	"crypto/x509"

	// "encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/repositories"
	// "github.com/santimpay/customer-loyality/internal/util"
)

type AuthRepository interface {
	RetrivePublicKey(phone string) (string, error)
}

func Auth(repo repositories.MerchantRepo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString, err := c.Cookie("auth-token")
			// merchantid := c.Get("merchantID")
			if err != nil {
				return c.JSON(http.StatusNotFound, "cookie not found")
			}
			var phone string
			token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					fmt.Print("1")
					return nil, fmt.Errorf("unexpected signing method")
				}
				if phoneNum, ok := token.Claims.(jwt.MapClaims)["PhoneNumber"].(string); ok {
					phone = phoneNum

				} else {
					fmt.Print(phoneNum)
					return nil, fmt.Errorf("can not find the phone number")
				}
				publicKey, err := repo.RetrivePublicKey(phone)
				if err != nil {
					fmt.Print("122")

					return false, err
				}
				publicKeyBlock, _ := pem.Decode([]byte(publicKey))
				if publicKeyBlock == nil || publicKeyBlock.Type != "PUBLIC KEY" {
					fmt.Print("11")
					return false, errors.New("invalid public key")
				}

				parsedKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
				if err != nil {
					fmt.Print("12")

					return false, err
				}

				Key, ok := parsedKey.(*rsa.PublicKey)
				if !ok {
					fmt.Print("13")

					return false, errors.New("failed to parse RSA public key")
				}

				fmt.Print(Key)
				return Key, nil
			})
			if err != nil {
				fmt.Print(err)
				fmt.Print("verr")
				return errors.New("error occured")
			}

			fmt.Print(token)

			if !token.Valid {
				return errors.New("invalid token")
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return errors.New("invalid token claims")
			}
			fmt.Print(claims)

			userRole, ok := claims["Role"]
			if !ok {
				return errors.New("role not found")

			}
			// fmt.Print(userRole)
			// merchantID, ok := claims["merchantid"]
			// if !ok {
			// 	return errors.New("merchant ID not found in token claims")
			// }
			// fmt.Print(merchantID)

			if userRole != "merchant" {
				return errors.New("unauthorized access")
			}
			// if merchantID != "29e8b28f-8125-4ae4-bc94-72e040cdfd99" {
			// 	return errors.New("login with this specific merchant")
			// }

			next(c)
			return nil
		}
	}

}
