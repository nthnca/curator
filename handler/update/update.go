package update

import (
	"log"
	"math/rand"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/client"
	"github.com/nthnca/curator/data/gcs"
	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/curator/util"

	"github.com/golang/protobuf/proto"
	"github.com/nthnca/datastore"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func SavePhotoSet(photos *message.PhotoSet) {
	clt, err := datastore.NewCloudClient(config.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	serialized, err := proto.Marshal(photos)
	if err != nil {
		// log.Infof("Marshaling failed: %v", err)
		return
	}

	key := clt.IncompleteKey("Tada", nil)
	entity := client.Proto{Proto: serialized}

	if _, err := clt.Put(key, &entity); err != nil {
		log.Fatalf("Failed to save: %v", err)
	}
}

func Handler(_ *kingpin.ParseContext) error {
	photoList, err := gcs.List(config.StorageBucket)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	clt, err := datastore.NewCloudClient(config.ProjectID)
	comparisons, _ := client.LoadAllComparisons(clt)

	photos := util.CalculateRankings(comparisons)
	for k := range photoList {
		if _, ok := photos[k]; ok {
			continue
		}
		photos[k] = util.Data{Key: k, Score: 1500}
	}

	var arr []util.Data
	for _, e := range photos {
		if e.Views == 0 {
			arr = append(arr, e)
		}
		/*
			if e.Views < 7 && e.Score > 4000 {
				arr = append(arr, e)
			}
		*/
	}

	if len(arr) < 1000 {
		log.Fatalf("Not enough photos")
	}
	log.Printf("Total photos: %v", len(photoList))
	log.Printf("Total compar: %v", len(comparisons))
	log.Printf("Total length: %v", len(arr))

	// Randomize arr
	for i := range arr {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}

	j := 0
	for i := 0; i < len(arr)/7; i++ {
		msg := message.PhotoSet{}
		for len(msg.Photo) < 7 {
			msg.Photo = append(msg.Photo, &message.Photo{
				Key: proto.String(arr[j].Key)})
			j++
		}
		SavePhotoSet(&msg)
	}
	return nil
}
