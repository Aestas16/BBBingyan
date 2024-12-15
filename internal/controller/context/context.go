package context

import (
    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"
)

func BindAndVali(c echo.Context, req interface{}) (err error) {
    err = c.Bind(req)
    if err != nil {
        return err
    }

    validate := validator.New()
    if err := validate.Struct(req); err != nil {
        return err
    }
    return nil
}