package middlewares

import (
	"net/http"

	"github.com/cyclopsci/apollo"
	"github.com/jmoiron/sqlx"
	"github.com/nleof/goyesql"
	"golang.org/x/net/context"
)

type JbConfig struct {
	Dbh   *sqlx.DB
	Query goyesql.Queries
}

func (c *JbConfig) NameMiddleware(h apollo.Handler) apollo.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	}
	return apollo.HandlerFunc(fn)
}
