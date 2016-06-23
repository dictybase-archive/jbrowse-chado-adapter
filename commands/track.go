// Package commands provides method for creating tracks configuration
// that will be served through a HTTP backend
package commands

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/dictybase/jbrowse-chado-adapter/model"
	"github.com/jmoiron/sqlx"
	"github.com/nleof/goyesql"
	"gopkg.in/codegangsta/cli.v1"
)

var coreTypes = []string{
	"chromosome",
	"supercontig",
	"contig",
	"gene",
	"tRNA",
	"ncRNA",
}

var subTypes = map[string]string{
	"mRNA":      "exon",
	"EST_match": "match_part",
}

//An example output the method will produce
//{
//"general": "PolysphondyliumpallidumK5",
//"tracks": [
//{
//"label": "reference",
//"key": "Reference sequence",
//"type": "SequenceTrack",
//"storeClass": "JBrowse/Store/SeqFeature/REST",
//"baseUrl": "http://rest-base-url",
//"query": {
//"retreival": "sequence"
//}
//},
//{
//"label": "supercontig",
//"key": "Supercontigs",
//"type": "CanvasFeatures",
//"storeClass": "JBrowse/Store/SeqFeature/REST",
//"baseUrl": "http://rest-base-url",
//"query": {
//"retreival": "feature",
//"type": "supercontig"
//}
//},
//{
//"label": "gene",
//"key": "Genes",
//"type": "CanvasFeatures",
//"storeClass": "JBrowse/Store/SeqFeature/REST",
//"baseUrl": "http://rest-base-url",
//"query": {
//"retreival": "feature",
//"type": "gene"
//}
//},
//{
//"label": "mRNA",
//"key": "mRNAs",
//"type": "CanvasFeatures",
//"subParts": "exon",
//"storeClass": "JBrowse/Store/SeqFeature/REST",
//"baseUrl": "http://rest-base-url",
//"query": {
//"retreival": "feature",
//"type": "mRNA"
//}
//}
//]
//}
func CreateTracksConf(c *cli.Context) error {
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

	dk := []model.DatasetKey{}
	err = dbh.Select(&dk, sf["get-jbrowse-dataset-ids"], c.String("name"))
	if err != nil {
		return err
	}
	if len(dk) == 0 {
		return fmt.Errorf("no jbrowse configuration with %s name", c.String("name"))
	}
	// get a map of so type names and their database ids
	type2Id, err := mapTypeToDbId(dbh, sf)
	if err != nil {
		return err
	}
	tx, err := dbh.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// It will iterate through all genomes/organism supported by this
	// jbrowse project
	for _, r := range dk {
		ds := model.Dataset{}
		err := tx.Get(&ds, sf["get-each-jbrowse-dataset"], c.String("name"), r.Id)
		if err != nil {
			return err
		}
		// fetch the organism id
		oid, err := strconv.Atoi(strings.Split(ds.Url, "/")[1])
		if err != nil {
			return err
		}
		var goid string
		// create record in the linking table
		err = tx.Get(&goid, sf["insert-jbrowse-organism"], oid, ds.JbrowseId, r.Id)
		if err != nil {
			return err
		}
		refStr, err := getRefseqTrack(c)
		if err != nil {
			return err
		}
		_, err = tx.Exec(sf["insert-jbrowse-track"], refStr, goid)
		if err != nil {
			return err
		}
		for _, typ := range coreTypes {
			exists, err := hasFeature(oid, typ, tx, sf)
			if err != nil {
				return err
			}
			if !exists {
				continue
			}
			fStr, err := getFeatureTrack(typ, c)
			if err != nil {
				return err
			}
			_, err = tx.Exec(sf["insert-jbrowse-track-with-type"], fStr, type2Id[typ], goid)
			if err != nil {
				return err
			}

		}
		for typ, subftyp := range subTypes {
			exists, err := hasFeatWithSubFeat(oid, typ, subftyp, tx, sf)
			if err != nil {
				return err
			}
			if !exists {
				continue
			}
			fStr, err := getFeatWithSubFeatTrack(typ, subftyp, c)
			if err != nil {
				return err
			}
			_, err = tx.Exec(sf["insert-jbrowse-track-with-type"], fStr, type2Id[typ], goid)
			if err != nil {
				return err
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func getRefseqTrack(c *cli.Context) (string, error) {
	ct := map[string]interface{}{
		"label":      "reference",
		"key":        "Reference sequence",
		"type":       "SequenceTrack",
		"storeClass": "JBrowse/Store/SeqFeature/REST",
		"baseUrl":    c.String("base-api"),
		"query": map[string]interface{}{
			"retreival": "sequence",
		},
	}
	js, err := json.Marshal(ct)
	if err != nil {
		return "", err
	}
	return string(js), nil
}

func getFeatureTrack(ftype string, c *cli.Context) (string, error) {
	ct := map[string]interface{}{
		"label":      ftype,
		"key":        fmt.Sprintf("%ss", strings.ToTitle(ftype)),
		"type":       "CanvasFeatures",
		"storeClass": "JBrowse/Store/SeqFeature/REST",
		"baseUrl":    c.String("base-api"),
		"query": map[string]interface{}{
			"retreival": "feature",
			"type":      ftype,
		},
	}
	js, err := json.Marshal(ct)
	if err != nil {
		return "", err
	}
	return string(js), nil
}

func getFeatWithSubFeatTrack(ftype string, subftype string, c *cli.Context) (string, error) {
	ct := map[string]interface{}{
		"label":      ftype,
		"key":        fmt.Sprintf("%ss", strings.ToTitle(ftype)),
		"type":       "CanvasFeatures",
		"subParts":   subftype,
		"storeClass": "JBrowse/Store/SeqFeature/REST",
		"baseUrl":    c.String("base-api"),
		"query": map[string]interface{}{
			"retreival": "feature",
			"type":      ftype,
		},
	}
	js, err := json.Marshal(ct)
	if err != nil {
		return "", err
	}
	return string(js), nil
}

func hasFeature(oid int, ftype string, dbh *sqlx.Tx, sf goyesql.Queries) (bool, error) {
	var c int
	err := dbh.Get(&c, sf["feature-exists"], oid, ftype)
	if err != nil {
		return false, err
	}
	if c > 0 {
		return true, nil
	}
	return false, nil
}

func hasFeatWithSubFeat(oid int, ftype string, subftype string, dbh *sqlx.Tx, sf goyesql.Queries) (bool, error) {
	var c int
	err := dbh.Get(&c, sf["feature-with-subfeat-exists"], oid, ftype, subftype)
	if err != nil {
		return false, err
	}
	if c > 0 {
		return true, nil
	}
	return false, nil
}

func mapTypeToDbId(dbh *sqlx.DB, sf goyesql.Queries) (map[string]int, error) {
	type2Id := make(map[string]int)
	for _, typ := range coreTypes {
		var id int
		err := dbh.Get(&id, sf["get-id-from-so-type"], typ)
		if err != nil {
			return type2Id, err
		}
		type2Id[typ] = id
	}
	for _, typ := range subTypes {
		var id int
		err := dbh.Get(&id, sf["get-id-from-so-type"], typ)
		if err != nil {
			return type2Id, err
		}
		type2Id[typ] = id
	}
	return type2Id, nil
}
