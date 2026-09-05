package handler

import (
	"github.com/labstack/echo/v4"
)

func SendSuccess(c echo.Context, status int, message string, data interface{}) error {
	res := echo.Map{
		"code":    status,
		"status":  "success",
		"message": message,
	}
	if data != nil {
		if m, ok := data.(echo.Map); ok {
			for k, v := range m {
				res[k] = v
			}
		} else {
			res["data"] = data
		}
	}
	return c.JSON(status, res)
}

func SendError(c echo.Context, status int, message string) error {
	return c.JSON(status, echo.Map{
		"code":    status,
		"status":  "error",
		"message": message,
	})
}
