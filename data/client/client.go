package client

import (
	"fmt"

	"github.com/nthnca/curator/data/message"

	"github.com/nthnca/datastore"
	"google.golang.org/api/iterator"
)

func CreateQueue(clt datastore.Client, msg *message.PhotoSet) (datastore.Key, error) {
	return put(clt, clt.IncompleteKey("Queue", nil), msg)
}

func LoadNextQueue(clt datastore.Client) ([]*message.Photo, error) {
	q := clt.NewQuery("Queue").Limit(1)
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

func ClearQueue(clt datastore.Client) error {
	q := clt.NewQuery("Queue")
	for it := clt.Run(q); ; {
		rv := &message.PhotoSet{}
		k, err := next(it, rv)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Iterator failed: %v", err)
		}

		// Delete this entry
		clt.Delete(k)
	}

	return nil
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
