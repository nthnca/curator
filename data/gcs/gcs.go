package gcs

import (
	"log"
	"strings"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func List(bucket string) (map[string]string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	m := make(map[string]string)
	bkt := client.Bucket(bucket)
	for it := bkt.Objects(ctx, nil); ; {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}

		m[strings.SplitN(objAttrs.Name, ".", 2)[0]] = objAttrs.Name
	}
	return m, nil
}
