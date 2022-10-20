package main

import (
	"log"

	"github.com/AnNosov/communications_info/config"
	"github.com/AnNosov/communications_info/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)

}
