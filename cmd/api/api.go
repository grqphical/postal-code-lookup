package main

import (
	"os"

	"github.com/grqphical/postal-code-lookup-api/internal/routes"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	e := routes.CreateRouter()

	e.Logger.Fatal(e.Start(os.Getenv("HOST_ADDR")))
}
