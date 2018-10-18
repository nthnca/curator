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
)

type Options struct {
	Filter string
	Max    int
	Tags   util.Tags
}

func Do(opts *Options) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	os, err := objectstore.New(ctx, client, config.MetadataBucket(), config.MetadataPath())
	if err != nil {
		log.Fatalf("New ObjectStore failed: %v", err)
	}

	opts.Tags.Normalize()
	opts.Tags.Validate(config.ValidLabels())

	count := 0
	os.ForEach(func(key string, value []byte) {
		if opts.Max != 0 && count >= opts.Max {
			return
		}

		var m message.Media
		if er := proto.Unmarshal(value, &m); er != nil {
			log.Fatalf("Unmarshalling proto: %v", er)
		}

		iter := m
		if !opts.Tags.Match(iter.Tags) {
			return
		}

		name := util.GetCanonicalName(&iter)
		if opts.Filter != "" && opts.Filter != name[:len(opts.Filter)] {
			return
		}

		count++
		fmt.Printf("%s %s\n", hex.EncodeToString(iter.Key), name)
	})

	log.Printf("Photos retrieved: %d", count)
}
