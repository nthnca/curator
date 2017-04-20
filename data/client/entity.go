package client

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/nthnca/datastore"
)

type protoEntry struct {
	Proto []byte `datastore:",noindex"`
}

func put(clt datastore.Client, key datastore.Key, p proto.Message) (datastore.Key, error) {
	data, err := proto.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("serializing proto: %v", err)
	}

	k, err := clt.Put(key, &protoEntry{Proto: data})
	if err != nil {
		return nil, fmt.Errorf("saving entry: %v", err)
	}
	return k, nil
}

func get(clt datastore.Client, key datastore.Key, p proto.Message) error {
	var entry protoEntry
	if err := clt.Get(key, &entry); err != nil {
		return fmt.Errorf("getting entry: %v", err)
	}

	if err := proto.Unmarshal(entry.Proto, p); err != nil {
		return fmt.Errorf("deserializing proto: %v", err)
	}
	return nil
}

func next(iter datastore.Iterator, p proto.Message) (datastore.Key, error) {
	var entry protoEntry
	key, err := iter.Next(&entry)
	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(entry.Proto, p)
	if err != nil {
		return nil, fmt.Errorf("deserializing proto: %v", err)
	}
	return key, nil
}
