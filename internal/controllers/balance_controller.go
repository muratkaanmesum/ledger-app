package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"ptm/internal/di"
	"ptm/internal/services"
	"ptm/pkg/jwt"
	"ptm/pkg/utils/response"
	"time"
)

type BalanceController interface {
	GetBalance(c echo.Context) error
	BalanceAtTime(c echo.Context) error
	GetHistoricalBalance(c echo.Context) error
}

type balanceController struct {
	service services.BalanceService
}

type BalanceRequest struct {
	Date time.Time `json:"date" validate:"required"`
}

func NewBalanceController() BalanceController {
	service := di.Resolve[services.BalanceService]()

	return &balanceController{
		service: service,
	}
}

func (b *balanceController) GetBalance(c echo.Context) error {
	user := jwt.GetUser(c)
	balance, err := b.service.GetUserBalance(user.Id)

	if err != nil {
		return err
	}
	return response.Ok(c, "Fetched Successfully", balance)
}

func (b *balanceController) BalanceAtTime(c echo.Context) error {
	var request BalanceRequest
	if err := c.Bind(&request); err != nil {
		return response.InternalServerError(c, "Failed to bind request body")
	}
	fmt.Println("Request", request)
	if err := c.Validate(&request); err != nil {
		return response.UnprocessableEntity(c, "Invalid request")
	}

	user := jwt.GetUser(c)
	history, err := b.service.GetBalanceAtTime(user.Id, request.Date)
	if err != nil {
		return err
	}

	return response.Ok(c, "History Fetched Successfully", history)
}

func (b *balanceController) GetHistoricalBalance(c echo.Context) error {
	user := jwt.GetUser(c)

	histories, err := b.service.GetUserBalanceHistory(user.Id)

	if err != nil {
		return err
	}

	return response.Ok(c, "Balance History Fetched Successfully", histories)
}
