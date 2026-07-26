package main

import (
	"log"

	"github.com/nongpal/Palge-Backend/internal/api"
)

func main() {
	var cfg api.Config
	api.NewConfig(&cfg)

	app, err := api.NewApplication(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer app.Close()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
