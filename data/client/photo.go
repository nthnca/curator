package client

import (
	"fmt"

	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/curator/util/need"

	"github.com/nthnca/datastore"
	"google.golang.org/api/iterator"
)

func GetPhoto(clt datastore.Client, key string) (message.Photo, error) {
	var p message.Photo
	err := get(clt, clt.NameKey("Photo", key, nil), &p)
	return p, err
}

func UpdatePhoto(clt datastore.Client, key string, msg *message.Photo) error {
	_, err := put(clt, clt.NameKey("Photo", key, nil), msg)
	return err
}

func GetPhotos(clt datastore.Client) ([]*message.Photo, error) {
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

func NeedPhotos(clt datastore.Client) func() []*message.Photo {
	n := needPhotoCacheData.Need(func() interface{} {
		mp, _ := GetPhotos(clt)
		return mp
	})
	return func() []*message.Photo {
		return n().([]*message.Photo)
	}
}

func CompactPhotoCache(clt datastore.Client) error {
	clearPhotoCache(clt)

	photos, _ := getAllPhotosExpensive(clt)
	set := message.PhotoSet{}
	for _, photo := range photos {
		set.Photo = append(set.Photo, photo)

		if len(set.Photo) > 1000 {
			_, err := createPhotoCache(clt, &set)
			if err != nil {
				return fmt.Errorf("CreatePhotoCache: %v", err)
			}
			set = message.PhotoSet{}
		}
	}

	_, err := createPhotoCache(clt, &set)
	if err != nil {
		return fmt.Errorf("CreatePhotoCache: %v", err)
	}

	return nil
}

func getAllPhotosExpensive(clt datastore.Client) ([]*message.Photo, error) {
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

func createPhotoCache(clt datastore.Client, msg *message.PhotoSet) (datastore.Key, error) {
	return put(clt, clt.IncompleteKey("PhotoCache", nil), msg)
}

func clearPhotoCache(clt datastore.Client) error {
	q := clt.NewQuery("PhotoCache")
	for it := clt.Run(q); ; {
		p := &message.PhotoSet{}
		k, err := next(it, p)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Iterator failed: %v", err)
		}

		clt.Delete(k)
	}

	return nil
}
