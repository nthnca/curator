package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/client"
	"github.com/nthnca/curator/data/entity"
	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/curator/util"
	"github.com/nthnca/datastore"

	"github.com/golang/protobuf/proto"
)

func SavePhotoSet2(photos *message.PhotoSet) {
	clt, err := datastore.NewCloudClient(config.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	serialized, err := proto.Marshal(photos)
	if err != nil {
		// log.Infof("Marshaling failed: %v", err)
		return
	}

	key := clt.IncompleteKey("Tada")
	entity := entity.Comparison{Proto: serialized}

	if _, err := clt.Put(key, &entity); err != nil {
		log.Fatalf("Failed to save: %v", err)
	}
}

func SavePhotoSet(count int, photos *message.PhotoSet) {
	clt, err := datastore.NewCloudClient(config.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	serialized, err := proto.Marshal(photos)
	if err != nil {
		// log.Infof("Marshaling failed: %v", err)
		return
	}

	key := clt.NameKey("PhotoSet", fmt.Sprintf("%v", count))
	entity := entity.Photo{Proto: serialized}

	if _, err := clt.Put(key, &entity); err != nil {
		log.Fatalf("Failed to save: %v", err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	photoList, err := client.LoadAllPhotos()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	clt, err := datastore.NewCloudClient(config.ProjectID)
	comparisons, _ := client.LoadAllComparisons(clt)

	photos := util.CalculateRankings(comparisons)
	for _, e := range photoList {
		if _, ok := photos[e.GetName()]; ok {
			continue
		}
		photos[e.GetName()] = util.Data{Key: e.GetName(), Score: 1500}
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
				Name: proto.String(arr[j].Key)})
			j++
		}
		SavePhotoSet2(&msg)
	}
}
