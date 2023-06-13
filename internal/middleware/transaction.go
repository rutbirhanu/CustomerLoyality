package middleware

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func StatusList(status int, statusList []int) bool {

	for _, value := range statusList {
		if value == status {
			return true
		}
	}
	return false
}

func DBTransactionMiddlware(db *gorm.DB) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc{
		return func(c echo.Context) error{

			if db==nil{
				log.Fatalln("db not specified")
			}
			tx:= db.Begin()

			defer func(){
				if r:=recover(); r!=nil{
					tx.Rollback()
				}
			}()
			
			c.Set("db_tx",tx)
			result := next(c)

			if StatusList(c.Response().Status, []int{http.StatusAccepted, http.StatusOK, http.StatusCreated}){
				if err:= tx.Commit().Error;err!=nil{
					log.Print("error while commiting")
				}
				return result
			}else{
				tx.Rollback()
				log.Print("rollback due to the returned status code")
				return result
			}
		}

	}

}
