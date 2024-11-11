package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/grqphical/postal-code-lookup-api/internal/database"
	"github.com/labstack/echo/v4"
)

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

	v1 := e.Group("/v1")

	v1.GET("/postal-code/:postalCode", s.PostalCodeInfoGetHandler)

	return e
}

func (s *server) IndexGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
