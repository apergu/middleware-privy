package main

import (
	"middleware/cmd/api"

	_ "github.com/joho/godotenv/autoload" // for development
)

func main() {
	api.Execute()
}
