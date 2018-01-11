package statphotos

import (
	"context"
	"fmt"
	"log"
	"sort"

	"cloud.google.com/go/storage"
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/mediainfo/store"
	"github.com/nthnca/curator/util"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func Register(app *kingpin.Application) {
	app.Command("stats", "analyze curator data").Action(
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

	arr := mi.All()
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Photo.EpochInSeconds < arr[j].Photo.EpochInSeconds
	})

	tagcount := make(map[string]int)
	var totalsize int64

	ycount := 0
	ytagcount := make(map[string]int)
	var ysize int64

	for i, y := range arr {
		name := util.GetCanonicalName(y)
		sort.Strings(y.Tags)

		var size int64
		for _, f := range y.File {
			size += f.SizeInBytes
		}

		k := fmt.Sprintf("%s", y.Tags)
		tagcount[k] = tagcount[k] + 1
		totalsize += size

		ytagcount[k] = ytagcount[k] + 1
		ycount++
		ysize += size

		if i+1 == len(arr) || name[:4] != util.GetCanonicalName(arr[i+1])[:4] {
			fmt.Printf("%s (%d)\n", name[:4], ysize/100000000)
			for k := range ytagcount {
				fmt.Printf("  %s %d\n", k, ytagcount[k])
			}

			ycount = 0
			ytagcount = make(map[string]int)
			ysize = 0
		}
	}

	fmt.Printf("Totals (%d)\n", totalsize/100000000)
	for k := range tagcount {
		fmt.Printf("  %s %d\n", k, tagcount[k])
	}
}
