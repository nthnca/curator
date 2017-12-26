package getphotos

import (
	"context"
	"encoding/hex"
	"fmt"
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

	fmt.Printf("mkdir .pics\n")
	size := len(mi.All())
	c := 0
	for i, _ := range mi.All() {
		iter := mi.All()[size-i-1]
		if iter.Deleted {
			continue
		}

		c++
		if c > 100 {
			break
		}

		fmt.Printf("gsutil cp gs://%s/%s .pics/\n",
			config.PhotoStorageBucket(), hex.EncodeToString(iter.Key))
		fmt.Printf("ln .pics/%s %s\n", hex.EncodeToString(iter.Key), iter.File[0].Filename)
	}
}
