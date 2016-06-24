package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/nleof/goyesql"
	"golang.org/x/net/context"
	"gopkg.in/unrolled/render.v1"
)

type JbConfig struct {
	Dbh    *sqlx.DB
	Query  goyesql.Queries
	Render *render.Render
}

func (c *JbConfig) GetHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

}

// Serves POST to /configurations
// Expected JSON structure
//      {
//			"name": "dictybase", # required
//			"general": {
//				"dataRoot": "data",
//				"datasets.DictyosteliumDiscoideumAX4": {
//					"url": "?data=genomes/2",
//					"name": "Dictyostelium Discoideum AX4"
//				},
//				"datasets.DictyosteliumFasciculatumSH3": {
//					"url": "?data=genomes/4",
//					"name": "Dictyostelium Fasciculatum SH3"
//				}
//          }
//		}
func (c *JbConfig) CreateHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Make sure the request body can be read twice
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	// check if name exists
	var f interface{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m := f.(map[string]interface{})
	var val string
	for k, v := range m {
		if k == "name" {
			val = v.(string)
		}
	}
	if len(val) == 0 {
		http.Error(w, fmt.Sprintf("name field is absent from json"), http.StatusForbidden)
		return
	}
	tx, err := c.Dbh.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// reread the body
	b2, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var jid int // primary key from database
	err = tx.Get(&jid, c.Query["insert-jbrowse"], val, string(b2))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u := fmt.Sprintf("%s/%s/%s/%d", r.URL.Scheme, r.URL.Host, r.URL.Path, jid)
	rs := map[string]interface{}{
		"links": map[string]interface{}{
			"self": u,
		},
	}
	// Add location header
	w.Header().Set("Location", u)
	c.Render.JSON(w, http.StatusCreated, rs)
	return
}

func (c *JbConfig) GetNamedHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
}

func (c *JbConfig) CreateNamedHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
}

func (c *JbConfig) UpdateNamedHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
}

func (c *JbConfig) DeleteNamedHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
}
