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

func GenerateToken(phone string, uid string, name string) (string, error) {
 tokenClaim:= & ClaimData{
	PhoneNumber: phone,
	UserId: uid,
	Name: name,
	StandardClaims: jwt.StandardClaims{
		ExpiresAt:time.Now().Add(time.Hour * time.Duration(5)).Unix(),
	},
 }

 token, err:=jwt.NewWithClaims( jwt.SigningMethodHS256, tokenClaim).SignedString([]byte("hello"))
 if err!=nil{
	return "", err
 }
 
 return token,nil
}
