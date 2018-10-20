package getphotos

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/pkg/config"
	"github.com/nthnca/curator/pkg/mediainfo"
	"github.com/nthnca/curator/pkg/util"
	objectstore "github.com/nthnca/object-store"
)

// Options allows you to modify the behavior of the GetPhotos action.
type Options struct {
	// Ctx is a valid context.Context to run this command under.
	Ctx context.Context

	// Storage is a Google Cloud Storage client.
	Storage *storage.Client

	// ObjStore is an ObjectStore client.
	ObjStore *objectstore.ObjectStore

	// Cfg is the configuration settings.
	Cfg *config.Config

	// Filter is a prefix which to match photots to.
	Filter string

	// Max is the maximum number of photos to list.
	Max int

	// Tags is a query that will be used to filter which photos to list.
	Tags util.Tags
}

// Do retrieves the set of photos that match the given parameters.
func Do(opts *Options) error {
	os := opts.ObjStore

	opts.Tags.Normalize()
	opts.Tags.Validate(opts.Cfg.ValidLabels())

	count := 0
	os.ForEach(func(key string, value []byte) {
		if opts.Max != 0 && count >= opts.Max {
			return
		}

		var m mediainfo.Media
		if er := proto.Unmarshal(value, &m); er != nil {
			log.Fatalf("Unmarshalling proto: %v", er)
		}

		iter := m
		if !opts.Tags.Match(iter.Tags) {
			return
		}

		name := util.GetCanonicalName(opts.Cfg, &iter)
		if opts.Filter != "" && opts.Filter != name[:len(opts.Filter)] {
			return
		}

		count++
		fmt.Printf("%s %s\n", hex.EncodeToString(iter.Key), name)
	})

	log.Printf("Photos retrieved: %d", count)
	return nil
}
