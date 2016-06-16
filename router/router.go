// A http router on top of httprouter package to make it compatible with context sensitive
// net/http handler. Currently only Get method is supported.
package router

import (
	"net/http"

	"github.com/cyclopsci/apollo"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

// A Router embeds httprouter
type Router struct {
	*httprouter.Router
}

func NewRouter() *Router {
	return &Router{httprouter.New()}
}

// Wraps around the GET method of httprouter
func (r *Router) Get(path string, handler apollo.Handler) {
	r.GET(path, wrapHandler(handler))
}

// Wrapper function to make httprouter compatible with context aware http handler
// The router parameter values are passed in the context
func wrapHandler(h apollo.Handler) httprouter.Handle {
	f := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h.ServeHTTP(context.WithValue(context.Background(), "params", ps), w, r)
	}
	return f
}
