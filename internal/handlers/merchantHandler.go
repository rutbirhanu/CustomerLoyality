package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/santimpay/customer-loyality/internal/service"
)

func GetAll(srvc service.MerchantService)echo.HandlerFunc{
	return func(c echo.Context)error{
		data,found:=srvc.GetAllMerchants()
		if !found{
			return c.JSON(http.StatusNotFound, "cannot find ")
		}
		c.JSON(http.StatusOK,data)
		return nil

	}
}


func FindMerchantById(srvc service.MerchantService) echo.HandlerFunc{
	return func (c echo.Context) error{
		merchantId:=c.Param("merchantid")
		data,found:=srvc.FindMerchantById(merchantId)
		if !found{
			return c.JSON(http.StatusBadRequest, "not found")
		}
		c.JSON(http.StatusAccepted,data)
		return nil
	}
}