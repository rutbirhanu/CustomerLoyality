package handlers

// import (
// 	"net/http"
// 	"github.com/labstack/echo/v4"
// 	"github.com/santimpay/customer-loyality/internal/repositories"
// 	// "gorm.io/gorm"
// )


// func DeleteAll(repo repositories.MerchantRepo) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		err := repo.DeleteAll()
// 		if err != nil {
// 			return c.JSON(http.StatusBadRequest, "no repo")
// 		}
// 		c.JSON(http.StatusAccepted, "nice")
// 		return nil
// 	}
// }
