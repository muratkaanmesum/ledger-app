package controllers

import (
	"fmt"
	"net/http"
	"ptm/repositories"
	"ptm/services"
	"ptm/utils/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TransactionRequest struct {
	SenderId   int `json:"sender_id" validate:"required"`
	ReceiverId int `json:"receiver_id" validate:"required"`
	Amount     int `json:"amount" validate:"required"`
}

func GetAllUsers(c echo.Context) error {
	userService := services.NewUserService(repositories.NewUserRepository())

	users, err := userService.GetAllUsers(10, 0)
	if err != nil {
		return response.BadRequest(c, "Error getting users", err)
	}
	return c.JSON(http.StatusOK, users)
}

func GetUserById(c echo.Context) error {
	userService := services.UserService{}
	idString := c.Param("id")
	num, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return response.BadRequest(c, "Error converting string to integer", err)
	}
	user, err := userService.GetUserById(num)

	if err != nil {
		return response.BadRequest(c, "User not found", err)
	}
	return c.JSON(http.StatusOK, user)
}
