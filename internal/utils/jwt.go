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
    UserId  uint64
    IsAdmin bool
    jwt.RegisteredClaims
}

func GenerateToken(user *model.User, isAdmin bool) (string, error) {
    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        UserID: user.ID,
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