package routes

import (
	"net/http"
	"strings"

	"github.com/grqphical/postal-code-lookup-api/internal/lookup"
	"github.com/labstack/echo/v4"
)

func PostalCodeInfoGetHandler(c echo.Context) error {
	postalCode := strings.ToLower(c.Param("postalCode"))

	postalCodeObj, err := lookup.NewPostalCode(postalCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, postalCodeObj)
}
