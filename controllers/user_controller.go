package controllers

import (
	"net/http"
	"ptm/models"

	"github.com/labstack/echo/v4"
)

func GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := models.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}
