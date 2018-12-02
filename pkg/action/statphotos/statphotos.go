package statphotos

import (
	"context"
	"fmt"
	"log"
	"sort"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/pkg/config"
	"github.com/nthnca/curator/pkg/mediainfo"
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

	// Cfg is the configuration settings.
	Cfg *config.Config
}

// Do outputs a basic set of stats (number of files, storage space used, tags, etc) about the set of photos stored.
func Do(opts *Options) error {
	os := opts.ObjStore

	var all []*mediainfo.Media
	os.ForEach(func(key string, value []byte) {
		var m mediainfo.Media
		if er := proto.Unmarshal(value, &m); er != nil {
			log.Fatalf("Unmarshalling proto: %v", er)
		}
		all = append(all, &m)
	})

	sort.Slice(all, func(i, j int) bool {
		return all[i].Photo.TimestampSeconds < all[j].Photo.TimestampSeconds
	})

	totaltagcount := make(map[string]int)
	var totalsize int64

	ycount := 0
	ytagcount := make(map[string]int)
	var ysize int64

	for i, y := range all {
		name := util.GetCanonicalName(opts.Cfg, y, 0)
		sort.Strings(y.Tags)

		var size int64
		for _, f := range y.File {
			size += f.SizeInBytes
		}

		k := fmt.Sprintf("%s", y.Tags)
		totaltagcount[k] = totaltagcount[k] + 1
		totalsize += size

		ytagcount[k] = ytagcount[k] + 1
		ycount++
		ysize += size

		if i+1 == len(all) || name[:4] != util.GetCanonicalName(opts.Cfg, all[i+1], 0)[:4] {
			fmt.Printf("%s (%s)\n", name[:4], gb(ysize))
			var taglist []string
			for k := range ytagcount {
				taglist = append(taglist, k)
			}
			sort.Strings(taglist)
			for k := range taglist {
				fmt.Printf("  %s %d\n", k, ytagcount[k])
			}

			ycount = 0
			ytagcount = make(map[string]int)
			ysize = 0
		}
	}

	fmt.Printf("Totals (%s)\n", gb(totalsize))
	for k := range totaltagcount {
		fmt.Printf("  %s %d\n", k, totaltagcount[k])
	}
	return nil
}

func gb(bytes int64) string {
	return fmt.Sprintf("%d.%d GB", bytes/1000000000, (bytes/100000000)%10)
}
