package cache

import (
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/client"
	"github.com/nthnca/datastore"
)

func Handler() {
	clt, _ := datastore.NewCloudClient(config.ProjectID)
	client.CompactPhotoCache(clt)
}
