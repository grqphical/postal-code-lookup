package routes

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/grqphical/postal-code-lookup-api/internal/database"
	"github.com/grqphical/postal-code-lookup-api/internal/lookup"
	"github.com/labstack/echo/v4"

	_ "github.com/grqphical/postal-code-lookup-api/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// httpError model info
// @Description JSON object returned whenever an error occurs
type HttpError struct {
	Message string `json:"message"`
}

// struct to store shared state between routes
type server struct {
	db *sql.DB
}

func CreateRouter() *echo.Echo {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	s := server{
		db,
	}

	e := echo.New()

	e.GET("/", s.IndexGetHandler)

	e.GET("/docs/*", echoSwagger.WrapHandler)

	v1 := e.Group("/v1")

	v1.GET("/postal-code/:postalCode", s.PostalCodeInfoGetHandler)

	return e
}

// IndexGetHandler godoc
// @Summary check the status of the API
// @Produce plain
// @Success 200 {} string
// @Router / [get]
func (s *server) IndexGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

// PostalCodeInfoGetHandler godoc
// @Summary extracts information about a postal code
// @Produce json
// @Param postalCode path string true "Postal Code to extract info from"
// @Success 200 {object} lookup.PostalCode
// @Failure 400 {object} HttpError
// @Router /postal-code/{postalCode} [get]
func (s *server) PostalCodeInfoGetHandler(c echo.Context) error {
	c.Response().Header().Add("Content-Type", "application/json")
	c.Response().Header().Add("Allow", "GET")

	postalCode := strings.ToLower(c.Param("postalCode"))

	postalCodeObj, err := lookup.NewPostalCode(postalCode, s.db)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, postalCodeObj)
}
