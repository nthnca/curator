package client

import (
	"fmt"

	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/curator/util/need"

	"github.com/nthnca/datastore"
	"google.golang.org/api/iterator"
)

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

func UpdatePhoto(clt datastore.Client, key string, msg *message.Photo) error {
	_, err := put(clt, clt.NameKey("Photo", key, nil), msg)
	return err
}

func GetPhoto(clt datastore.Client, key string) (message.Photo, error) {
	var p message.Photo
	err := get(clt, clt.NameKey("Photo", key, nil), &p)
	return p, err
}

func ReadAllPhoto(clt datastore.Client) ([]*message.Photo, error) {
	var rv []*message.Photo

	q := clt.NewQuery("Photo")
	for it := clt.Run(q); ; {
		p := &message.Photo{}
		_, err := next(it, p)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Iterator failed: %v", err)
		}

		rv = append(rv, p)
	}

	return rv, nil
}

var needPhotoState need.NeedData

func NeedPhoto(clt datastore.Client) func() []*message.Photo {
	n := needPhotoState.Need(func() interface{} {
		mp, _ := ReadAllPhoto(clt)
		return mp
	})
	return func() []*message.Photo {
		return n().([]*message.Photo)
	}
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
