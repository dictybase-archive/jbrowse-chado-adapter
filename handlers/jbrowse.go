// Package handlers providers all the http handlers that
// responds to web requests
package handlers

import (
	"database/sql"
	"net/http"

	"github.com/nleof/goyesql"
	"golang.org/x/net/context"
)

type Jbrowse struct {
	Dbh   *sql.DB
	Query goyesql.Queries
}

func (jb *Jbrowse) GlobalStatsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

}
