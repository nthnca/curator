package code

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/nthnca/curator/handler/analysis"
	"github.com/nthnca/curator/handler/images"
	"github.com/nthnca/curator/handler/index"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", index.Handler)
	http.HandleFunc("/images", images.Handler)
	http.HandleFunc("/best", analysis.Handler)
}
