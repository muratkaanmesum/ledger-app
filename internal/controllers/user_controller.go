package controllers

import (
	"ptm/internal/di"
	"ptm/internal/services"
	"ptm/internal/utils/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserController interface {
	GetAllUsers(c echo.Context) error
	GetUserById(c echo.Context) error
}

type userController struct {
	userService services.UserService
}

type TransactionRequest struct {
	SenderId   int `json:"sender_id" validate:"required"`
	ReceiverId int `json:"receiver_id" validate:"required"`
	Amount     int `json:"amount" validate:"required"`
}

func NewUserController() UserController {
	service := di.Resolve[services.UserService]()

	return &userController{
		userService: service,
	}
}

func (uc *userController) GetAllUsers(c echo.Context) error {
	users, err := uc.userService.GetAllUsers(10, 0)
	if err != nil {
		return err
	}
	return response.Ok(c, "Successfully Fetched", users)
}

func (uc *userController) GetUserById(c echo.Context) error {
	idString := c.Param("id")
	num, err := strconv.Atoi(idString)
	if err != nil {
		return response.BadRequest(c, "Error converting string to integer", err)
	}
	user, err := uc.userService.GetUserById(uint(num))
	if err != nil {
		return err
	}
	return response.Ok(c, "User Found", user)
}
