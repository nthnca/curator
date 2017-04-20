package client

import (
	"fmt"

	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/curator/util/need"

	"github.com/golang/protobuf/proto"
	"github.com/nthnca/datastore"
	"google.golang.org/api/iterator"
)

func put(clt datastore.Client, key datastore.Key, p proto.Message) (datastore.Key, error) {
	data, err := proto.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize proto: %v", err)
	}

	k, err := clt.Put(key, &Proto{Proto: data})
	if err != nil {
		return nil, fmt.Errorf("Failed to save entry: %v", err)
	}

	return k, nil
}

func get(clt datastore.Client, key datastore.Key, p proto.Message) error {
	var entry Proto
	if err := clt.Get(key, &entry); err != nil {
		return fmt.Errorf("Failed to load entry: %v", err)
	}

	if err := proto.Unmarshal(entry.Proto, p); err != nil {
		return fmt.Errorf("Failed to deserialize: %v", err)
	}

	return nil
}

func next(iter datastore.Iterator, p proto.Message) (datastore.Key, error) {
	var entry Proto
	key, err := iter.Next(&entry)
	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(entry.Proto, p)
	if err != nil {
		return nil, fmt.Errorf("Unmarshalling error: %v", err)
	}
	return key, nil
}

func CreateTada(clt datastore.Client, msg *message.PhotoSet) (datastore.Key, error) {
	return put(clt, clt.IncompleteKey("Tada", nil), msg)
}

func LoadNextTada(clt datastore.Client) ([]*message.Photo, error) {
	q := clt.NewQuery("Tada") //.Limit(1)
	for it := clt.Run(q); ; {
		rv := &message.PhotoSet{}
		k, err := next(it, rv)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Iterator failed: %v", err)
		}

		// Delete this entry
		clt.Delete(k)
		return rv.GetPhoto(), nil
	}

	return nil, fmt.Errorf("No results found.")
}

func SaveComparison(clt datastore.Client, p *message.Comparison) error {
	key := clt.IncompleteKey("ComparisonSet", nil)
	_, err := put(clt, key, p)
	if err != nil {
		return fmt.Errorf("Iterator failed: %v", err)
	}

	return nil
}

func LoadAllComparisons(clt datastore.Client) ([]*message.ComparisonEntry, error) {
	var rv []*message.ComparisonEntry

	q := clt.NewQuery("ComparisonSet")
	for it := clt.Run(q); ; {
		p := &message.Comparison{}
		_, err := next(it, p)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Iterator failed: %v", err)
		}

		rv = append(rv, p.GetEntry()...)
	}

	return rv, nil
}

func CreatePhotoCache(clt datastore.Client, msg *message.PhotoSet) (datastore.Key, error) {
	return put(clt, clt.IncompleteKey("PhotoCache", nil), msg)
}

func ReadAllPhotoCache(clt datastore.Client) ([]*message.Photo, error) {
	var rv []*message.Photo

	q := clt.NewQuery("PhotoCache")
	for it := clt.Run(q); ; {
		p := &message.PhotoSet{}
		_, err := next(it, p)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Iterator failed: %v", err)
		}

		rv = append(rv, p.GetPhoto()...)
	}

	return rv, nil
}

var needPhotoCacheData need.NeedData

func NeedPhotoCache(clt datastore.Client) func() []*message.Photo {
	n := needPhotoCacheData.Need(func() interface{} {
		mp, _ := ReadAllPhotoCache(clt)
		return mp
	})
	return func() []*message.Photo {
		return n().([]*message.Photo)
	}
}
