package statphotos

import (
	"context"
	"fmt"
	"log"
	"sort"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/pkg/mediainfo/message"
	"github.com/nthnca/curator/pkg/util"
	objectstore "github.com/nthnca/object-store"
)

// Options allows you to modify the behavior of the StatPhotos action.
type Options struct {
	// Ctx is a valid context.Context to run this command under.
	Ctx context.Context

	// Storage is a Google Cloud Storage client.
	Storage *storage.Client

	// ObjStore is an ObjectStore client
	ObjStore *objectstore.ObjectStore
}

// Do outputs a basic set of stats (number of files, storage space used, tags, etc) about the set of photos stored.
func Do(opts *Options) {
	os := opts.ObjStore

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

func gb(bytes int64) string {
	return fmt.Sprintf("%d.%d GB", bytes/1000000000, (bytes/100000000)%10)
}
