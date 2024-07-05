package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateRouter() *echo.Echo {
	e := echo.New()

	e.GET("/", IndexGetHandler)

	v1 := e.Group("/v1")

	v1.GET("/valid/:postalCode", ValidPostalCodeGetHandler)

	return e
}

func IndexGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
