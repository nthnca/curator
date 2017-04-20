package client

import (
	"fmt"

	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/datastore"
	"google.golang.org/api/iterator"
)

// PutQueue saves a PhotoSet message to the datastore.
func PutQueue(clt datastore.Client, msg *message.PhotoSet) (datastore.Key, error) {
	return put(clt, clt.IncompleteKey("Queue", nil), msg)
}

// GetNextQueue gets and deletes a PhotoSet message from the datastore and then
// returns that Photo message.
func GetNextQueue(clt datastore.Client) ([]*message.Photo, error) {
	q := clt.NewQuery("Queue").Limit(1)
	for it := clt.Run(q); ; {
		rv := &message.PhotoSet{}
		k, err := next(it, rv)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		clt.Delete(k)
		return rv.GetPhoto(), nil
	}

	return nil, fmt.Errorf("no results found")
}

// ClearQueue deletes all the entries in the Queue datastore table.
func ClearQueue(clt datastore.Client) error {
	q := clt.NewQuery("Queue")
	for it := clt.Run(q); ; {
		rv := &message.PhotoSet{}
		k, err := next(it, rv)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		clt.Delete(k)
	}

	return nil
}
