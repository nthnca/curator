package mutatephotos

import (
	"bufio"
	"context"
	"encoding/hex"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/pkg/config"
	"github.com/nthnca/curator/pkg/mediainfo"
	"github.com/nthnca/curator/pkg/util"
	objectstore "github.com/nthnca/object-store"
)

// Options allows you to modify the behavior of the MutatePhotos action.
type Options struct {
	// Ctx is a valid context.Context to run this command under.
	Ctx context.Context

	// Storage is a Google Cloud Storage client.
	Storage *storage.Client

	// ObjStore is an ObjectStore client.
	ObjStore *objectstore.ObjectStore

	// Cfg is the configuration settings.
	Cfg *config.Config

	// Tags is the set of photos to add and remove.
	Tags util.Tags

	// DryRun if True then actually perform the modifications, otherwise just print what would happen.
	DryRun bool
}

// Do performs a set of mutations on the listed photos.
func Do(opts *Options) error {
	ctx := opts.Ctx
	mi := opts.ObjStore

	opts.Tags.Normalize()
	opts.Tags.Validate(opts.Cfg.ValidLabels())

	log.Printf("--add %s", opts.Tags.A)
	log.Printf("--remove %s", opts.Tags.B)

	objs := make(map[string][]byte)
	reader := bufio.NewReader(os.Stdin)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			if str != "" {
				mutate(mi, objs, str, &opts.Tags)
			}
			break
		}
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}

		// Strip newline off the end.
		mutate(mi, objs, str[:len(str)-1], &opts.Tags)
	}

	if !opts.DryRun {
		err := mi.InsertBulk(ctx, objs)
		if err != nil {
			log.Fatalf("Oops %v", err)
		}
	} else {
		log.Printf("Running in dry run mode, no changes made")
	}

	return nil
}

func mutate(mi *objectstore.ObjectStore, objs map[string][]byte, line string, tags *util.Tags) {
	var changed bool

	b, err := hex.DecodeString(line)
	if err != nil {
		log.Fatalf("Unable to decode sha256: %s", line)
	}

	value := mi.Get(hex.EncodeToString(b))

	var p1 mediainfo.Media
	if er := proto.Unmarshal(value, &p1); er != nil {
		log.Fatalf("Unmarshalling proto: %v", er)
	}

	p1.Tags, changed = tags.Modify(p1.Tags)
	if changed {
		data, err := proto.Marshal(&p1)
		if err != nil {
			log.Fatalf("Failed to marshal proto")
		}

		objs[hex.EncodeToString(b)] = data

		log.Printf("Modifying %q\n", p1.File[0].Filename)
	}
}
