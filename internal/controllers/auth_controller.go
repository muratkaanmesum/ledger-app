package controllers

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"ptm/internal/models"
	"ptm/internal/repositories"
	"ptm/internal/services"
	"ptm/internal/utils/jwt"
	"ptm/internal/utils/logger"
	"ptm/internal/utils/response"
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
	userService := services.NewUserService(repositories.NewUserRepository())

	logger.Logger.Info("Handling RegisterUser request")

	if err := c.Bind(&req); err != nil {
		logger.Logger.Error("Failed to bind request body", zap.Error(err))
		return response.BadRequest(c, "Error", err)
	}

	logger.Logger.Debug("Request body parsed successfully", zap.String("username", req.Username), zap.String("email", req.Email))

	if err := c.Validate(req); err != nil {
		logger.Logger.Warn("Validation failed", zap.Error(err))
		return response.BadRequest(c, "Validation error", err)
	}

	user, err := userService.RegisterUser(&models.User{
		Username:     req.Username,
		Email:        req.Email,
		Role:         req.Role,
		PasswordHash: req.Password,
	})
	if err != nil {
		logger.Logger.Error("Failed to register user", zap.String("username", req.Username), zap.Error(err))
		return response.BadRequest(c, "Error", err)
	}

	logger.Logger.Info("User registered successfully", zap.String("username", req.Username), zap.String("email", req.Email))
	balanceService := services.NewBalanceService(repositories.NewBalanceRepository())
	_, err = balanceService.CreateBalance(user)
	if err != nil {
		return response.BadRequest(c, "Error", err)
	}

	return response.Ok(c, "User created", user)
}

func AuthenticateUser(c echo.Context) error {
	var req AuthenticateUserRequest
	userService := services.NewUserService(repositories.NewUserRepository())

	if err := c.Bind(&req); err != nil {
		logger.Logger.Error("Failed to bind request body", zap.Error(err))
		return response.BadRequest(c, "Error", err)
	}

	logger.Logger.Debug("Request body parsed successfully", zap.String("username", req.Username))

	if err := c.Validate(req); err != nil {
		logger.Logger.Warn("Validation failed", zap.Error(err))
		return response.BadRequest(c, "Validation error", err)
	}

	user, err := userService.GetUserByUsername(req.Username)
	if err != nil {
		logger.Logger.Error("Failed to find the user by username", zap.String("username", req.Username), zap.Error(err))
		return response.BadRequest(c, "Error", err)
	}

	if err := user.VerifyUser(req.Password); err != nil {
		logger.Logger.Error("Failed to verify user", zap.String("username", req.Username), zap.Error(err))
		return response.BadRequest(c, "Authentication Error", err)
	}
	token, err := jwt.GenerateJWT(user)

	if err != nil {
		logger.Logger.Error("Failed to generate token", zap.String("username", req.Username), zap.Error(err))
		return response.InternalServerError(c, "Internal Server Error", err)
	}

	return response.Ok(c, "User authenticated", map[string]interface{}{
		"access_token": token,
		"user":         user,
	})
}
