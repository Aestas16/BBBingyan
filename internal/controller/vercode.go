package controller

import (
    "time"
    "github.com/labstack/echo/v4"

    "user-management-system/internal/config"
    "user-management-system/internal/model"
    "user-management-system/internal/utils"
)

func SendVerCode(c echo.Context) error {
    var req struct {
        Username    string    `json:"username"`
    }
    if err := c.Bind(&req); err != nil {
        return echo.ErrBadRequest
    }
    user, err := model.FindUserByName(req.Username)
    if err == model.ErrUserNotFound {
        return echo.NewHTTPError(403, "user not found")
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
        return c.JSON(200, &resp)
    } else if err != nil {
        return echo.ErrInternalServerError
    }
    nowTime := time.Now().Unix()
    if nowTime - lastVerCode.Time < config.Config.Server.Email.Interval {
        return echo.NewHTTPError(403, "access denied")
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
    return c.JSON(200, &resp) 
}