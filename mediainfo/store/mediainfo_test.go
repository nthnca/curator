package store

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/curator/util"
	"google.golang.org/api/iterator"
)

func TestMediaInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	bucketName := "test-bucket-mediainfo"

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	cleanup(ctx, client, bucketName)
	time.Sleep(2 * time.Second)

	baseCapabilities(ctx, client, bucketName)
}

func cleanup(ctx context.Context, client *storage.Client, bucketName string) {
	log.Printf("cleanup")
	bucket := client.Bucket(bucketName)
	for it := bucket.Objects(ctx, nil); ; {
		obj, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			log.Fatalf("Failed to iterate through objects: %v", err)
		}

		err = client.Bucket(bucketName).Object(obj.Name).Delete(ctx)
		if err != nil {
			log.Fatalf("Failed to delete object: %v (%s)", err, obj.Name)
		}
	}
}
func baseCapabilities(ctx context.Context, client *storage.Client, bucketName string) {
	insert := func(mi *MediaInfo, time int64, count int) {
		log.Printf("Insert time: %d count: %d", time, count)
		for i := 0; i < count; i++ {
			var photo message.Media
			shasum := sha256.Sum256([]byte(fmt.Sprintf("%v", i)))
			photo.Key = shasum[:]
			photo.TimestampSecondsSinceEpoch = time
			mi.Insert(ctx, client, &photo)
		}
	}

	check := func(mi *MediaInfo, time int64, count int) {
		log.Printf("Check time: %d count: %d", time, count)
		index := make(map[[32]byte]int)
		for _, e := range mi.All() {
			index[util.Sha256(e.Key)] = 1
			if e.TimestampSecondsSinceEpoch != time {
				log.Fatalf("Wrong entry value: %d", e.TimestampSecondsSinceEpoch)

			}
		}
		if len(index) != count {
			log.Fatalf("Wrong number of entries: %d (%d)", len(mi.All()), len(index))
		}
	}

	log.Printf("baseCapabilities")
	{
		mi, err := New(ctx, client, bucketName)
		if err != nil {
			log.Fatalf("NewMediaInfo failed: %v", err)
		}

		insert(mi, 1, 3)
		check(mi, 1, 3)
		insert(mi, 2, 3)
		check(mi, 2, 3)
		insert(mi, 1, 3)
		check(mi, 2, 3)
	}

	// Note that at this point we don't even have a master file.
	time.Sleep(2 * time.Second)
	log.Printf("reload")
	mi, err := New(ctx, client, bucketName)
	if err != nil {
		log.Fatalf("NewMediaInfo failed: %v", err)
	}

	check(mi, 2, 3)

	t := time.Now().UnixNano()
	insert(mi, 3, 50)
	mi.Flush(ctx, client)
	log.Printf("Adding 50 entries took %v seconds",
		float64(time.Now().UnixNano()-t)/1000000000.0)

	mi, err = New(ctx, client, bucketName)
	if err != nil {
		log.Fatalf("NewMediaInfo failed: %v", err)
	}
	check(mi, 3, 50)
}
