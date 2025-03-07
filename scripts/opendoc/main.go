package main

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/fivetentaylor/pointy/pkg/rogue"
	"github.com/fivetentaylor/pointy/pkg/storage/redis"
	"github.com/fivetentaylor/pointy/pkg/storage/s3"
)

func main() {
	docID := os.Args[1]
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	redis, err := redis.NewRedis()
	if err != nil {
		log.Fatalf("failed to create redis client: %s", err)
	}

	s3, err := s3.NewS3()
	if err != nil {
		log.Fatalf("failed to create s3 client: %s", err)
	}

	docStore := rogue.NewDocStore(s3, nil, redis)

	ctx := context.Background()
	_, doc, err := docStore.GetCurrentDoc(ctx, docID)
	if err != nil {
		log.Fatalf("failed to create load doc: %s", err)
	}

	log.Printf("DOC TEXT:\n%q", doc.GetText())
}
