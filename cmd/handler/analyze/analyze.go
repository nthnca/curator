package analyze

import (
	"log"
	"sort"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/client"
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

func Handler() {
	clt, err := datastore.NewCloudClient(config.ProjectID)
	if err != nil {
		log.Printf("Creating cloud client failed: %v", err)
	}

	comparisons, err := client.LoadAllComparisons(clt)
	if err != nil {
		log.Printf("Failed to load comparisons: %v", err)
	}

	score := util.CalculateRankings(comparisons)

	var data []util.Data
	for _, y := range score {
		data = append(data, y)
	}

	sort.Sort(ByLength(data))

	for _, y := range data {
		log.Printf("%v: (%v) %v", y.Key, y.Views, y.Score)
	}
	log.Printf("LEN: %v", len(data))
}
