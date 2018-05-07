package getphotos

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"github.com/nthnca/curator/pkg/config"
	"github.com/nthnca/curator/pkg/mediainfo/store"
	"github.com/nthnca/curator/pkg/util"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	filter string
	max    int
	tags   util.Tags
)

func Register(app *kingpin.Application) {
	cmd := app.Command("get", "Create script for copying photos")
	cmd.Action(
		func(_ *kingpin.ParseContext) error {
			handler()
			return nil
		})
	cmd.Flag("filter", "description").StringVar(&filter)
	cmd.Flag("max", "The maximum number of results to return").IntVar(&max)
	cmd.Flag("has", "Has labels").StringsVar(&tags.A)
	cmd.Flag("not", "Not labels").StringsVar(&tags.B)
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

	tags.Normalize()
	tags.Validate(config.ValidLabels())

	count := 0
	for i := len(mi.All()) - 1; i >= 0; i-- {
		iter := mi.All()[i]
		if !tags.Match(iter.Tags) {
			continue
		}

		name := util.GetCanonicalName(iter)
		if filter != "" && filter != name[:len(filter)] {
			continue
		}

		count++
		fmt.Printf("%s %s\n", hex.EncodeToString(iter.Key), name)

		if max != 0 && count >= max {
			break
		}
	}

	log.Printf("Photos retrieved: %d", count)
}
