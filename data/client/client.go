package client

import (
	"fmt"

	"github.com/nthnca/curator/data/message"

	"github.com/golang/protobuf/proto"
	"github.com/nthnca/datastore"
	"google.golang.org/api/iterator"
)

func Put(clt datastore.Client, key datastore.Key, p proto.Message) (datastore.Key, error) {
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

func Get(clt datastore.Client, key datastore.Key, p proto.Message) error {
	var entry Proto
	if err := clt.Get(key, &entry); err != nil {
		return fmt.Errorf("Failed to load entry: %v", err)
	}

	if err := proto.Unmarshal(entry.Proto, p); err != nil {
		return fmt.Errorf("Failed to deserialize: %v", err)
	}

	return nil
}

func LoadNextTada(clt datastore.Client) ([]*message.Photo, error) {
	q := clt.NewQuery("Tada") //.Limit(1)
	for it := clt.Run(q); ; {
		var entry Proto
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

func SaveComparison(clt datastore.Client, p *message.Comparison) error {
	key := clt.IncompleteKey("ComparisonSet", nil)
	_, err := Put(clt, key, p)
	if err != nil {
		return fmt.Errorf("Iterator failed: %v", err)
	}

	return nil
}

func LoadAllComparisons(clt datastore.Client) ([]*message.ComparisonEntry, error) {
	var rv []*message.ComparisonEntry

	q := clt.NewQuery("ComparisonSet")
	for it := clt.Run(q); ; {
		var entry Proto
		_, err := it.Next(&entry)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Iterator failed: %v", err)
		}

		p := &message.Comparison{}
		err = proto.Unmarshal(entry.Proto, p)
		if err != nil {
			return nil, fmt.Errorf("Unmarshalling error: %v", err)
		}

		rv = append(rv, p.GetEntry()...)
	}

	return rv, nil
}
