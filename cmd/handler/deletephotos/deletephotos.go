package deletephotos

import (
	"bufio"
	"context"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
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
	mediaInfo *store.MediaInfo
)

func Handler() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	mi, err := store.New(ctx, client, config.MediaInfoBucket())
	mediaInfo = mi
	if err != nil {
		log.Fatalf("New MediaInfo store failed: %v", err)
	}

	var b Blah
	reader := bufio.NewReader(os.Stdin)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			if str != "" {
				b.Foo(str)
			}
			break
		}
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		b.Foo(str[:len(str)-1])
	}
}

type Blah struct {
	inodes []string
}

func (b *Blah) Foo(line string) {
	parts := strings.Split(line, " ")
	if len(parts) != 2 {
		return
	}

	if len(parts[1]) > 6 && parts[1][0:6] == ".pics/" {
		for _, x := range b.inodes {
			if x == parts[0] {
				return
			}
		}
		b, _ := hex.DecodeString(parts[1][6:])
		p1 := mediaInfo.Get(util.Sha256(b))
		log.Printf("Delete %q\n", p1.File[0].Filename)
		p1.TimestampSecondsSinceEpoch = time.Now().Unix()
		p1.Deleted = true
		// log.Printf("P1 %+v", p1)
	} else {
		b.inodes = append(b.inodes, parts[0])
	}
}
