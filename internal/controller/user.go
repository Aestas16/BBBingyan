package controller

import (
    "fmt"
    "time"
    "strconv"
    "crypto/md5"
    "github.com/labstack/echo/v4"
    "github.com/golang-jwt/jwt/v4"

    "user-management-system/internal/config"
    "user-management-system/internal/model"
    "user-management-system/internal/utils"
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
    if !claims.IsAdmin && claims.User.ID != id {
        return echo.NewHTTPError(403, "access denied")
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
    return c.JSON(200, &resp)
}

func UpdateUser(c echo.Context) error {
    claims := c.Get("claims").(*utils.Claims)
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        return echo.ErrNotFound
    }
    if !claims.IsAdmin && claims.User.ID != id {
        return echo.NewHTTPError(403, "access denied")
    }
    user, err := model.FindUserById(id)
    if err == model.ErrUserNotFound {
        return echo.ErrNotFound
    } else if err != nil {
        return echo.ErrInternalServerError
    }
    req := User{}
    if err := c.Bind(&req); err != nil {
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
    return c.JSON(200, &resp)
}

func DeleteUser(c echo.Context) error {
    claims := c.Get("claims").(utils.Claims)
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        return echo.ErrNotFound
    }
    if !claims.IsAdmin {
        return echo.NewHTTPError(403, "access denied")
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
    return c.JSON(200, &resp)
}

func RegisterUser(c echo.Context) error {
    req := User{}
    if err := c.Bind(&req); err != nil {
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
        return echo.NewHTTPError(403, "access denied")
    } else if err != nil {
        return echo.ErrInternalServerError
    }
    var resp struct {
        Message string  `json:"message"`
    }
    resp.Message = "Success!"
    return c.JSON(201, &resp)
}

func LoginUser(c echo.Context) error {
    var req struct {
        Username    string  `json:"username"`
        Password    string  `json:"password"`
    }
    if err := c.Bind(&req); err != nil {
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
        return c.JSON(200, &resp)
    }
    req.Password = fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))
    user, err := model.FindUserByName(req.Username)
    if err == model.ErrUserNotFound {
        return echo.NewHTTPError(403, "user not found")
    } else if err != nil {
        return echo.ErrInternalServerError
    }
    if user.Password != req.Password {
        return echo.NewHTTPError(403, "wrong password")
    }
    tokenString, err := utils.GenerateToken(user, false)
    if err != nil {
        return echo.ErrInternalServerError
    }
    var resp struct {
        TokenString string  `json:"token"`
    }
    resp.TokenString = tokenString
    return c.JSON(200, &resp)
}

func RefreshToken(c echo.Context) error {
    claims := c.Get("claims").(utils.Claims)
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
    return c.JSON(200, &resp)
}