package statphotos

import (
	"context"
	"fmt"
	"log"
	"sort"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/pkg/config"
	"github.com/nthnca/curator/pkg/mediainfo/message"
	"github.com/nthnca/curator/pkg/util"
	objectstore "github.com/nthnca/object-store"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func Register(app *kingpin.Application) {
	app.Command("stats", "analyze curator data").Action(
		func(_ *kingpin.ParseContext) error {
			handler()
			return nil
		})
}

func gb(bytes int64) string {
	return fmt.Sprintf("%d.%d GB", bytes/1000000000, (bytes/100000000)%10)
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

	var arr []*message.Media
	os.ForEach(func(key string, value []byte) {
		var m message.Media
		if er := proto.Unmarshal(value, &m); er != nil {
			log.Fatalf("Unmarshalling proto: %v", er)
		}
		arr = append(arr, &m)
	})

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Photo.TimestampSeconds < arr[j].Photo.TimestampSeconds
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
			fmt.Printf("%s (%s)\n", name[:4], gb(ysize))
			for k := range ytagcount {
				fmt.Printf("  %s %d\n", k, ytagcount[k])
			}

			ycount = 0
			ytagcount = make(map[string]int)
			ysize = 0
		}
	}

	fmt.Printf("Totals (%s)\n", gb(totalsize))
	for k := range tagcount {
		fmt.Printf("  %s %d\n", k, tagcount[k])
	}
}
