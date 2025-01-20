package controllers

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/di"
	"ptm/internal/services"
	"ptm/internal/utils/customError"
	"ptm/internal/utils/response"
)

type BalanceController interface {
	GetBalance(c echo.Context) error
}

type balanceController struct {
	service services.BalanceService
}

type BalanceRequest struct {
	Id uint `json:"id"`
}

func NewBalanceController() BalanceController {
	service := di.Resolve[services.BalanceService]()

	return &balanceController{
		service: service,
	}
}

func (b *balanceController) GetBalance(c echo.Context) error {
	var request BalanceRequest
	if err := c.Bind(&request); err != nil {
		return customError.New(customError.BadRequest, err)
	}

	if err := c.Validate(request); err != nil {
		return customError.New(customError.BadRequest, err)
	}

	balance, err := b.service.GetUserBalance(request.Id)

	if err != nil {
		return customError.New(customError.NotFound, err)
	}

	return response.Ok(c, "Fetched Successfully", balance)
}
