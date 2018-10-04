package marc

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

type Record struct {
	Identifier      string
	Title           string
	AlternateTitles []string
	Creator         []string
	Contributor     []*Contributor
	Url             []string
	Subject         []string
	Isbn            []string
	Issn            []string
	Doi             []string
	Country         string
	Language        []string
	Year            string
	ContentType     string
	CallNumber      []string
	RelatedItems    []RelatedItem
	Links           []Link
	Holdings        []Holdings
}

type Contributor struct {
	Kind  string
	Value []string
}

type RelatedItem struct {
	Kind  string
	Value string
}

type Link struct {
	Kind         string
	Text         string
	Url          string
	Restrictions string
}

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

var consumed int
var ingested int

type MarcParser struct {
	file  io.Reader
	rules []*Rule
	out   chan Record
}

func (m *MarcParser) Parse() {
	for {
		record, err := marc21.ReadRecord(m.file)

		// if we get an error, log it
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Println("An error occured processing the", ingested, "record.")
			log.Fatal(err)
		}

		ingested++

		m.out <- marcToRecord(record, m.rules)
	}
	close(m.out)
}

// Process kicks off the MARC processing
func Process(marcfile io.Reader, rulesfile string) {
	rules, err := RetrieveRules(rulesfile)

	if err != nil {
		spew.Dump(err)
		return
	}

	out := make(chan Record)
	done := make(chan bool, 1)

	p := MarcParser{file: marcfile, rules: rules, out: out}
	go p.Parse()
	go ConsumeRecords(out, done)

	// wait until the ConsumeRecords routine reports it is done via `done` channel
	<-done

	log.Println("Ingested ", ingested, "records")
	log.Println("Finished", consumed, "records")
}

// ConsumeRecords currently just prints record titles
func ConsumeRecords(rec <-chan Record, done chan<- bool) {
	for r := range rec {
		consumed++
		log.Println(r.Title)
	}

	// indicate over done channel this routine is complete
	done <- true
}

// trasforms a single marc21 record into our internal record struct
func marcToRecord(marcRecord *marc21.Record, rules []*Rule) Record {
	r := Record{}

	r.Identifier = marcRecord.Identifier()

	title := getFields(marcRecord, rules, "title")
	if title != nil {
		r.Title = title[0]
	}
	r.AlternateTitles = getFields(marcRecord, rules, "alternate_titles")
	r.Creator = getFields(marcRecord, rules, "creators")
	r.Contributor = getContributors(marcRecord, rules, "contributors")

	// urls 856:4[0|1] $u
	// only take 856 fields where first indicator is 4
	// only take 856 fields where second indicator is 0 or 1
	// possibly filter out any matches where $3 or $z is "table of contents" or "Publisher description"
	// todo: this does not follow the noted rules yet and instead just grabs anything in 856$u
	// r.url = getFields(marcRecord, rules, "url")

	r.Subject = getFields(marcRecord, rules, "subjects")

	//isbn
	r.Isbn = getFields(marcRecord, rules, "isbns")
	r.Issn = getFields(marcRecord, rules, "issns")
	r.Doi = getFields(marcRecord, rules, "dois")

	country := getFields(marcRecord, rules, "country_of_publication")
	if country != nil {
		r.Country = country[0]
	}

	r.Language = getFields(marcRecord, rules, "languages")
	r.CallNumber = getFields(marcRecord, rules, "call_numbers")

	// publication year
	year := getFields(marcRecord, rules, "year")
	if year != nil {
		r.Year = year[0]
	}

	// content type LDR/06:1
	r.ContentType = contentType(marcRecord.Leader.Type)
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
