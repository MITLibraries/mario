package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/miku/marc21"
)

// RetrieveRules for parsing MARC
func RetrieveRules(rulefile string) ([]*Rule, error) {
	// Open the file.
	file, err := os.Open(rulefile)
	if err != nil {
		return nil, err
	}

	// Schedule the file to be closed once
	// the function returns.
	defer file.Close()

	// Decode the file into a slice of pointers
	// to Feed values.
	var rules []*Rule
	err = json.NewDecoder(file).Decode(&rules)
	// We don't need to check for errors, the caller can do this.
	return rules, err
}

type marcparser struct {
	file  io.Reader
	rules []*Rule
}

//MarcGenerator parses binary MARC records.
type MarcGenerator struct {
	marcfile  io.Reader
	rulesfile string
}

//Generate a channel of Records.
func (m *MarcGenerator) Generate() <-chan Record {
	rules, err := RetrieveRules(m.rulesfile)

	if err != nil {
		spew.Dump(err)
	}

	out := make(chan Record)
	p := marcparser{file: m.marcfile, rules: rules}
	go p.parse(out)
	return out
}

func (m *marcparser) parse(out chan Record) {
	for {
		record, err := marc21.ReadRecord(m.file)

		// if we get an error, log it
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		out <- marcToRecord(record, m.rules)
	}
	close(out)
}

// trasforms a single marc21 record into our internal record struct
func marcToRecord(marcRecord *marc21.Record, rules []*Rule) Record {
	r := Record{}

	r.Identifier = marcRecord.Identifier()
	r.OclcNumber = getFields(marcRecord, rules, "oclc_number")

	lccn := getFields(marcRecord, rules, "lccn")
	if lccn != nil {
		r.Lccn = lccn[0]
	}

	title := getFields(marcRecord, rules, "title")
	if title != nil {
		r.Title = title[0]
	}
	r.AlternateTitles = getFields(marcRecord, rules, "alternate_titles")
	r.Creator = getFields(marcRecord, rules, "creators")
	r.Contributor = getContributors(marcRecord, rules, "contributors")

	r.RelatedPlace = getFields(marcRecord, rules, "related_place")

	r.RelatedItems = getRelatedItems(marcRecord, rules, "related_items")

	r.InBibliography = getFields(marcRecord, rules, "in_bibliography")

	// urls 856:4[0|1] $u
	// only take 856 fields where first indicator is 4
	// only take 856 fields where second indicator is 0 or 1
	// possibly filter out any matches where $3 or $z is "table of contents" or "Publisher description"
	// todo: this does not follow the noted rules yet and instead just grabs anything in 856$u
	// r.url = getFields(marcRecord, rules, "url")

	// TODO: Links may be best represented by extracting a few values from 856 and _not_ contatanating them but instead filtering on some values and storing them in the Link structs

	r.Subject = getFields(marcRecord, rules, "subjects")

	r.Isbn = getFields(marcRecord, rules, "isbns")
	r.Issn = getFields(marcRecord, rules, "issns")
	r.Doi = getFields(marcRecord, rules, "dois")

	country := getFields(marcRecord, rules, "country_of_publication")
	if country != nil {
		r.Country = country[0]
	}

	// TODO: use lookup tables to translate returned codes to values
	r.Language = getFields(marcRecord, rules, "languages")

	r.CallNumber = getFields(marcRecord, rules, "call_numbers")

	edition := getFields(marcRecord, rules, "edition")
	if edition != nil {
		r.Edition = edition[0]
	}

	r.Imprint = getFields(marcRecord, rules, "imprint")

	description := getFields(marcRecord, rules, "physical_description")
	if description != nil {
		r.PhysicalDescription = description[0]
	}

	r.PublicationFrequency = getFields(marcRecord, rules, "publication_frequency")

	// publication year
	date := getFields(marcRecord, rules, "publication_date")
	if date != nil {
		r.PublicationDate = date[0]
	}

	numbering := getFields(marcRecord, rules, "numbering")
	if numbering != nil {
		r.Numbering = numbering[0]
	}

	r.Notes = getFields(marcRecord, rules, "notes")

	r.Contents = getFields(marcRecord, rules, "contents")

	r.Summary = getFields(marcRecord, rules, "summary")

	// TODO: use lookup tables to translate returned codes to values
	r.Format = getFields(marcRecord, rules, "format")

	// TODO: use lookup tables to translate returned codes to values
	r.ContentType = contentType(marcRecord.Leader.Type)

	lf := getFields(marcRecord, rules, "literary_form")
	r.LiteraryForm = literaryForm(lf)
	return r
}

// returns slice of string representations of marc fields taking into account the rules for which fields and subfields we care about as defined in marc_rules.json
func getFields(marcRecord *marc21.Record, rules []*Rule, field string) []string {
	recordFieldRule := getRules(rules, field)
	return applyRule(recordFieldRule, marcRecord)
}

// returns slice of contributors of marc fields taking into account the rules for which fields and subfields we care about as defined in marc_rules.json
func getContributors(marcRecord *marc21.Record, rules []*Rule, field string) []*Contributor {
	recordFieldRule := getRules(rules, field)
	var c []*Contributor
	for _, r := range recordFieldRule.Fields {
		y := new(Contributor)
		y.Kind = r.Kind
		y.Value = collectSubfields(r, marcRecord)
		if y.Value != nil {
			c = append(c, y)
		}
	}
	return c
}

// returns slice of related items of marc fields taking into account the rules for which fields and subfields we care about as defined in marc_rules.json
func getRelatedItems(marcRecord *marc21.Record, rules []*Rule, field string) []*RelatedItem {
	recordFieldRule := getRules(rules, field)
	var c []*RelatedItem
	for _, r := range recordFieldRule.Fields {
		y := new(RelatedItem)
		y.Kind = r.Kind
		y.Value = collectSubfields(r, marcRecord)
		if y.Value != nil {
			c = append(c, y)
		}
	}
	return c
}

// returns all rules that match a supplied fieldname
func getRules(rules []*Rule, label string) *Rule {
	for _, v := range rules {
		if v.Label == label {
			return v
		}
	}
	return nil // TODO: this will lead to a panic and end the world. While this is ultimately an appropriate response to failing to find rules we expect to find, it would be better to handle that explictly and log something that explains it before terminating cleanly.
}

// takes a supplied marc rule and marcrecord returns an array of stringified subfields
func applyRule(rule *Rule, marcRecord *marc21.Record) []string {
	var field []string
	for _, r := range rule.Fields {
		field = append(field, collectSubfields(r, marcRecord)...)
	}
	return field
}

// takes our local Field structure that contains our processing rules and a MARC21.Record and returns a slice of stringified representations of the fields we are interested in
func collectSubfields(field *Field, marcrecord *marc21.Record) []string {
	fields := marcrecord.GetFields(field.Tag)
	var r []string
	for _, f := range fields {
		subs := stringifySelectSubfields(f, []byte(field.Subfields))
		if field.Bytes != "" && subs != "" {
			f := strings.Split(field.Bytes, ":")
			first, _ := strconv.Atoi(f[0])
			take, _ := strconv.Atoi(f[1])
			r = append(r, subs[first:(first+take)])
		} else {
			r = append(r, subs)
		}
	}
	return r
}

// keeps only supplied subfields (effectively filtering out unwanted subfields) while maintaining order of subfields in supplied marc21.Field and returns them by joining them into a string
func stringifySelectSubfields(field marc21.Field, subfields []byte) string {
	var keep []string
	switch f := field.(type) {
	case *marc21.DataField:
		for _, s := range f.SubFields {
			if Contains(subfields, s.Code) {
				keep = append(keep, s.Value)
			}
		}
	case *marc21.ControlField:
		keep = append(keep, f.Data)
	}
	return strings.Join(keep, " ")
}

// Contains tells whether a contains x.
func Contains(a []byte, x byte) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func literaryForm(x []string) string {
	var t string
	if x == nil {
		return ""
	}
	switch x[0] {
	case "0", "s", "e":
		t = "nonfiction"
	default:
		t = "fiction"
	}
	return t
}

// Content type mappings
func contentType(x byte) string {
	var t string
	switch x {
	case 'c':
		t = "Musical score"
	case 'd':
		t = "Musical score"
	case 'e':
		t = "Cartographic material"
	case 'f':
		t = "Cartographic material"
	case 'g':
		t = "Moving image"
	case 'i':
		t = "Sound recording"
	case 'j':
		t = "Sound recording"
	case 'k':
		t = "Still image"
	case 'm':
		t = "Computer file"
	case 'o':
		t = "Kit"
	case 'p':
		t = "Mixed materials"
	case 'r':
		t = "Object"
	default:
		t = "Text"
	}
	return t
}
