package client

import (
	"fmt"
	"log"
	"strings"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/entity"
	"github.com/nthnca/curator/data/message"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	ds "google.golang.org/appengine/datastore"
)

func LoadPhotos(ctx context.Context) ([]*message.Photo, error) {
	var entry entity.Photo
	key := ds.NewKey(ctx, "PhotoSet", "0", 0, nil)
	if err := ds.Get(ctx, key, &entry); err != nil {
		return nil, fmt.Errorf("Iterator failed: %v", err)
	}

	p := &message.PhotoSet{}
	err := proto.Unmarshal(entry.Proto, p)
	if err != nil {
		return nil, fmt.Errorf("Iterator failed: %v", err)
	}

	return p.GetPhoto(), nil
}

func LoadNextTada(ctx context.Context) ([]*message.Photo, error) {
	q := ds.NewQuery("Tada").Limit(1)
	for it := q.Run(ctx); ; {
		var entry entity.Comparison
		k, err := it.Next(&entry)
		if err == ds.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Iterator failed: %v", err)
		}

		rv := &message.PhotoSet{}
		err = proto.Unmarshal(entry.Proto, rv)
		if err != nil {
			return nil, fmt.Errorf("Unmarshalling error: %v", err)
		}

		// Delete this entry
		ds.Delete(ctx, k)
		return rv.GetPhoto(), nil
	}

	return nil, fmt.Errorf("No results found.")
}

func SaveComparison(ctx context.Context, p *message.ComparisonSet) error {
	data, err := proto.Marshal(p)
	if err != nil {
		return fmt.Errorf("Iterator failed: %v", err)
	}
	entry := entity.Comparison{Proto: data}
	key := ds.NewKey(ctx, "ComparisonSet", "", 0, nil)
	if _, err := ds.Put(ctx, key, &entry); err != nil {
		return fmt.Errorf("Iterator failed: %v", err)
	}
	return nil
}

func LoadAllComparisons2(ctx context.Context) ([]*message.Comparison, error) {
	var rv []*message.Comparison

	q := ds.NewQuery("ComparisonSet")
	for it := q.Run(ctx); ; {
		var entry entity.Comparison
		_, err := it.Next(&entry)
		if err == ds.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Iterator failed: %v", err)
		}

		p := &message.ComparisonSet{}
		err = proto.Unmarshal(entry.Proto, p)
		if err != nil {
			return nil, fmt.Errorf("Unmarshalling error: %v", err)
		}

		rv = append(rv, p.GetComparison()...)
	}

	return rv, nil
}

func LoadAllPhotos() ([]message.Photo, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	var rv []message.Photo
	for it := client.Bucket(config.StorageBucket).Objects(ctx, nil); ; {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// It seems there is a bug and we get one error
			// before iterator done.
			continue
			// return nil, fmt.Errorf("Iterator failed: %v", err)
		}

		rv = append(rv, message.Photo{
			Name: proto.String(
				strings.SplitN(objAttrs.Name, ".", 2)[0])})
	}
	return rv, nil
}

func LoadAllComparisons() ([]*message.Comparison, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	var rv []*message.Comparison

	q := datastore.NewQuery("ComparisonSet")
	for it := client.Run(ctx, q); ; {
		var entry entity.Comparison
		_, err := it.Next(&entry)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Iterator failed: %v", err)
		}

		p := &message.ComparisonSet{}
		err = proto.Unmarshal(entry.Proto, p)
		if err != nil {
			return nil, fmt.Errorf("Unmarshalling error: %v", err)
		}

		rv = append(rv, p.GetComparison()...)
	}

	return rv, nil
}

func LoadAllPhotoSets(ctx context.Context) ([]*message.Photo, error) {
	client, err := datastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	var rv []*message.Photo

	q := datastore.NewQuery("PhotoSet")
	for it := client.Run(ctx, q); ; {
		var entry entity.Photo
		_, err := it.Next(&entry)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Iterator failed: %v", err)
		}

		p := &message.PhotoSet{}
		err = proto.Unmarshal(entry.Proto, p)
		if err != nil {
			return nil, fmt.Errorf("Unmarshalling error: %v", err)
		}

		rv = append(rv, p.GetPhoto()...)
	}

	return rv, nil
}
