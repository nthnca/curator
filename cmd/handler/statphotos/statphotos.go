package statphotos

import (
	"context"
	"log"
	"sort"

	"cloud.google.com/go/storage"
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/mediainfo/store"
	"github.com/nthnca/curator/util"
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

	var tags util.Tags
	tags.A = []string{"keep"}

	var totalsize int64
	total := len(mi.All())
	del := 0
	yeart := make(map[string]int)
	yeard := make(map[string]int)
	years := make(map[string]int64)
	for _, y := range mi.All() {
		name := util.GetCanonicalName(y)
		var size int64
		for _, f := range y.File {
			size += f.SizeInBytes
		}

		totalsize += size
		years[name[:4]] += size
		if !tags.Match(y.Tags) {
			yeard[name[:4]]++
			del++
		} else {
			yeart[name[:4]]++
		}
	}
	var keys []string
	for k := range yeart {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		log.Printf("Photos %v: %v (%v) %d", k, yeart[k], yeard[k], years[k]/100000000)
	}
	log.Printf("Photos %d (deleted: %d, total: %d) %d",
		total-del, del, total, totalsize/100000000)
}
