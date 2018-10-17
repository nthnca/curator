package mutatephotos

import (
	"bufio"
	"context"
	"encoding/hex"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/pkg/config"
	"github.com/nthnca/curator/pkg/mediainfo/message"
	"github.com/nthnca/curator/pkg/util"
	objectstore "github.com/nthnca/object-store"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	tags   util.Tags
	actual *bool
)

func Register(app *kingpin.Application, actualx *bool) {
	actual = actualx
	cmd := app.Command("mutate", "Mutate")
	cmd.Action(
		func(_ *kingpin.ParseContext) error {
			handler()
			return nil
		})
	cmd.Flag("add", "Labels to add").Short('a').StringsVar(&tags.A)
	cmd.Flag("remove", "Labels to remove").Short('r').StringsVar(&tags.B)
}

func handler() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	mi, err := objectstore.New(ctx, client, config.MetadataBucket(), config.MetadataPath())
	if err != nil {
		log.Fatalf("New ObjectStore failed: %v", err)
	}

	tags.Normalize()
	tags.Validate(config.ValidLabels())

	log.Printf("--add %s", tags.A)
	log.Printf("--remove %s", tags.B)

	objs := make(map[string][]byte)
	reader := bufio.NewReader(os.Stdin)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			if str != "" {
				mutate(mi, objs, str)
			}
			break
		}
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}

		// Strip newline off the end.
		mutate(mi, objs, str[:len(str)-1])
	}

	if *actual {
		err := mi.InsertBulk(ctx, objs)
		if err != nil {
			log.Fatalf("Oops %v", err)
		}
	} else {
		log.Printf("Running in dry run mode, no changes made")
	}
}

func mutate(mi *objectstore.ObjectStore, objs map[string][]byte, line string) {
	var changed bool

	b, err := hex.DecodeString(line)
	if err != nil {
		log.Fatalf("Unable to decode sha256: %s", line)
	}

	value := mi.Get(hex.EncodeToString(b))

	var p1 message.Media
	if er := proto.Unmarshal(value, &p1); er != nil {
		log.Fatalf("Unmarshalling proto: %v", er)
	}

	p1.Tags, changed = tags.Modify(p1.Tags)
	if changed {
		p1.TimestampSeconds = time.Now().Unix()

		data, err := proto.Marshal(&p1)
		if err != nil {
			log.Fatalf("Failed to marshal proto")
		}

		objs[hex.EncodeToString(b)] = data

		log.Printf("Modifying %q\n", p1.File[0].Filename)
	}
}
