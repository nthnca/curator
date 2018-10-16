package getphotos

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/pkg/config"
	"github.com/nthnca/curator/pkg/mediainfo/message"
	"github.com/nthnca/curator/pkg/util"
	objectstore "github.com/nthnca/object-store"
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

	os, err := objectstore.New(ctx, client, config.MetadataBucket(), config.MetadataPath())
	if err != nil {
		log.Fatalf("New ObjectStore failed: %v", err)
	}

	tags.Normalize()
	tags.Validate(config.ValidLabels())

	count := 0
	os.ForEach(func(key string, value []byte) {
		if max != 0 && count >= max {
			return
		}

		var m message.Media
		if er := proto.Unmarshal(value, &m); er != nil {
			log.Fatalf("Unmarshalling proto: %v", er)
		}

		iter := m
		if !tags.Match(iter.Tags) {
			return
		}

		name := util.GetCanonicalName(&iter)
		if filter != "" && filter != name[:len(filter)] {
			return
		}

		count++
		fmt.Printf("%s %s\n", hex.EncodeToString(iter.Key), name)
	})

	log.Printf("Photos retrieved: %d", count)
}
