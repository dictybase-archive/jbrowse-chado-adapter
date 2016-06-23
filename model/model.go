// Package model provides various structure for modeling data from database
package model

type Organism struct {
	OrganismId int
	Species    string
	Genus      string
}

type DatasetKey struct {
	Id string
}

type Dataset struct {
	Url       string
	Name      string
	JbrowseId int
}
