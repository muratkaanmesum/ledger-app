package controllers

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/di"
	"ptm/internal/services"
	"ptm/internal/utils/customError"
	"ptm/internal/utils/response"
	"ptm/pkg/jwt"
	"time"
)

type BalanceController interface {
	GetBalance(c echo.Context) error
	BalanceAtTime(c echo.Context) error
}

type balanceController struct {
	service services.BalanceService
}

type BalanceRequest struct {
	Date time.Time `json:"date,omitempty"`
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
		return customError.New(customError.NotFound, err)
	}

	return response.Ok(c, "Fetched Successfully", balance)
}

func (b *balanceController) BalanceAtTime(c echo.Context) error {
	var request BalanceRequest
	if err := c.Bind(&request); err != nil {
		return customError.New(customError.InternalServerError, err)
	}

	if err := c.Validate(request); err != nil {
		return customError.New(customError.BadRequest)
	}
	user := jwt.GetUser(c)

	history, err := b.service.GetBalanceAtTime(user.Id, request.Date)

	if err != nil {
		return customError.New(customError.InternalServerError, err)
	}

	return response.Ok(c, "History Fetched Successfully", history)
}
