package main

import (
	"log"

	"weaver-edge/config"
	"weaver-edge/internal/app"
)

func main() {
	cfg := config.MustLoad()

	if err := app.Migrate(cfg); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}
}
