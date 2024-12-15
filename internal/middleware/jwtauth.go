package middleware

import (
    "errors"
	"net/http"
    "github.com/labstack/echo/v4"
)

var ErrInvalidToken = errors.New("invalid token")
var ErrTokenExpired = errors.New("token expired")

func JWTAuthMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            tokenString := c.Request().Header.Get("Authorization")
            if tokenString == "" {
                return echo.NewHTTPError(http.StatusUnauthorized, "token not found")
            }
            claims, err := ParseToken(tokenString)
            if err == ErrInvalidToken {
                return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
            } else if err == ErrTokenExpired {
                return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
            } else if err != nil {
                return echo.ErrInternalServerError
            }
            c.Set("claims", claims)
            return next(c);
        }
    }
}