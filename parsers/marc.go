package marc

import (
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
}

// Process kicks off the MARC processing
func Process() {

	var records []record

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
		records = append(records, marcToRecord(record))
	}
	spew.Dump(records)
	log.Println("Processed ", count, "records")
}

func marcToRecord(marcRecord *marc21.Record) record {
	var subfields []byte
	r := record{}

	r.identifier = marcRecord.Identifier()

	// main entry
	subfields = []byte{'a', 'b', 'f', 'g', 'k', 'n', 'p', 's'}
	r.title = concatSubfields("245", subfields, marcRecord)[0]

	// author
	subfields = []byte{'a', 'b', 'c', 'd', 'e', 'q'}
	r.author = append(r.author, concatSubfields("100", subfields, marcRecord)...)

	// contributors
	subfields = []byte{'a', 'b', 'c', 'd', 'e', 'q'}
	r.contributor = append(r.contributor, concatSubfields("700", subfields, marcRecord)...)

	// 856 $u
	urls := marcRecord.GetFields("856")
	for _, url := range urls {
		keep := []byte{'a', 'b', 'c', 'd', 'e', 'q'}
		r.url = append(r.url, stringifySelectSubfields(url.GetSubfields(), keep))
	}

	subfields = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'x', 'y', 'z'}
	r.subject = append(r.subject, concatSubfields("600", subfields, marcRecord)...)
	r.subject = append(r.subject, concatSubfields("610", subfields, marcRecord)...)

	subfields = []byte{'a', 'v', 'x', 'y', 'z'}
	r.subject = append(r.subject, concatSubfields("650", subfields, marcRecord)...)
	r.subject = append(r.subject, concatSubfields("651", subfields, marcRecord)...)

	return r
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
