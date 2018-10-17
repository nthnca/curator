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

	os, err := objectstore.New(ctx, client, config.MetadataBucket(), config.MetadataPath())
	if err != nil {
		log.Fatalf("New ObjectStore failed: %v", err)
	}

	totalObjects := 0
	extraObjects := 0
	missingObjects := 0
	wanted := make(map[string]bool)
	have := make(map[string]bool)

	os.ForEach(func(key string, value []byte) {
		var m message.Media
		if er := proto.Unmarshal(value, &m); er != nil {
			log.Fatalf("Unmarshalling proto: %v", er)
		}
		for _, f := range m.File {
			name := hex.EncodeToString(f.Sha256Sum)
			wanted[name] = true
			totalObjects += 1
		}
	})

	bkt := client.Bucket(config.PhotoStorageBucket())
	for it := bkt.Objects(ctx, nil); ; {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate through objects: %v", err)
		}

		have[obj.Name] = true
		if wanted[obj.Name] == false {
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

	del := [][]byte{}
	os.ForEach(func(key string, value []byte) {
		var m message.Media
		if er := proto.Unmarshal(value, &m); er != nil {
			log.Fatalf("Unmarshalling proto: %v", er)
		}

		flag := false
		for _, f := range m.File {
			name := hex.EncodeToString(f.Sha256Sum)
			if have[name] == false {
				flag = true
			}
		}

		if !flag {
			return
		}

		if missingObjects == 0 {
			fmt.Printf("Missing objects:\n")
		}
		fmt.Printf("  %s\n", util.GetCanonicalName(&m))
		del = append(del, m.Key)
		for _, f := range m.File {
			name := hex.EncodeToString(f.Sha256Sum)
			fmt.Printf("    %s %v\n", name, have[name])
			if have[name] == false {
				missingObjects += 1
			}
		}
	})

	/*
		        // Uncomment this if you want to delete Info relating to missing objects.
			for _, b := range del {
				log.Printf("%v", b)
				mi.DeleteFast(util.Sha256(b))
			}
			mi.Flush(ctx, client)
	*/

	log.Printf("Total Objects: %d", totalObjects)
	log.Printf("Extra Objects: %d", extraObjects)
	log.Printf("Missing Objects: %d", missingObjects)

}
