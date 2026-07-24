package main

import (
	"log"

	"github.com/nongpal/Palge-Backend/internal/api"
)

func main() {
	cfg := api.Config{
		Port: 4000,
		Env:  "development",
	}

	app := api.NewApplication(cfg)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
