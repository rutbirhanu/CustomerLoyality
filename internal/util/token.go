package util

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	Admin    Role = "admin"
	User     Role = "user"
	Merchant Role = "merchant"
)

type ClaimData struct {
	PhoneNumber string
	// Name        string
	UserId      string
	jwt.StandardClaims
	Role       Role
	MerchantID string `json:"merchantid,omitempty"`
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

func GenerateToken(phone string, uid string, privateKeyBytes []byte, role Role, merchantId string) (string, error) {
	tokenClaim := &ClaimData{
		PhoneNumber: phone,
		UserId:      uid,
		// Name:        name,
		Role:        role,
		MerchantID:  merchantId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(5)).Unix(),
		},
	}

	//  privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	// 	if err != nil {
	// 	    return "", err
	// 	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return "", errors.New("no PEM block found")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, tokenClaim).SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}
