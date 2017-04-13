package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/disk"
	"github.com/nthnca/curator/handler/analyze"
	sy "github.com/nthnca/curator/handler/sync"
	"github.com/nthnca/curator/handler/update"
	"github.com/nthnca/curator/util"

	"github.com/nthnca/gobuild"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := kingpin.New(
		"curator",
		"Photo organizational system that run in Google AppEngine")
	gobuild.RegisterCommands(app, config.Path, config.ProjectID)
	app.Command("sync", "list entries").Action(sy.Handler)
	app.Command("analyze", "list entries").Action(analyze.Handler)
	app.Command("update", "list entries").Action(update.Handler)
	app.Command("test", "test").Action(test)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- proto.Message) {
	defer wg.Done()
	for j := range jobs {
		photo, err := util.IdentifyPhoto(j)
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}

		results <- photo
	}
}

func test(_ *kingpin.ParseContext) error {
	jobs := make(chan string, 100)
	results := make(chan proto.Message, 100)

	wg := &sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go worker(wg, jobs, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		m := disk.List(config.PhotoPath)
		for _, path := range m {
			jobs <- path

		}
		close(jobs)
	}()

	for r := range results {
		fmt.Printf("%v\n", r.String())
	}
	return nil
}
