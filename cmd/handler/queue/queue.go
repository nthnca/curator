package queue

import (
	"log"
	"sort"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/client"
	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/datastore"
)

func Handler() {
	clt, err := datastore.NewCloudClient(config.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	photos, _ := client.GetPhotos(clt)
	epoch := func(x int) int64 {
		return photos[x].GetProperties().GetOriginalEpoch()
	}
	sort.Slice(photos, func(i, j int) bool { return epoch(i) < epoch(j) })

	client.ClearQueue(clt)
	start := 0
	count := 0
	for i := range photos {
		if i+1 >= len(photos) || epoch(i)+30*60 < epoch(i+1) {
			if i-start >= 40 {
				log.Printf("%d %d", start, i)
				if count > 30 {
					return
				}
				count++
				SavePhotoSet(clt, photos[start:i])
			}
			start = i + 1
			continue
		}

	}
}

func SavePhotoSet(clt datastore.Client, photos []*message.Photo) {
	queue := message.PhotoSet{}
	queue.Photo = photos

	if _, err := client.CreateQueue(clt, &queue); err != nil {
		log.Fatalf("Failed to save: %v", err)
	}
}

//comparisons, _ := client.LoadAllComparisons(clt)
//rankings := util.CalculateRankings(comparisons)

//	for k := range photos {
//		if _, ok := rankings[k]; ok {
//			continue
//		}
//		rankings[k] = util.Data{Key: k, Score: 1500}
//	}
//
//	var arr []util.Data
//	for _, e := range photos {
//		if e.Views == 0 {
//			arr = append(arr, e)
//		}
//		/*
//			if e.Views < 7 && e.Score > 4000 {
//				arr = append(arr, e)
//			}
//		*/
//	}
//
//	if len(arr) < 1000 {
//		log.Fatalf("Not enough photos")
//	}
//
//	log.Printf("Total photos: %v", len(photoList))
//	log.Printf("Total compar: %v", len(comparisons))
//	log.Printf("Total length: %v", len(arr))

/*
	// Randomize arr
	for i := range arr {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
*/
