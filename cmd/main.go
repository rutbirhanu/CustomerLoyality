package main

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/configs"
	"github.com/santimpay/customer-loyality/internal/db_utils"

	"github.com/santimpay/customer-loyality/internal/handlers"
	"github.com/santimpay/customer-loyality/internal/handlers/auth"
	"github.com/santimpay/customer-loyality/internal/repositories"
	"github.com/santimpay/customer-loyality/internal/service"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	Config := configs.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		Username: os.Getenv("DB_USERNAME"),
		DbName:   os.Getenv("DB_NAME"),
	}
	db, err := db_utils.Connect(&Config)
	if err != nil {
		log.Fatal(err)
	}
	var userRepo repositories.UserRepo
	merchantRepo := repositories.NewMerchantRepo(db,userRepo)
	userRepo = repositories.NewUserRepo(db, merchantRepo)
	userSrvc:=service.NewUserSrvc(userRepo)
	merchantSrvc := service.NewMerchantSrvc(merchantRepo)

	app := echo.New()
	app.GET("/allMerchant", handlers.GetAll(merchantSrvc))
	app.GET("/getMerchant", handlers.FindMerchantById(merchantSrvc))
	app.POST("/addMerchant/:merchantid",handlers.AddMerchant(userSrvc,userRepo))
	app.POST("/signup", auth.Signup(merchantSrvc , merchantRepo))
	app.POST("/login", auth.Login(merchantSrvc,merchantRepo))
	app.POST("/createUser/:merchantid", handlers.RegisterUser(userSrvc, userRepo))
	serverPort := os.Getenv("SERVER_PORT")
	app.Logger.Fatal(app.Start(fmt.Sprintf(":%s", serverPort)))

}