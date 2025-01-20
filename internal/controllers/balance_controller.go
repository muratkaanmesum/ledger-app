package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"ptm/internal/di"
	"ptm/internal/services"
	"ptm/internal/utils/customError"
	"ptm/internal/utils/response"
	"ptm/pkg/jwt"
)

type BalanceController interface {
	GetBalance(c echo.Context) error
}

type balanceController struct {
	service services.BalanceService
}

type BalanceRequest struct {
	Date string `json:"date,omitempty"`
}

func NewBalanceController() BalanceController {
	service := di.Resolve[services.BalanceService]()

	return &balanceController{
		service: service,
	}
}

func (b *balanceController) GetBalance(c echo.Context) error {
	user := c.Get("user").(*jwt.CustomClaims)
	fmt.Println(user)
	balance, err := b.service.GetUserBalance(user.Id)

	if err != nil {
		return customError.New(customError.NotFound, err)
	}

	return response.Ok(c, "Fetched Successfully", balance)
}
