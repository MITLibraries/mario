package marc

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/MITLibraries/marc21"
	"github.com/davecgh/go-spew/spew"
)

type record struct {
	identifier  string
	title       string
	author      []string
	contributor []string
	url         []string
	subject     []string
	isbn        []string
}

const marcRules = "../fixtures/marc_rules.json"

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
func Process(rulesfile string) {

	var records []record

	rules, err := RetrieveRules(rulesfile)
	if err != nil {
		spew.Dump(err)
		return
	}

	// loop over all records
	count := 0
	for {
		record, err := marc21.ReadRecord(os.Stdin)

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
	spew.Dump(records)
	log.Println("Processed ", count, "records")
}

func marcToRecord(marcRecord *marc21.Record, rules []*Rules) record {
	r := record{}

	r.identifier = marcRecord.Identifier()

	// main entry
	rule := rules[0]
	r.title = concatSubfields(rule.Tag, []byte(rule.Subfields), marcRecord)[0]

	// author
	r.author = toRecord(r.author, rules[1], marcRecord)

	// contributors
	r.contributor = toRecord(r.contributor, rules[2], marcRecord)

	// urls
	r.url = toRecord(r.url, rules[3], marcRecord)

	// subjects
	r.subject = toRecord(r.subject, rules[4], marcRecord)
	r.subject = toRecord(r.subject, rules[5], marcRecord)
	r.subject = toRecord(r.subject, rules[6], marcRecord)
	r.subject = toRecord(r.subject, rules[7], marcRecord)

	//isbn
	r.isbn = toRecord(r.isbn, rules[8], marcRecord)
	return r
}

func toRecord(field []string, rule *Rules, marcRecord *marc21.Record) []string {
	field = append(field, concatSubfields(rule.Tag, []byte(rule.Subfields), marcRecord)...)
	return field
}

// takes a mark field tag and subfields of interest for a supplied marc record and returns them concatenated
func concatSubfields(marcfield string, subfields []byte, marcrecord *marc21.Record) []string {
	x := marcrecord.GetFields(marcfield)
	var r []string
	for _, y := range x {
		r = append(r, stringifySelectSubfields(y.GetSubfields(), subfields))
	}
	return r
}

// Returns specified subfields concatenated in order they appear in the field
func stringifySelectSubfields(subs []*marc21.SubField, keep []byte) string {
	var stringified []string
	for _, f := range subs {
		if !Contains(keep, f.Code) {
			continue
		}
		stringified = append(stringified, f.Value)
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
