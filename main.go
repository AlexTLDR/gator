package main

import (
	"fmt"
	"gator/internal/config"
	"log"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	err = cfg.SetUser("Alex")
	if err != nil {
		log.Fatalf("Failed to set user: %v", err)
	}

	updatedCfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to reload configuration: %v", err)
	}

	fmt.Printf("%+v\n", updatedCfg)
}
