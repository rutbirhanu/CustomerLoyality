package middleware

import (
	// "crypto/rsa"
	"crypto/ecdsa"
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
			merchantid := c.Param("merchantid")

			if err != nil {
				return c.JSON(http.StatusNotFound, "cookie not found")
			}
			var phone string
			// token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
			// 	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			// 		return nil, fmt.Errorf("unexpected signing method")
			// 	}
			// 	if phoneNum, ok := token.Claims.(jwt.MapClaims)["PhoneNumber"].(string); ok {
			// 		phone = phoneNum

			// 	} else {
			// 		return nil, fmt.Errorf("can not find the phone number")
			// 	}
			// 	publicKey, err := repo.RetrivePublicKey(phone)
			// 	if err != nil {

			// 		return false, err
			// 	}

			// 	block, _ := pem.Decode([]byte(publicKey))
			// 	if block == nil {
			// 		log.Fatalf("No PEM blob found")
			// 	}

			// 	Key, err := x509.ParsePKIXPublicKey(block.Bytes)

			// 	if err != nil {
			// 		log.Fatalf("Failed to parse public key: %v", err)
			// 		return "", err
			// 	}
			// 	// fmt.Print(publicKey)

			// 	return Key, nil
			// })

			token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				if phoneNum, ok := token.Claims.(jwt.MapClaims)["PhoneNumber"].(string); ok {
					phone = phoneNum
				} else {
					return nil, fmt.Errorf("cannot find the phone number")
				}
				publicKey, err := repo.RetrivePublicKey(phone)
				if err != nil {
					return nil, err
				}
				block, _ := pem.Decode([]byte(publicKey))
				if block == nil {
					log.Fatalf("No PEM blob found")
				}
				pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
				if err != nil {
					log.Fatalf("Failed to parse public key: %v", err)
				}
				ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
				if !ok {
					return nil, fmt.Errorf("unexpected public key type")
				}
				return ecdsaPubKey, nil
			})

			if err != nil {
				return errors.New("error occured")
			}

			if !token.Valid {
				return errors.New("invalid token")
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return errors.New("invalid token claims")
			}

			userRole, ok := claims["Role"]
			if !ok {
				return errors.New("role not found")

			}
			merchantID, ok := claims["merchantid"]
			if !ok {
				return errors.New("merchant ID not found in token claims")
			}

			if userRole != "merchant" {
				return errors.New("unauthorized access")
			}
			if merchantID != merchantid {
				return errors.New("login with this specific merchant")
			}

			next(c)
			return nil
		}
	}

}
