package client

import (
	"fmt"

	"github.com/nthnca/curator/data/message"

	"github.com/golang/protobuf/proto"
	"github.com/nthnca/datastore"
	"google.golang.org/api/iterator"
)

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
	data, err := proto.Marshal(p)
	if err != nil {
		return fmt.Errorf("Iterator failed: %v", err)
	}
	entry := Proto{Proto: data}
	key := clt.IncompleteKey("ComparisonSet", nil)
	if _, err := clt.Put(key, &entry); err != nil {
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
