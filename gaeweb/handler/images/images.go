package images

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/client"
	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/curator/util"
	"github.com/nthnca/datastore"

	humanize "github.com/dustin/go-humanize"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

type PhotoInfoJson struct {
	Date        string
	Time        string
	Camera      string
	X           string
	Y           string
	Bytes       string
	Aperture    string
	Shutter     string
	FocalLength string
	ISO         string
}

type PhotoJson struct {
	ID   string        `json:"id"`
	Src  string        `json:"src"`
	Info PhotoInfoJson `json:"info"`
}

type PhotoRequest struct {
	ID     string `json:"id"`
	Result int    `json:"result"`
	Angle  int    `json:"angle"`
	Flag   int    `json:"flag"`
}

func tester() []PhotoJson {
	return []PhotoJson{
		PhotoJson{ID: "a", Src: "http://i.imgur.com/xfrYauH.jpg", Info: PhotoInfoJson{ISO: "iso100"}},
		PhotoJson{ID: "b", Src: "http://i.imgur.com/Oci11T4.jpg", Info: PhotoInfoJson{ISO: "iso200"}},
		PhotoJson{ID: "c", Src: "http://i.imgur.com/Z9sNBeP.jpg", Info: PhotoInfoJson{ISO: "iso300"}},
		PhotoJson{ID: "d", Src: "http://i.imgur.com/402I0Z8.jpg", Info: PhotoInfoJson{ISO: "iso400"}},
	}
}

func ParseBody(ctx context.Context, r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("Couldn't parse request body: %s", err)
	}

	dec := json.NewDecoder(strings.NewReader(string(body)))
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("Couldn't decode post: %v", err)
	}

	return nil
}

func processImageResults(ctx context.Context, r *http.Request) {
	var ir []PhotoRequest
	if err := ParseBody(ctx, r, &ir); err != nil {
		log.Warningf(ctx, "Couldn't decode post: %v", err)
		return
	}

	clt := datastore.NewGaeClient(ctx)
	test := message.Comparison{
		Epoch: proto.Int64(time.Now().Unix())}
	for _, e1 := range ir {
		if e1.Flag == 1 {
			p, _ := client.GetPhoto(clt, e1.ID)
			p.UserHide = proto.Bool(true)
			client.UpdatePhoto(clt, e1.ID, &p)
		}
		for _, e2 := range ir {
			if e1.Result <= e2.Result {
				continue
			}
			test.Entry = append(test.Entry, &message.ComparisonEntry{
				Photo1: proto.String(e1.ID),
				Photo2: proto.String(e2.ID),
				Score:  proto.Int32(int32(e1.Result - e2.Result))})
		}
	}
	if len(test.Entry) > 0 {
		client.SaveComparison(clt, &test)
	}
}

func gcd(x, y int32) int32 {
	for y != 0 {
		x, y = y, x%y
	}
	if x == 0 {
		return 1
	}
	return x
}

func generateImageSet(ctx context.Context, w http.ResponseWriter) {
	clt := datastore.NewGaeClient(ctx)

	var d []PhotoJson
	list, _ := client.LoadNextTada(clt)
	for i := range list {
		n := list[i].GetKey()
		p, _ := client.GetPhoto(clt, n)
		if p.GetUserHide() {
			continue
		}
		prop := p.GetProperties()
		exp_n, exp_d := prop.GetExposureTime().GetNumerator(),
			prop.GetExposureTime().GetDenominator()
		exp_gcd := gcd(exp_n, exp_d)
		date := time.Unix(prop.GetOriginalEpoch(), 0).Format(time.RFC1123)
		date = strings.Replace(date, " 0", " ", -1)
		date = strings.Join(strings.Split(date, ":")[:2], ":")

		i := PhotoInfoJson{
			Date:   date,
			Camera: prop.GetMake() + " " + prop.GetModel(),
			X:      fmt.Sprintf("%d", prop.GetWidth()),
			Y:      fmt.Sprintf("%d", prop.GetHeight()),
			Bytes:  humanize.Bytes(uint64(p.GetBytes())),
			Aperture: fmt.Sprintf("f/%.1f",
				float64(prop.GetAperture().GetNumerator())/float64(prop.GetAperture().GetDenominator())),
			Shutter: fmt.Sprintf("%d/%d", exp_n/exp_gcd, exp_d/exp_gcd),
			FocalLength: fmt.Sprintf("%.1f mm",
				float64(prop.GetFocalLength().GetNumerator())/float64(prop.GetFocalLength().GetDenominator())),
			ISO: fmt.Sprintf("iso%d", prop.GetIso()),
		}
		f := PhotoJson{
			ID: n,
			Src: fmt.Sprintf(
				"https://storage.googleapis.com/%v/%v.jpg",
				config.StorageBucket, n),
			Info: i}
		d = append(d, f)
	}
	if util.IsDevAppServer() {
		log.Warningf(ctx, "IN TEST MODE")
		d = tester()
		time.Sleep(2000000000)
		log.Warningf(ctx, "DONE SLEEPING")
	}
	jData, _ := json.Marshal(d)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	defer r.Body.Close()

	processImageResults(ctx, r)
	generateImageSet(ctx, w)

	t := taskqueue.NewPOSTTask("/worker", url.Values{
		"key": {"key"},
	})
	taskqueue.Add(ctx, t, "") // add t to the default queue
}
