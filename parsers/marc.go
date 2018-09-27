package marc

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/miku/marc21"
)

type record struct {
	identifier   string
	title        string
	author       []string
	contributor  []string
	url          []string
	subject      []string
	isbn         []string
	year         string
	content_type string
}

// Rules defines where the rules are in JSON
type Rules struct {
	Field     string `json:"field"`
	Tag       string `json:"tag"`
	Subfields string `json:"subfields"`
}

// RetrieveRules for parsing MARC
func RetrieveRules(rulefile string) ([]*Rules, error) {
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
	var rules []*Rules
	err = json.NewDecoder(file).Decode(&rules)

	// We don't need to check for errors, the caller can do this.
	return rules, err
}

// Process kicks off the MARC processing
func Process(marcfile io.Reader, rulesfile string) {

	var records []record

	rules, err := RetrieveRules(rulesfile)
	if err != nil {
		spew.Dump(err)
		return
	}

	// loop over all records
	count := 0
	for {
		record, err := marc21.ReadRecord(marcfile)

		// if we get an error, log it
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Println("An error occured processing the", count, "record.")
			log.Fatal(err)
		}

		count++

		// we probably don't want to make this in memory representation of the
		// combined data but instead will probably want to open a JSON file for
		// writing at the start of the loop, write to it on each iteration, and
		// close it when we are done. Or something. Channels?
		// For now I'm just throwing everything into a slice and dumping it because
		// :shrug:
		records = append(records, marcToRecord(record, rules))
	}
	// spew.Dump(records)
	log.Println("Processed ", count, "records")
}

// returns slice of string representations of a given marc field taking into account the rules for which subfields we care about as defined in marc_rules.json
func getFields(marcRecord *marc21.Record, rules []*Rules, field string) []string {
	fieldRules := getRules(rules, field)
	var things []string
	for _, x := range fieldRules {
		things = toRecord(things, x, marcRecord)
	}
	return things
}

func marcToRecord(marcRecord *marc21.Record, rules []*Rules) record {
	r := record{}

	r.identifier = marcRecord.Identifier()

	title := getFields(marcRecord, rules, "title")
	if title != nil {
		r.title = title[0]
	}
	r.author = getFields(marcRecord, rules, "author")
	r.contributor = getFields(marcRecord, rules, "contributor")

	// urls 856:4[0|1] $u
	// only take 856 fields where first indicator is 4
	// only take 856 fields where second indicator is 0 or 1
	// possibly filter out any matches where $3 or $z is "table of contents" or "Publisher description"
	// todo: this does not follow the noted rules yet and instead just grabs anything in 856$u
	r.url = getFields(marcRecord, rules, "url")

	r.subject = getFields(marcRecord, rules, "subject")

	//isbn
	r.isbn = getFields(marcRecord, rules, "isbn")

	// publication year
	// Go to 008 field, 7th byte, grab 4 characters\
	year := getFields(marcRecord, rules, "year")
	if year != nil {
		r.year = year[0][7:11]
	}

	// content type LDR/06:1
	r.content_type = contentType(marcRecord.Leader.Type)
	return r
}

// returns all rules that match a supplied field
func getRules(rules []*Rules, field string) []*Rules {
	var r []*Rules
	for _, v := range rules {
		if v.Field == field {
			r = append(r, v)
		}
	}
	return r
}

func toRecord(field []string, rule *Rules, marcRecord *marc21.Record) []string {
	field = append(field, collectSubfields(rule.Tag, []byte(rule.Subfields), marcRecord)...)
	return field
}

// takes a mark field tag and subfields of interest for a supplied marc record and returns a slice of stringified representations of them
func collectSubfields(marcfield string, subfields []byte, marcrecord *marc21.Record) []string {
	fields := marcrecord.GetFields(marcfield)
	var r []string
	for _, f := range fields {
		r = append(r, stringifySelectSubfields(f, subfields))
	}
	return r
}

func stringifySelectSubfields(field marc21.Field, subfields []byte) string {
	var stringified []string
	switch f := field.(type) {
	case *marc21.DataField:
		for _, s := range f.SubFields {
			if Contains(subfields, s.Code) {
				stringified = append(stringified, s.Value)
			}
		}
	case *marc21.ControlField:
		stringified = append(stringified, f.Data)
	}
	return strings.Join(stringified, " ")
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
