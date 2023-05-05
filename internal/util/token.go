package util

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type ClaimData struct {
	PhoneNumber string
	Name        string
	UserId      string
	jwt.StandardClaims
}

func HashPassword(password string) string {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(hashedPass)
}


func VerifyPassword(userPass string, providedPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(providedPass), []byte(userPass))
	if err != nil {
		return false
	}
	return true

}

func GenerateToken(phone string, uid string, name string, privateKeyBytes []byte) (string, error) {
 tokenClaim:= & ClaimData{
	PhoneNumber: phone,
	UserId: uid,
	Name: name,
	StandardClaims: jwt.StandardClaims{
		ExpiresAt:time.Now().Add(time.Hour * time.Duration(5)).Unix(),
	},
 }

 privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
	    return "", err
	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
// 	tokenString, err := token.SignedString(privateKey)
// 	if err != nil {
// 	    return "", err
// 	}

// 	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
// 	if err != nil {
// 	    return nil, err
// 	}

// 	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
// 	    return publicKey, nil
// 	})

// 	if err != nil {
// 	    return nil, err
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
	


 token, err:=jwt.NewWithClaims( jwt.SigningMethodRS256, tokenClaim).SignedString(privateKey)
 if err!=nil{
	return "", err
 }
 
 return token,nil
}