package router

import (
    "fmt"
    "github.com/labstack/echo/v4"

    "user-management-system/internal/config"
    "user-management-system/internal/controller"
    "user-management-system/internal/utils"
)

func InitRouter(e *echo.Echo) {
    apiVersion := config.Config.Server.Version
    apiUser := e.Group(fmt.Sprintf("/%s/user", apiVersion))
    apiUser.Use(utils.JWTAuthMiddleware())
    apiUser.GET("/:id", controller.UserInfo)
    apiUser.POST("/:id", controller.UserInfo)
    apiUser.POST("/:id/update", controller.UpdateUser)
    apiUser.POST("/:id/delete", controller.DeleteUser)
    apiUser.POST("/token", controller.RefreshToken)
    api := e.Group(fmt.Sprintf("/%s", apiVersion))
    api.POST("/login", controller.LoginUser)
    api.POST("/register", controller.RegisterUser)
}