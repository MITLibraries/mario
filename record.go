package main

// Record struct stores our internal mappings of data and is used to when
// mapping various external data sources before sending to elasticsearch
type Record struct {
	Identifier           string
	Title                string
	AlternateTitles      []string
	Creator              []string
	Contributor          []*Contributor
	URL                  []string
	Subject              []string
	Isbn                 []string
	Issn                 []string
	Doi                  []string
	OclcNumber           []string
	Lccn                 string
	Country              string
	Language             []string
	PublicationDate      string
	ContentType          string
	CallNumber           []string
	Edition              string
	Imprint              []string
	PhysicalDescription  string
	PublicationFrequency []string
	Numbering            string
	Notes                []string
	Contents             []string
	Summary              []string
	Format               []string
	LiteraryForm         string
	RelatedPlace         []string
	InBibliography       []string
	RelatedItems         []*RelatedItem
	Links                []Link
	Holdings             []Holdings
}

// Contributor is a port of a Record
type Contributor struct {
	Kind  string
	Value []string
}

// RelatedItem is a port of a Record
type RelatedItem struct {
	Kind  string
	Value []string
}

// Link is a port of a Record
type Link struct {
	Kind         string
	Text         string
	URL          string
	Restrictions string
}

// Holdings is a port of a Record
type Holdings struct {
	Location   string
	CallNumber string
	Status     string
}

// Rule defines where the rules are in JSON
type Rule struct {
	Label  string   `json:"label"`
	Array  bool     `json:"array"`
	Fields []*Field `json:"fields"`
}

// Field defines where the Fields within a Rule are in JSON
type Field struct {
	Tag       string `json:"tag"`
	Subfields string `json:"subfields"`
	Bytes     string `json:"bytes"`
	Kind      string `json:"kind"`
}

// Parser defines an interface common to parsers
type Parser interface {
	Parse(chan Record)
}

// Processor is an interface that allows converting from custom data into
// our Record structure
type Processor interface {
	Process()
}

var ingested int
