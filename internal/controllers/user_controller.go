package controllers

import (
	"ptm/internal/di"
	"ptm/internal/dtos"
	"ptm/internal/services"
	"ptm/internal/utils/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserController interface {
	GetAllUsers(c echo.Context) error
	GetUserById(c echo.Context) error
	UpdateUser(c echo.Context) error
}

type userController struct {
	userService services.UserService
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
func (uc *userController) UpdateUser(c echo.Context) error {
	var request dtos.UpdateUserRequest

	if err := c.Bind(&request); err != nil {
		return response.BadRequest(c, "Error parsing request", err)
	}

	if err := c.Validate(request); err != nil {
		return response.BadRequest(c, "Error validating request", err)
	}

	idString := c.Param("id")
	num, err := strconv.Atoi(idString)
	if err != nil {
		return response.BadRequest(c, "Error converting string to integer", err)
	}

	updatedUser, err := uc.userService.UpdateUser(uint(num), &request)

	if err != nil {
		return err
	}

	return response.Ok(c, "Successfully Updated", updatedUser)
}

func (uc *userController) DeleteUser(c echo.Context) error {
	idString := c.Param("id")

	num, err := strconv.Atoi(idString)

	if err != nil {
		return response.BadRequest(c, "Error converting string to integer", err)
	}

	if err := uc.userService.DeleteUser(uint(num)); err != nil {
		return err
	}

	return response.Ok(c, "Successfully Deleted", nil)
}
