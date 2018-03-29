package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/nthnca/curator/cmd/handler/fsckphotos"
	"github.com/nthnca/curator/cmd/handler/getphotos"
	"github.com/nthnca/curator/cmd/handler/mutatephotos"
	"github.com/nthnca/curator/cmd/handler/newphotos"
	"github.com/nthnca/curator/cmd/handler/statphotos"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	actual bool
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := kingpin.New(
		"curator",
		"Photo organizational system that run in Google AppEngine")
	app.UsageWriter(os.Stdout)
	app.Flag("go", "Actually do things").BoolVar(&actual)

	newphotos.Register(app, &actual)
	getphotos.Register(app)
	mutatephotos.Register(app, &actual)
	statphotos.Register(app)
	fsckphotos.Register(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
