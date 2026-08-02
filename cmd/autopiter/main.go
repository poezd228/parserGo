package main

import (
	"log"

	"parser/internal/app/autopiter"
	autopiterconfig "parser/internal/config/autopiter"
)

func main() {
	cfg, err := autopiterconfig.Load("configs/autopiter/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	a := autopiter.NewApp(cfg)
	a.Start()
}
