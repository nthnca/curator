package getphotos

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/util"
)

const (
	dryRun = false
)

var (
	photoData util.PhotoInfo
)

func Handler() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	photoData.Load(ctx, client)

	fmt.Printf("mkdir .pics\n")
	size := len(photoData.Data)
	for i, _ := range photoData.Data {
		if i >= 25 {
			break
		}

		iter := photoData.Data[size-i-1]
		fmt.Printf("gsutil cp gs://%s/%s .pics/%s\n",
			config.PhotoStorageBucket(), iter.GetPath(), hex.EncodeToString(iter.GetSha256Sum()))
		fmt.Printf("ln .pics/%s %s\n", hex.EncodeToString(iter.GetSha256Sum()),
			iter.GetPath())
	}
}
