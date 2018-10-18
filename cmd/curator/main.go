package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/nthnca/curator/pkg/action/fsckphotos"
	"github.com/nthnca/curator/pkg/action/getphotos"
	"github.com/nthnca/curator/pkg/action/mutatephotos"
	"github.com/nthnca/curator/pkg/action/newphotos"
	"github.com/nthnca/curator/pkg/action/statphotos"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	actual bool
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := kingpin.New(
		"curator",
		"Photo storage and organization tool")
	app.UsageWriter(os.Stdout)
	app.Flag("go", "Actually make modifications").BoolVar(&actual)

	newphotos.Register(app, &actual)
	getphotos.Register(app)
	mutatephotos.Register(app, &actual)
	statphotos.Register(app)
	fsckphotos.Register(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
