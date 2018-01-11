package mutatephotos

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/mediainfo/store"
	"github.com/nthnca/curator/util"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	ctx       context.Context
	client    *storage.Client
	mediaInfo *store.MediaInfo
	actual    *bool
	add       []string
	remove    []string
)

func Register(app *kingpin.Application, actualx *bool) {
	actual = actualx
	cmd := app.Command("mutate", "Mutate")
	cmd.Action(
		func(_ *kingpin.ParseContext) error {
			handler()
			return nil
		})
	cmd.Flag("add", "Labels to add").Short('a').StringsVar(&add)
	cmd.Flag("remove", "Labels to remove").Short('r').StringsVar(&remove)
}

func validateLabels(add, remove []string) {
	allowed := make(map[string]bool)

	for _, t := range config.ValidLabels() {
		allowed[t] = true
	}

	valid := func(arr []string) {
		for _, x := range arr {
			if !allowed[x] {
				log.Fatalf("Invalid tag: %s", x)
			}
		}
	}

	valid(add)
	valid(remove)

	duplicates := map[string]bool{}

	dups := func(arr []string) {
		for _, x := range arr {
			if duplicates[x] {
				log.Fatalf("Duplicate tag: %s", x)
			}
			duplicates[x] = true
		}
	}

	dups(add)
	dups(remove)
}

func modifyTags(tags, add, remove []string) ([]string, bool) {
	changed := false
	tagMap := make(map[string]bool)

	for _, t := range tags {
		tagMap[t] = true
	}
	for _, t := range add {
		_, has := tagMap[t]
		changed = changed || !has
		tagMap[t] = true
	}
	for _, t := range remove {
		_, has := tagMap[t]
		changed = changed || has
		delete(tagMap, t)
	}
	var rv []string
	for t, _ := range tagMap {
		rv = append(rv, t)
	}
	return rv, changed
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

	normalize := func(arr []string) []string {
		var tmp []string
		for _, t := range arr {
			tmp = append(tmp, strings.Split(t, ",")...)
		}
		return tmp
	}

	add = normalize(add)
	remove = normalize(remove)

	validateLabels(add, remove)

	log.Printf("--add %s", add)
	log.Printf("--remove %s", remove)

	reader := bufio.NewReader(os.Stdin)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			if str != "" {
				mutate(str, add, remove)
			}
			break
		}
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}

		// Strip newline off the end.
		mutate(str[:len(str)-1], add, remove)
	}

	if *actual {
		mi.Flush(ctx, client)
	}
}

func mutate(line string, add, remove []string) {
	b, _ := hex.DecodeString(line)
	p1 := mediaInfo.Get(util.Sha256(b))
	tags, changed := modifyTags(p1.Tags, add, remove)
	if changed {
		p1.TimestampSecondsSinceEpoch = time.Now().Unix()
		p1.Tags = tags
		log.Printf("Success %q\n", p1)

		mediaInfo.InsertFast(p1)
	} else {
		log.Printf("Nope %q\n", p1.File[0].Filename)
	}
}
