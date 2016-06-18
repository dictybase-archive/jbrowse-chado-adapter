package commands

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cyclopsci/apollo"
	"github.com/dictybase/jbrowse-chado-adapter/handlers"
	"github.com/dictybase/jbrowse-chado-adapter/middlewares"
	"github.com/dictybase/jbrowse-chado-adapter/query"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/nleof/goyesql"
	"github.com/rs/cors"
	"golang.org/x/net/context"
	"gopkg.in/codegangsta/cli.v1"
)

// Runs the http server
func RunServer(c *cli.Context) error {
	for _, p := range []string{"user", "password", "host", "database"} {
		if !c.IsSet(p) {
			return fmt.Errorf("argument %s is REQUIRED ", p)
		}
	}

	var logMw *middlewares.Logger
	if c.IsSet("log") {
		w, err := os.Create(c.String("log"))
		defer w.Close()
		if err != nil {
			return fmt.Errorf("cannot open log file %q\n", err)
		}
		logMw = middlewares.NewFileLogger(w)
	} else {
		logMw = middlewares.NewLogger()
	}

	cors := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
	})

	// sql file
	sf, err := getSqlResource(c)
	if err != nil {
		return err
	}
	// db connection
	dbh, err := getDbHandler(c)
	if err != nil {
		return err
	}
	jb := &handlers.Jbrowse{dbh, sf}
	r := mux.NewRouter()
	sgChain := apollo.New(
		apollo.Wrap(cors.Handler),
		apollo.Wrap(logMw.LoggerMiddleware)).
		With(context.Background()).
		ThenFunc(jb.GlobalStatsHandler)
	r.Handle("/stats/global", sgChain).Methods("GET")
	http.Handle("/", r)
	log.Printf("Starting web server on port %d\n", c.Int("port"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Int("port")), nil))
	return nil
}

// Parse and return a reader for sql file
func getSqlResource(c *cli.Context) (goyesql.Queries, error) {
	var sf goyesql.Queries
	if c.IsSet("sql-file") {
		sf, err := goyesql.ParseFile(c.String("sql-file"))
		return sf, err
	}
	b, err := query.Asset("resources/jbrowse.sql")
	if err != nil {
		return sf, err
	}
	sf, err = goyesql.ParseBytes(b)
	return sf, err
}

func getDbHandler(c *cli.Context) (*sql.DB, error) {
	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		c.String("user"), c.String("password"),
		c.String("database"), c.String("host"))
	return sql.Open("postgres", connString)
}
