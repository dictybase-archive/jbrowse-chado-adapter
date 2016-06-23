package handlers

import (
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

func (c *JbConfig) CreateHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

}

func (c *JbConfig) GetNamedHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

}

func (c *JbConfig) CreateNamedHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

}
