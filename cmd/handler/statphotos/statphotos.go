package statphotos

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/mediainfo/store"
)

func Handler() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	mi, err := store.New(ctx, client, config.MediaInfoBucket())
	if err != nil {
		log.Fatalf("New MediaInfo store failed: %v", err)
	}

	total := len(mi.All())
	del := 0
	for _, y := range mi.All() {
		if y.Deleted {
			del++
		}
	}
	log.Printf("Photos %d (deleted: %d, total: %d)", total-del, del, total)
}
