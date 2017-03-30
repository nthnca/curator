package main

import (
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/client"
	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/curator/util"
	"github.com/nthnca/datastore"
)

type ByLength []util.Data

func (s ByLength) Len() int {
	return len(s)
}
func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
	return s[i].Score < s[j].Score
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup

	var data []*message.Comparison
	wg.Add(1)
	go func() {
		clt, _ := datastore.NewCloudClient(config.ProjectID)
		data, _ = client.LoadAllComparisons(clt)
		wg.Done()
	}()

	wg.Wait()

	score := util.CalculateRankings(data)

	var xyz []util.Data
	for _, y := range score {
		xyz = append(xyz, y)
	}

	sort.Sort(ByLength(xyz))

	for _, y := range xyz {
		log.Printf("%v: (%v) %v", y.Key, y.Views, y.Score)
	}
	log.Printf("LEN: %v", len(xyz))
}
