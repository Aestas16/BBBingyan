package main

import (
    "strconv"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"

    "user-management-system/internal/config"
    "user-management-system/internal/model"
    "user-management-system/internal/router"
)

func main() {
    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    config.InitConfig()
    model.InitDB()
    router.InitRouter(e)
    e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.Config.Server.Port)))
}