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
	}
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "runs the jbrowse backend server",
			Action: commands.RunServer,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "sql-file",
					Usage: "Text file with sql queries, default is bundled with the package",
				},
			},
		},
	}
	app.Run(os.Args)
}
