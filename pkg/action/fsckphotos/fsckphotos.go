package fsckphotos

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
	"google.golang.org/api/iterator"
)

// Options allows you to modify the behavior of the FsckPhotos action.
type Options struct {
	// Ctx is a valid context.Context to run this command under.
	Ctx context.Context

	// Storage is a Google Cloud Storage client.
	Storage *storage.Client

	// ObjStore is an ObjectStore client
	ObjStore *objectstore.ObjectStore

	// Cfg is the configuration settings.
	Cfg *config.Config
}

// Do performs a system integrity check.
func Do(opts *Options) error {
	ctx := opts.Ctx
	client := opts.Storage
	os := opts.ObjStore

	totalObjects := 0
	wantedObjects := 0
	extraObjects := 0
	missingObjects := 0
	wanted := make(map[string][]byte)
	have := make(map[string]bool)

	// What does the object store have.
	os.ForEach(func(key string, value []byte) {
		var m message.Media
		if er := proto.Unmarshal(value, &m); er != nil {
			log.Fatalf("Unmarshalling proto: %v", er)
		}
		if key != hex.EncodeToString(m.File[0].Sha256Sum) {
			log.Fatalf("oops")
		}
		if key != hex.EncodeToString(m.Key) {
			log.Fatalf("oops")
		}
		for _, f := range m.File {
			name := hex.EncodeToString(f.Sha256Sum)
			wanted[name] = f.Md5Sum
			wantedObjects += 1
		}
	})

	// Any storage entries that we didn't expect?
	bkt := client.Bucket(opts.Cfg.PhotoStorageBucket())
	for it := bkt.Objects(ctx, nil); ; {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate through objects: %v", err)
		}

		have[obj.Name] = true
		totalObjects += 1
		check := wanted[obj.Name]
		if len(check) > 0 {
			if util.MD5(check) != util.MD5(obj.MD5) {
				log.Fatalf("Ooops")
			}
		} else {
			if extraObjects == 0 {
				fmt.Printf("Extra Objects:\n")
			}
			fmt.Printf("  %v\n", obj.Name)
			extraObjects += 1

			// Uncomment this at your own risk, you can use it to clean things up.
			// err := client.Bucket(config.PSBucket()).Object(obj.Name).Delete(ctx)
			// err := obj.Delete(ctx)
			// log.Printf("%v", err)
			// exit(1)
		}
	}

	// Were any storage entries missing?
	for k, v := range wanted {
		if len(v) > 0 || have[k] {
			continue
		}

		if missingObjects == 0 {
			fmt.Printf("Missing objects:\n")
		}

		fmt.Printf("  %s\n", k)
		missingObjects += 1
	}

	/*
		        // Uncomment this if you want to delete Info relating to missing objects.
			for _, b := range del {
				log.Printf("%v", b)
				mi.DeleteFast(util.Sha256(b))
			}
			mi.Flush(ctx, client)
	*/

	log.Printf("Total Objects: %d", totalObjects)
	log.Printf("Wanted Objects: %d", wantedObjects)
	log.Printf("Extra Objects: %d", extraObjects)
	log.Printf("Missing Objects: %d", missingObjects)

	return nil
}
