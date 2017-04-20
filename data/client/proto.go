package client

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/nthnca/datastore"
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
