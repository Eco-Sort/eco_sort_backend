package db

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
)

var GcStorage *storage.Client

func InitGcStorage() {
	ctx := context.Background()
	storage, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	GcStorage = storage
}
