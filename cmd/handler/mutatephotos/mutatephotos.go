package mutatephotos

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/nthnca/curator/pkg/config"
	"github.com/nthnca/curator/pkg/mediainfo/store"
	"github.com/nthnca/curator/pkg/util"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	ctx       context.Context
	client    *storage.Client
	mediaInfo *store.MediaInfo
	tags      util.Tags
	actual    *bool
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
	ctx = context.Background()
	err := fmt.Errorf("") // Next line can't use :=
	client, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	mi, err := store.New(ctx, client, config.MediaInfoBucket())
	mediaInfo = mi
	if err != nil {
		log.Fatalf("New MediaInfo store failed: %v", err)
	}

	tags.Normalize()
	tags.Validate(config.ValidLabels())

	log.Printf("--add %s", tags.A)
	log.Printf("--remove %s", tags.B)

	reader := bufio.NewReader(os.Stdin)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			if str != "" {
				mutate(str)
			}
			break
		}
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}

		// Strip newline off the end.
		mutate(str[:len(str)-1])
	}

	if *actual {
		mi.Flush(ctx, client)
	} else {
		log.Printf("Running in dry run mode, no changes made")
	}
}

func mutate(line string) {
	var changed bool

	b, err := hex.DecodeString(line)
	if err != nil {
		log.Fatalf("Unable to decode sha256: %s", line)
	}
	p1 := mediaInfo.Get(util.Sha256(b))
	p1.Tags, changed = tags.Modify(p1.Tags)
	if changed {
		p1.TimestampSeconds = time.Now().Unix()
		mediaInfo.InsertFast(p1)

		log.Printf("Modifying %q\n", p1.File[0].Filename)
	}
}
