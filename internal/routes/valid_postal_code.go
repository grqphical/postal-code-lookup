package routes

import (
	"net/http"
	"strings"

	"github.com/grqphical/postal-code-lookup-api/internal/lookup"
	"github.com/labstack/echo/v4"
)

func ValidPostalCodeGetHandler(c echo.Context) error {
	postalCode := strings.ToLower(c.Param("postalCode"))

	if !lookup.IsValidPostalCode(postalCode) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Canadian Postal Code")
	}

	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "Valid",
	})
}
