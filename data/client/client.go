package client

import (
	"fmt"
	"log"
	"strings"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/entity"
	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/datastore"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func LoadNextTada(clt datastore.Client) ([]*message.Photo, error) {
	q := clt.NewQuery("Tada") //.Limit(1)
	for it := clt.Run(q); ; {
		var entry entity.Comparison
		k, err := it.Next(&entry)
		if err == iterator.Done {
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
		clt.Delete(k)
		return rv.GetPhoto(), nil
	}

	return nil, fmt.Errorf("No results found.")
}

func SaveComparison(clt datastore.Client, p *message.ComparisonSet) error {
	data, err := proto.Marshal(p)
	if err != nil {
		return fmt.Errorf("Iterator failed: %v", err)
	}
	entry := entity.Comparison{Proto: data}
	key := clt.IncompleteKey("ComparisonSet")
	if _, err := clt.Put(key, &entry); err != nil {
		return fmt.Errorf("Iterator failed: %v", err)
	}
	return nil
}

func LoadAllComparisons(clt datastore.Client) ([]*message.Comparison, error) {
	var rv []*message.Comparison

	q := clt.NewQuery("ComparisonSet")
	for it := clt.Run(q); ; {
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

func LoadAllPhotoSets(clt datastore.Client) ([]*message.Photo, error) {
	var rv []*message.Photo

	q := clt.NewQuery("PhotoSet")
	for it := clt.Run(q); ; {
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
