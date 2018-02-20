package fsckphotos

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/mediainfo/store"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func Register(app *kingpin.Application) {
	cmd := app.Command("fsck", "Validate photos are intact")
	cmd.Action(
		func(_ *kingpin.ParseContext) error {
			handler()
			return nil
		})
}

func handler() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	mi, err := store.New(ctx, client, config.MediaInfoBucket())
	if err != nil {
		log.Fatalf("New MediaInfo store failed: %v", err)
	}

	for _, m := range mi.All() {
		for _, f := range m.File {
			name := hex.EncodeToString(f.Sha256Sum)

			_, err := client.Bucket(config.PhotoStorageBucket()).Object(name).Attrs(ctx)
			if err == nil {
				continue
			}
			if err == storage.ErrObjectNotExist {
				fmt.Printf("Missing %s: %s\n", name, f.Filename)
				continue
			}
			log.Fatalf("Error checking for file: %v, %v", name, err)
		}
	}
}
