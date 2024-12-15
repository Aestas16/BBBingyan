package controller

import (
    "fmt"
    "time"
    "strconv"
    "crypto/md5"
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/golang-jwt/jwt/v4"

    "user-management-system/internal/config"
    "user-management-system/internal/model"
    "user-management-system/internal/utils"
    "user-management-system/internal/controller/param"
    "user-management-system/internal/controller/context"
)

type User struct {
    Username    string  `json:"username"`
    Password    string  `json:"password"`
    Email       string  `json:"email"`
}

func UserInfo(c echo.Context) error {
    claims := c.Get("claims").(*utils.Claims)
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        return echo.ErrNotFound
    }
    if !claims.IsAdmin && claims.UserId != id {
        return echo.NewHTTPError(http.StatusForbidden, "access denied")
    }
    user, err := model.FindUserById(id)
    if err == model.ErrUserNotFound {
        return echo.ErrNotFound
    } else if err != nil {
        return echo.ErrInternalServerError
    }
    var resp struct {
        Username    string  `json:"username"`
        Email       string  `json:"email"`
    }
    resp.Username = user.Username
    resp.Email = user.Email
    return c.JSON(http.StatusOK, &resp)
}

func UpdateUser(c echo.Context) error {
    claims := c.Get("claims").(*utils.Claims)
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        return echo.ErrNotFound
    }
    if !claims.IsAdmin && claims.UserId != id {
        return echo.NewHTTPError(http.StatusForbidden, "access denied")
    }
    user, err := model.FindUserById(id)
    if err == model.ErrUserNotFound {
        return echo.ErrNotFound
    } else if err != nil {
        return echo.ErrInternalServerError
    }
    req := new(param.UserRequest)
    if err := context.BindAndVali(c, req); err != nil {
        return echo.ErrBadRequest
    }
    user.Username = req.Username
    user.Password = fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))
    user.Email = req.Email
    if err := model.SaveUser(user); err != nil {
        return echo.ErrInternalServerError
    }
    var resp struct {
        Message string  `json:"message"`
    }
    resp.Message = "Success!"
    return c.JSON(http.StatusOK, &resp)
}

func DeleteUser(c echo.Context) error {
    claims := c.Get("claims").(*utils.Claims)
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        return echo.ErrNotFound
    }
    if !claims.IsAdmin {
        return echo.NewHTTPError(http.StatusForbidden, "access denied")
    }
    _, err = model.FindUserById(id)
    if err == model.ErrUserNotFound {
        return echo.ErrNotFound
    } else if err != nil {
        return echo.ErrInternalServerError
    }
    if err := model.DeleteUserById(id); err != nil {
        return echo.ErrInternalServerError
    }
    var resp struct {
        Message string  `json:"message"`
    }
    resp.Message = "Success!"
    return c.JSON(http.StatusOK, &resp)
}

func RegisterUser(c echo.Context) error {
    req := new(param.UserRequest)
    if err := context.BindAndVali(c, req); err != nil {
        return echo.ErrBadRequest
    }
    if req.Username == "" || req.Password == "" {
        return echo.ErrBadRequest
    }
    user := model.User{}
    user.Username = req.Username
    user.Password = fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))
    user.Email = req.Email
    err := model.CreateUser(&user)
    if err == model.ErrUserAlreadyExist {
        return echo.NewHTTPError(http.StatusForbidden, "access denied")
    } else if err != nil {
        return echo.ErrInternalServerError
    }
    var resp struct {
        Message string  `json:"message"`
    }
    resp.Message = "Success!"
    return c.JSON(http.StatusCreated, &resp)
}

func LoginUser(c echo.Context) error {
    req := new(param.LoginUserRequest)
    if err := context.BindAndVali(c, req); err != nil {
        return echo.ErrBadRequest
    }
    if req.Username == config.Config.Server.Admin.Username && req.Password == config.Config.Server.Admin.Password {
        tokenString, err := utils.GenerateToken(&model.User{}, true)
        if err != nil {
            return echo.ErrInternalServerError
        }
        var resp struct {
            TokenString string  `json:"token"`
        }
        resp.TokenString = tokenString
        return c.JSON(http.StatusOK, &resp)
    }
    req.Password = fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))
    user, err := model.FindUserByName(req.Username)
    if err == model.ErrUserNotFound {
        return echo.NewHTTPError(http.StatusForbidden, "user not found")
    } else if err != nil {
        return echo.ErrInternalServerError
    }
    verCode, err := model.FindVerCodeByUserId(user.ID)
    if err == model.ErrVerCodeNotFound {
        return echo.NewHTTPError(http.StatusForbidden, "verification code not found")
    } else if err != nil {
        return echo.ErrInternalServerError
    }
    nowTime := time.Now().Unix()
    if nowTime - verCode.Time > config.Config.Server.Email.Interval {
        return echo.NewHTTPError(http.StatusForbidden, "verification code expired")
    }
    if req.Code != verCode.Code {
        return echo.NewHTTPError(http.StatusForbidden, "wrong verification code")
    }
    if user.Password != req.Password {
        return echo.NewHTTPError(http.StatusForbidden, "wrong password")
    }
    tokenString, err := utils.GenerateToken(user, false)
    if err != nil {
        return echo.ErrInternalServerError
    }
    var resp struct {
        TokenString string  `json:"token"`
    }
    resp.TokenString = tokenString
    return c.JSON(http.StatusOK, &resp)
}

func RefreshToken(c echo.Context) error {
    claims := c.Get("claims").(*utils.Claims)
    expirationTime := time.Now().Add(5 * time.Minute)
    claims.ExpiresAt.Time = expirationTime
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(utils.JWTKey)
    if err != nil {
        return echo.ErrInternalServerError
    }
    var resp struct {
        TokenString string  `json:"token"`
    }
    resp.TokenString = tokenString
    return c.JSON(http.StatusOK, &resp)
}