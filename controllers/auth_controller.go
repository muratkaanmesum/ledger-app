package controllers

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"ptm/models"
	"ptm/services"
	"ptm/utils"
	"ptm/utils/response"
)

type createUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role"`
}

type AuthenticateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func RegisterUser(c echo.Context) error {
	var req createUserRequest

	utils.Logger.Info("Handling RegisterUser request")

	if err := c.Bind(&req); err != nil {
		utils.Logger.Error("Failed to bind request body", zap.Error(err))
		return response.BadRequest(c, "Error", err)
	}

	utils.Logger.Debug("Request body parsed successfully", zap.String("username", req.Username), zap.String("email", req.Email))

	if err := c.Validate(req); err != nil {
		utils.Logger.Warn("Validation failed", zap.Error(err))
		return response.BadRequest(c, "Validation error", err)
	}

	user, err := services.RegisterUser(&models.User{
		Username:     req.Username,
		Email:        req.Email,
		Role:         req.Role,
		PasswordHash: req.Password,
	})
	if err != nil {
		utils.Logger.Error("Failed to register user", zap.String("username", req.Username), zap.Error(err))
		return response.BadRequest(c, "Error", err)
	}

	utils.Logger.Info("User registered successfully", zap.String("username", req.Username), zap.String("email", req.Email))

	return response.Ok(c, "User created", user)
}

func AuthenticateUser(c echo.Context) error {
	var req AuthenticateUserRequest

	if err := c.Bind(&req); err != nil {
		utils.Logger.Error("Failed to bind request body", zap.Error(err))
		return response.BadRequest(c, "Error", err)
	}

	utils.Logger.Debug("Request body parsed successfully", zap.String("username", req.Username))

	if err := c.Validate(req); err != nil {
		utils.Logger.Warn("Validation failed", zap.Error(err))
		return response.BadRequest(c, "Validation error", err)
	}

	user, err := services.GetUserByUsername(req.Username)
	if err != nil {
		utils.Logger.Error("Failed to find the user by username", zap.String("username", req.Username), zap.Error(err))
		return response.BadRequest(c, "Error", err)
	}

	if err := user.VerifyUser(req.Password); err != nil {
		utils.Logger.Error("Failed to verify user", zap.String("username", req.Username), zap.Error(err))
		return response.BadRequest(c, "Authentication Error", err)
	}

	return response.Ok(c, "User authenticated", user)
}
