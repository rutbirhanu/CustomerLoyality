// package middleware

// import (
// 	"errors"
// 	"fmt"
// 	"net/http"

// 	"github.com/golang-jwt/jwt"
// 	"github.com/labstack/echo/v4"
// 	"github.com/santimpay/customer-loyality/internal/repositories"
// )

// func Auth(next echo.HandlerFunc) echo.HandlerFunc {
// return func(c echo.Context ) error{
// 	tokenString,err := c.Cookie("auth-token")
// 	if err!=nil{
// 		return c.JSON(http.StatusNotFound, "cookie not found")
// 	} 
// 	var phone string
// 	token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
// 			return nil, fmt.Errorf("Unexpected signing method")
// 		}
// 		if phoneNum, ok := token.Claims.(jwt.MapClaims)["PhoneNumber"].(string); ok{
// 			phone=phoneNum}else{
// 			return nil, fmt.Errorf("can not find the phone number")
// 		}
// 		publicKey,err := repositories.MerchantRepo.RetrivePublicKey(phone)
// 		if err!=nil{
// 			return nil, fmt.Errorf("can not get public key")
// 		}
// 		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
// 		return publicKey, nil
// 	})
	
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		fmt.Println(claims["PhoneNumber"], claims["Name"])
// 	} else {
// 		fmt.Println(err)
// 	}
// 	if !token.Valid {
// 		return errors.New("invalid token")
// 	}
	
// 	next(c)
// 	return nil
// }
// }

