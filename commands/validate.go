// Package commands provides methods for validating
// command line arguments
package commands

import (
	"fmt"

	"gopkg.in/codegangsta/cli.v1"
)

func ValidateDbArg(c *cli.Context) error {
	for _, p := range []string{"user", "password", "host", "database"} {
		if !c.IsSet(p) {
			return fmt.Errorf("argument %s is REQUIRED ", p)
		}
	}
	return nil
}
