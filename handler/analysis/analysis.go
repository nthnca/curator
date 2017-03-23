package analysis

import (
	"html/template"
	"net/http"
	"sort"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/client"
	"github.com/nthnca/curator/util"
	"google.golang.org/appengine"
)

const tpl = `<!DOCTYPE html>
<html><head>
  <meta charset="UTF-8">
</head><body>
<p>Matches: {{.Matches}}</p>
<p>Pairs: {{.Pairs}}</p>
<p>Images: {{len .Data}}</p>
{{range .Data}}
  <div><a href="{{$.Url}}/{{.Key}}.jpg">
    {{.Key}}</a> Views:{{.Views}} Score:{{.Score}} Next:{{.Next}}
  </div>
  <div>&nbsp;&nbsp;
  {{range .Games}}
    {{.Result}} {{(index $.Map .Opponent).Score}}
  {{end}}
  </div>
{{end}}
</body></html>`

type ByLength []util.Data

func (s ByLength) Len() int {
	return len(s)
}
func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
	return s[i].Score > s[j].Score
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	comparisons, _ := client.LoadAllComparisons2(ctx)
	score := util.CalculateRankings(comparisons)

	data := struct {
		Matches int
		Pairs   int
		Data    []util.Data
		Map     map[string]util.Data
		Url     string
	}{
		Matches: 2,
		Pairs:   len(comparisons),
		Url:     "https://storage.googleapis.com/" + config.StorageBucket,
		Map:     score,
	}

	for _, y := range score {
		data.Data = append(data.Data, y)
	}
	sort.Sort(ByLength(data.Data))

	t, _ := template.New("webpage").Parse(tpl)
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	t.Execute(w, data)
}
