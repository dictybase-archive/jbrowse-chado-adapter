// Package commands provides method for handing jbrowe_conf.json configuration.
package commands

import (
	"encoding/json"
	"fmt"

	"github.com/dictybase/jbrowse-chado-adapter/model"
	"gopkg.in/codegangsta/cli.v1"
)

// Generates and save a new jbrowse_conf.json configuration
// in the postgresql database.
// It primarilly generates the dataset json configuration by
// reading the list of available genomes in the database.
//	Example
//      {
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
func CreateConf(c *cli.Context) error {
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
	//query
	org := []model.Organism{}
	err = dbh.Select(&org, sf["get-organism-with-features"])
	if err != nil {
		return err
	}
	ct := map[string]interface{}{
		"dataRoot": c.String("data-root"),
	}
	for _, odata := range org {
		data := fmt.Sprintf("datasets.%s%s", odata.Genus, odata.Species)
		url := fmt.Sprintf("?data=%s/%d", c.String("genome-root"), odata.OrganismId)
		name := fmt.Sprintf("%s %s", odata.Genus, odata.Species)
		ct[data] = map[string]interface{}{
			"url":  url,
			"name": name,
		}
	}
	config := map[string]interface{}{
		"general": ct,
	}
	str, err := json.Marshal(config)
	if err != nil {
		return err
	}
	tx, err := dbh.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec("INSERT INTO jbrowse(name, configuration) VALUES (?, ?)", c.String("name"), string(str))
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
