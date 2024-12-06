package main

import (
	"ptm/config"
	"ptm/db"
	"ptm/routes"
	"ptm/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	config.InitConfig()
	utils.InitLogger()

	db.InitDB()

	e := echo.New()

	routes.InitRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
