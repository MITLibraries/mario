package main

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/MITLibraries/marc21"
)

type item struct {
	identifier  string
	title       string
	author      []string
	contributor []string
	url         []string
	subject     []string
}

// example: cat MARC_FILE_WITH_ERRORS | go run marc/marc.go
func main() {

	var items []item

	// loop over all records
	count := 0
	for {
		record, err := marc21.ReadRecord(os.Stdin)
		count++

		// if we get to the end of the file, stop doing stuff
		if err == io.EOF {
			break
		}

		if count == 100 {
			break
		}

		// if we get an error, log it
		if err != nil {
			log.Println("An error occured processing the", count, "record.")
			log.Fatal(err)
		}

		// we probably don't want to make this in memory representation of the
		// combined data but instead will probably want to open a JSON file for
		// writing at the start of the loop, write to it on each iteration, and
		// close it when we are done. Or something.
		// For now I'm just throwing everything into a slice and dumping it because
		// :shrug:
		items = append(items, marcToRecord(record))
	}
}

func marcToRecord(record *marc21.Record) item {
	var subfields []byte
	i := item{}

	i.identifier = record.Identifier()

	// main entry
	subfields = []byte{'a', 'b', 'f', 'g', 'k', 'n', 'p', 's'}
	i.title = concatSubfields("245", subfields, record)[0]

	// author
	subfields = []byte{'a', 'b', 'c', 'd', 'e', 'q'}
	i.author = append(i.author, concatSubfields("100", subfields, record)...)

	// contributors
	subfields = []byte{'a', 'b', 'c', 'd', 'e', 'q'}
	i.contributor = append(i.contributor, concatSubfields("700", subfields, record)...)

	// 856 $u
	urls := record.GetFields("856")
	for _, url := range urls {
		keep := []byte{'a', 'b', 'c', 'd', 'e', 'q'}
		i.url = append(i.url, stringifySelectSubfields(url.GetSubfields(), keep))
	}

	subfields = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'x', 'y', 'z'}
	i.subject = append(i.subject, concatSubfields("600", subfields, record)...)
	i.subject = append(i.subject, concatSubfields("610", subfields, record)...)

	subfields = []byte{'a', 'v', 'x', 'y', 'z'}
	i.subject = append(i.subject, concatSubfields("650", subfields, record)...)
	i.subject = append(i.subject, concatSubfields("651", subfields, record)...)

	return i
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
