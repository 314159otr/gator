package main

import (
	"fmt"
	"log"
	"github.com/314159otr/gator/internal/config"
)

func main() {
	var cfg config.Config
	cfg, err := config.Read() 
	if err != nil {
		log.Fatalf("error reading file. Error: %s", err)
	}

	err = cfg.SetUser("piotr")
	if err != nil {
		log.Fatalf("error reading file. Error: %s", err)
	}

	cfg, err = config.Read() 
	if err != nil {
		log.Fatalf("error reading file. Error: %s", err)
	}

	fmt.Print(cfg)
}
