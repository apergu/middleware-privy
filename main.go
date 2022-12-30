package main

import (
	_ "github.com/joho/godotenv/autoload" // for development
	"gitlab.com/mohamadikbal/project-privy/cmd/api"
)

func main() {
	api.Execute()
}
