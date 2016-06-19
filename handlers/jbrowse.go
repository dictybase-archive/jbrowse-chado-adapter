// Package handlers providers all the http handlers that
// responds to web requests
package handlers

import (
	"net/http"

	"gopkg.in/unrolled/render.v1"

	"github.com/jmoiron/sqlx"
	"github.com/nleof/goyesql"
	"golang.org/x/net/context"
)

type Jbrowse struct {
	Dbh    *sqlx.DB
	Query  goyesql.Queries
	Render *render.Render
}

func (jb *Jbrowse) GlobalStatsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

}
