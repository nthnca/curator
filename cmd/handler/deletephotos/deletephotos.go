package deletephotos

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
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/mediainfo/store"
	"github.com/nthnca/curator/util"
)

const (
	dryRun = false
)

var (
	ctx       context.Context
	client    *storage.Client
	mediaInfo *store.MediaInfo
)

func Handler() {
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

	reader := bufio.NewReader(os.Stdin)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			if str != "" {
				Foo(str)
			}
			break
		}
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}

		// Strip newline off the end.
		Foo(str[:len(str)-1])
	}
	mi.Flush(ctx, client)
}

func Foo(line string) {
	b, _ := hex.DecodeString(line)
	p1 := mediaInfo.Get(util.Sha256(b))
	if p1.Deleted {
		log.Printf("Already deleted %q\n", p1.File[0].Filename)
	} else {
		log.Printf("Delete %q\n", p1.File[0].Filename)
		p1.TimestampSecondsSinceEpoch = time.Now().Unix()
		p1.Deleted = true
		mediaInfo.Insert(ctx, client, p1)
	}
}
