package controller

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/redis/go-redis/v9"

    "user-management-system/internal/config"
    "user-management-system/internal/model"
    "user-management-system/internal/utils"
    "user-management-system/internal/controller/param"
    "user-management-system/internal/controller/context"
)

func SendVerCode(c echo.Context) error {
    req := new(param.SendVerCodeRequest)
    if err := context.BindAndVali(c, req); err != nil {
        return echo.ErrBadRequest
    }
    user, err := model.FindUserByName(req.Username)
    if err == model.ErrUserNotFound {
        return echo.NewHTTPError(http.StatusForbidden, "user not found")
    } else if err != nil {
        return echo.ErrInternalServerError
    }
    exists, err := utils.CheckVerCodeExist(user.Email)
    if err != nil {
        return echo.ErrInternalServerError
    } else if exists == true {
        return echo.NewHTTPError(http.StatusForbidden, "verification code already exist")
    }
    if err := utils.SendVerCode(user.Email); err != nil {
        return echo.ErrInternalServerError
    }
    var resp struct {
        Message string  `json:"message"`
    }
    resp.Message = "Success!"
    return c.JSON(http.StatusOK, &resp) 
}