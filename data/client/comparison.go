package client

import (
	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/datastore"
	"google.golang.org/api/iterator"
)

// PutComparison saves a Comparison message to the datastore.
func PutComparison(clt datastore.Client, p *message.Comparison) error {
	key := clt.IncompleteKey("ComparisonSet", nil)
	_, err := put(clt, key, p)
	return err
}

// GetComparisons loads all Comparison messages from the datastore.
func GetComparisons(clt datastore.Client) ([]*message.ComparisonEntry, error) {
	var rv []*message.ComparisonEntry

	q := clt.NewQuery("ComparisonSet")
	for it := clt.Run(q); ; {
		p := &message.Comparison{}
		_, err := next(it, p)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		rv = append(rv, p.GetEntry()...)
	}

	return rv, nil
}
