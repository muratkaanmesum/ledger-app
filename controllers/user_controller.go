package controllers

import (
	"fmt"
	"net/http"
	"ptm/models"
	"ptm/services"
	"ptm/utils/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role"`
}

type TransactionRequest struct {
	SenderId   int `json:"sender_id" validate:"required"`
	ReceiverId int `json:"receiver_id" validate:"required"`
	Amount     int `json:"amount" validate:"required"`
}

func RegisterUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(req); err != nil {
		return response.BadRequest(c, "Error", err)
	}

	if err := c.Validate(req); err != nil {
		return response.BadRequest(c, "Validation error", err)
	}

	user, err := services.RegisterUser(&models.User{
		Username:     req.Username,
		Email:        req.Email,
		Role:         req.Role,
		PasswordHash: req.Password,
	})
	if err != nil {
		return response.BadRequest(c, "Error", err)
	}

	return response.Ok(c, "User created", user)
}

func GetAllUsers(c echo.Context) error {
	users, err := services.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func GetUserById(c echo.Context) error {
	idString := c.Param("id")
	num, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return response.BadRequest(c, "Error converting string to integer", err)
	}
	user, err := services.GetUserById(num)

	if err != nil {
		return response.BadRequest(c, "User not found", err)
	}
	return c.JSON(http.StatusOK, user)
}
