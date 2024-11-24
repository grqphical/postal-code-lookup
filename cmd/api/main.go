package main

import (
	"os"

	"github.com/grqphical/postal-code-lookup-api/internal/routes"
	_ "github.com/joho/godotenv/autoload"
)

// @title Canadian Postal Code Lookup API
// @version 1.0
// @description An API to extract information about Canadian postal codes

// @contact.name grqphical
// @contact.url https://github.com/grqphical/postal-code-lookup

//	@license.name	GNU Public License Version 3.0
//	@license.url	https://github.com/grqphical/postal-code-lookup/blob/main/LICENSE

//	@host		localhost:8000
//	@BasePath	/v1
//  @produce json

func main() {
	e := routes.CreateRouter()

	e.Logger.Fatal(e.Start(os.Getenv("HOST_ADDR")))
}
