package main

import (
	"os"

	"github.com/dictybase/jbrowse-chado-adapter/commands"
	"gopkg.in/codegangsta/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "jbc"
	app.Usage = "A jbrowse backend server for chado database"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log,l",
			Usage: "Name of the log file(optional), default goes to stderr",
		},
		cli.StringFlag{
			Name:   "user, u",
			Usage:  "chado database user[REQUIRED]",
			EnvVar: "CHADO_USER",
		},
		cli.StringFlag{
			Name:   "password, p",
			Usage:  "chado database password[REQUIRED]",
			EnvVar: "CHADO_PASS",
		},
		cli.StringFlag{
			Name:   "database, db",
			Usage:  "chado database name[REQUIRED]",
			EnvVar: "CHADO_DB",
		},
		cli.StringFlag{
			Name:   "host, h",
			Usage:  "chado database host[REQUIRED]",
			EnvVar: "CHADO_HOST",
		},
		cli.IntFlag{
			Name:  "port, p",
			Usage: "server port",
			Value: 9998,
		},
		cli.StringFlag{
			Name:  "sql-file",
			Usage: "Text file with sql queries, default is bundled with the package",
		},
	}
	app.Before = commands.ValidateDbArg
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "runs the jbrowse backend server",
			Action: commands.RunServer,
		},
		{
			Name:   "bootstrap-conf",
			Usage:  "Generates and saves a new jbrowse_conf.json configuration in the postgresql database",
			Action: commands.CreateConf,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "data-root",
					Usage: "jbrowse data directory",
					Value: "data",
				},
				cli.StringFlag{
					Name:  "genome-root",
					Usage: "The url path prepended for genome dataset for looking up of individual trackList.json",
					Value: "genome",
				},
			},
		},
	}
	app.Run(os.Args)
}
