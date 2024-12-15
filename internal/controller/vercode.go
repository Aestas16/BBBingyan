package controller

import (
    "time"
    "net/http"
    "github.com/labstack/echo/v4"

    "user-management-system/internal/config"
    "user-management-system/internal/model"
    "user-management-system/internal/utils"
    "user-management-system/internal/controller/param"
)

func SendVerCode(c echo.Context) error {
    req := new(param.SendVerCodeRequest)
    if err := c.Bind(&req); err != nil {
        return echo.ErrBadRequest
    }
    user, err := model.FindUserByName(req.Username)
    if err == model.ErrUserNotFound {
        return echo.NewHTTPError(http.StatusForbidden, "user not found")
    } else if err != nil {
        return echo.ErrInternalServerError
    }
    lastVerCode, err := model.FindVerCodeByUserId(user.ID)
    if err == model.ErrVerCodeNotFound {
        verCode := model.VerCode{}
        verCode.UserId = user.ID
        verCode.Code = utils.GenerateVerCode(6)
        verCode.Time = time.Now().Unix()
        err = model.CreateVerCode(&verCode)
        if err != nil {
            return echo.ErrInternalServerError
        }
        if err := utils.SendVerCode(verCode.Code, user.Email); err != nil {
            return echo.ErrInternalServerError
        }
        var resp struct {
            Message string  `json:"message"`
        }
        resp.Message = "Success!"
        return c.JSON(http.StatusOK, &resp)
    } else if err != nil {
        return echo.ErrInternalServerError
    }
    nowTime := time.Now().Unix()
    if nowTime - lastVerCode.Time < config.Config.Server.Email.Interval {
        return echo.NewHTTPError(http.StatusForbidden, "access denied")
    }
    lastVerCode.Code = utils.GenerateVerCode(6)
    lastVerCode.Time = nowTime
    err = model.SaveVerCode(lastVerCode)
    if err != nil {
        return echo.ErrInternalServerError
    }
    if err := utils.SendVerCode(lastVerCode.Code, user.Email); err != nil {
        return echo.ErrInternalServerError
    }
    var resp struct {
        Message string  `json:"message"`
    }
    resp.Message = "Success!"
    return c.JSON(http.StatusOK, &resp) 
}