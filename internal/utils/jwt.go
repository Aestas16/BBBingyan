package utils

import (
    "time"
    "errors"
    "github.com/golang-jwt/jwt/v4"
    "github.com/labstack/echo/v4"

    "user-management-system/internal/model"
    "user-management-system/internal/config"
)

var JWTKey = []byte(config.Config.Server.JwtKey)
var ErrInvalidToken = errors.New("invalid token")
var ErrTokenExpired = errors.New("token expired")

type Claims struct {
    User *model.User
    IsAdmin bool
    jwt.RegisteredClaims
}

func GenerateToken(user *model.User, isAdmin bool) (string, error) {
    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        User: user,
        IsAdmin: isAdmin,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(JWTKey)
    return tokenString, err
}

func ParseToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return JWTKey, nil
    })
    if !token.Valid || err != nil {
        return claims, ErrInvalidToken
    }
    if time.Until(claims.ExpiresAt.Time) < 0 {
        return claims, ErrTokenExpired
    }
    return claims, err
}

func JWTAuthMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            tokenString := c.Request().Header.Get("Authorization")
            if tokenString == "" {
                return echo.NewHTTPError(401, "token not found")
            }
            claims, err := ParseToken(tokenString)
            if err == ErrInvalidToken {
                return echo.NewHTTPError(401, "invalid token")
            } else if err == ErrTokenExpired {
                return echo.NewHTTPError(401, "token expired")
            } else if err != nil {
                return echo.ErrInternalServerError
            }
            c.Set("claims", claims)
            return next(c);
        }
    }
}