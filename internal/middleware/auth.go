package middleware

import (
	// "crypto/rsa"
	"crypto/x509"
	"log"

	// "encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/repositories"
)

type AuthRepository interface {
	RetrivePublicKey(phone string) (string, error)
}

func Auth(repo repositories.MerchantRepo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString, err := c.Cookie("auth-token")
			merchantid:= c.Get("merchantID")

			if err != nil {
				return c.JSON(http.StatusNotFound, "cookie not found")
			}
			var phone string
			token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				if phoneNum, ok := token.Claims.(jwt.MapClaims)["PhoneNumber"].(string); ok {
					phone = phoneNum

				} else {
					return nil, fmt.Errorf("can not find the phone number")
				}
				publicKey, err := repo.RetrivePublicKey(phone)
				fmt.Print(publicKey)
				if err != nil {

					return false, err
				}

				block, _ := pem.Decode([]byte(publicKey))
				if block == nil {
					log.Fatalf("No PEM blob found")
				}

				Key, err := x509.ParsePKIXPublicKey(block.Bytes)

				if err != nil {
					log.Fatalf("Failed to parse public key: %v", err)
					return "", err
				}
				// fmt.Print(publicKey)

				fmt.Print(Key)
				fmt.Print(merchantid)
				return Key, nil
			})

			if err != nil {
				fmt.Print(err)
				fmt.Print("errorrr")
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

			// userRole, ok := claims["Role"]
			// if !ok {
			// 	return errors.New("role not found")

			// }
			// merchantID, ok := claims["merchantid"]
			// if !ok {
			// 	return errors.New("merchant ID not found in token claims")
			// }

			// if userRole != "merchant" {
			// 	return errors.New("unauthorized access")
			// }
			// if merchantID != "29e8b28f-8125-4ae4-bc94-72e040cdfd99" {
			// 	return errors.New("login with this specific merchant")
			// }

			next(c)
			return nil
		}
	}

}
