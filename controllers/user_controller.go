package controllers

import (
	"fmt"
	"net/http"
	"ptm/services"
	"ptm/utils/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CreateUserRequest struct {
	Name string `json:"name" validate:"required"`
	Role string `json:"role"`
}

type TransactionRequest struct {
	SenderId   int `json:"sender_id" validate:"required"`
	ReceiverId int `json:"receiver_id" validate:"required"`
	Amount     int `json:"amount" validate:"required"`
}

func CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := services.CreateUser(req.Name, req.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Return the created user
	return c.JSON(http.StatusCreated, user)
}

func GetAllUsers(c echo.Context) error {
	users, err := services.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func SendToUser(c echo.Context) error {
	var req TransactionRequest
	if err := c.Bind(&req); err != nil {
		return response.InternalServerError(c, "Error while binding", err)
	}
	if err := c.Validate(req); err != nil {
		return response.BadRequest(c, "validation Error", err)
	}
	if err := services.Send(req.SenderId, req.ReceiverId, req.Amount); err != nil {
		return response.BadRequest(c, "Bad Request", err)
	}
	return c.JSON(http.StatusCreated, nil)
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
