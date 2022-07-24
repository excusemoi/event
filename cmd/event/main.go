package main

import (
	"event/pkg/config"
	"event/pkg/models"
	"event/pkg/repository/cache"
	"event/pkg/server"
	"log"
)

var (
	cfg = &config.Config{
		Host: "localhost",
		Port: "8000",
	}
)

func main() {
	db := cache.NewEventsCache(make(map[uint64]*models.Event))
	s := server.NewServer(db)
	s.InitServer(cfg)
	log.Println("Server starting...")
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
