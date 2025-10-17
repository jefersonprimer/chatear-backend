package main

import (
	"fmt"
	"log"

	"github.com/jefersonprimer/chatear-backend/config"
	"github.com/jefersonprimer/chatear-backend/pkg/api"
)

func main() {
	cfg := config.LoadConfig()

	r, err := api.SetupServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("🚀 Chatear Backend running on :%d", cfg.Port)
	r.Run(fmt.Sprintf(":%d", cfg.Port))
}