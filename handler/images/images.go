package images

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/client"
	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/curator/util"
	"github.com/nthnca/datastore"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

func tester() []*message.Photo {
	return []*message.Photo{
		&message.Photo{Key: proto.String("ABC")},
		&message.Photo{Key: proto.String("CDE")},
		&message.Photo{Key: proto.String("ABD")},
		&message.Photo{Key: proto.String("Fake")},
		&message.Photo{Key: proto.String("FIxmE")}}
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
	var ir []struct {
		Image1, Image2 string
		Result         int32
	}

	if err := ParseBody(ctx, r, &ir); err != nil {
		log.Warningf(ctx, "Couldn't decode post: %v", err)
		return
	}

	if len(ir) == 0 {
		return
	}

	test := message.Comparison{
		Epoch: proto.Int64(time.Now().Unix())}
	for _, k := range ir {
		test.Entry = append(test.Entry, &message.ComparisonEntry{
			Photo1: proto.String(k.Image1),
			Photo2: proto.String(k.Image2),
			Score:  proto.Int32(k.Result)})
	}
	clt := datastore.NewGaeClient(ctx)
	client.SaveComparison(clt, &test)
}

func generateImageSet(ctx context.Context, w http.ResponseWriter) {
	clt := datastore.NewGaeClient(ctx)
	list, err := client.LoadNextTada(clt)
	if err != nil {
		log.Warningf(ctx, "%v", err)
		return
	}

	if util.IsDevAppServer() {
		log.Warningf(ctx, "IN TEST MODE")
		list = tester()
		time.Sleep(2000000000)
		log.Warningf(ctx, "DONE SLEEPING")
	}

	var d []map[string]string
	for i := range list {
		n := list[i].GetKey()
		e := map[string]string{
			"id": n,
			"src": fmt.Sprintf(
				"https://storage.googleapis.com/%v/%v.jpg",
				config.StorageBucket, n)}
		d = append(d, e)
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
