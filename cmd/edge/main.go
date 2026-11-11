package main

import (
	"log"

	"weaver-edge/config"
	"weaver-edge/internal/app"
)

func main() {
	cfg := config.MustLoad()

	if err := app.Run(cfg); err != nil {
		log.Fatalf("application failed: %v", err)
	}
}
