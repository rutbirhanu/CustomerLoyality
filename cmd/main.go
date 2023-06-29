package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/configs"

	"github.com/santimpay/customer-loyality/internal/db_utils"
	// "github.com/santimpay/customer-loyality/internal/util"

	"github.com/santimpay/customer-loyality/internal/handlers"
	"github.com/santimpay/customer-loyality/internal/handlers/auth"
	"github.com/santimpay/customer-loyality/internal/middleware"
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

	userRepo := repositories.NewUserRepo(db)
	merchantRepo := repositories.NewMerchantRepo(db, userRepo)
	trxRepo := repositories.NewTransactionRepo(db, userRepo, merchantRepo)
	adminRepo := repositories.NewAdminRepo(db, userRepo, merchantRepo)
	userSrvc := service.NewUserSrvc(userRepo)
	merchantSrvc := service.NewMerchantSrvc(merchantRepo)

	app := echo.New()

	trxRoute := app.Group("/trx")
	trxRoute.Use(middleware.Auth(merchantRepo))
	trxRoute.Use(middleware.DBTransactionMiddlware(db))
	trxRoute.POST("/collect/:merchantid", handlers.PointCollection(trxRepo))
	trxRoute.POST("/transfer/:merchantid/:userid", handlers.TransferPoint(trxRepo))
	trxRoute.POST("/:merchantid/charity/:userid", handlers.Donate(trxRepo))
	trxRoute.GET("/find/:merchantid/:userid", handlers.FindUserMer(trxRepo))

	merchantRoute := app.Group("/merchant")
	merchantRoute.Use(middleware.Auth(merchantRepo))
	app.GET("/getMerchant/:merchantid", handlers.FindMerchantById(merchantSrvc))
	merchantRoute.POST("/addMerchant/:merchantid/:userid", handlers.Login(adminRepo, userSrvc, merchantRepo))
	app.POST("/signup", auth.Signup(merchantSrvc, merchantRepo))
	app.POST("/login", auth.Login(merchantSrvc, merchantRepo))
	merchantRoute.POST("/createUser/:merchantid", handlers.RegisterUser(merchantRepo))
	merchantRoute.GET("/allMerchant", handlers.GetAll(merchantSrvc))

	userRoute := app.Group("/user")
	userRoute.GET("/getUser/:userid", handlers.GetUserById(userSrvc))
	userRoute.POST("/userRegistration", handlers.UserRegistration(userRepo))

	adminRoute := app.Group("/admin")
	adminRoute.GET("/getWallet/:Walletid", handlers.GetWalletById(adminRepo))

	serverPort := os.Getenv("SERVER_PORT")
	app.Logger.Fatal(app.Start(fmt.Sprintf(":%s", serverPort)))

}
