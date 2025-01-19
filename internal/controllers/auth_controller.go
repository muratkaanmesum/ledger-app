package controllers

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"ptm/internal/di"
	"ptm/internal/models"
	"ptm/internal/services"
	"ptm/internal/utils/customError"
	"ptm/internal/utils/response"
	"ptm/pkg/jwt"
	"ptm/pkg/logger"
)

type AuthController interface {
	RegisterUser(c echo.Context) error
	AuthenticateUser(c echo.Context) error
}

type authController struct {
	userService services.UserService
}

func NewAuthController() AuthController {
	service, ok := di.Resolve((*services.UserService)(nil)).(services.UserService)

	if !ok || service == nil {
		panic("Failed to resolve UserService")
	}
	return &authController{
		userService: service,
	}
}

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

func (ac *authController) RegisterUser(c echo.Context) error {
	var req createUserRequest
	logger.Logger.Info("Handling RegisterUser request")

	if err := c.Bind(&req); err != nil {
		return customError.New(customError.BadRequest, err)
	}

	logger.Logger.Debug("Request body parsed successfully", zap.String("username", req.Username), zap.String("email", req.Email))

	if err := c.Validate(req); err != nil {
		logger.Logger.Warn("Validation failed", zap.Error(err))
		return customError.New(customError.BadRequest, err)
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		Role:         req.Role,
		PasswordHash: req.Password,
	}

	createdUser, err := ac.userService.RegisterUser(user)
	if err != nil {
		return customError.New(customError.InternalServerError, err)
	}

	logger.Logger.Info("User registered successfully", zap.String("username", createdUser.Username))
	return response.Ok(c, "User created successfully", createdUser)
}

func (ac *authController) AuthenticateUser(c echo.Context) error {
	var req AuthenticateUserRequest

	if err := c.Bind(&req); err != nil {
		return customError.New(customError.InternalServerError, err)
	}

	logger.Logger.Debug("Request body parsed successfully", zap.String("username", req.Username))

	if err := c.Validate(req); err != nil {
		return customError.New(customError.BadRequest, err)
	}

	user, err := ac.userService.GetUserByUsername(req.Username)
	if err != nil {
		return customError.New(customError.InternalServerError, err)
	}

	if err := user.VerifyUser(req.Password); err != nil {
		return customError.New(customError.Forbidden, err)
	}

	token, err := jwt.GenerateJWT(user)
	if err != nil {
		return customError.New(customError.InternalServerError, err)
	}

	return response.Ok(c, "User authenticated successfully", map[string]interface{}{
		"access_token": token,
		"user":         user,
	})
}
