package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/MITLibraries/fml"
	"github.com/davecgh/go-spew/spew"
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
	file          io.Reader
	rules         []*Rule
	languageCodes map[string]string
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

	languageCodes, err := RetrieveLanguageCodelist()
	if err != nil {
		spew.Dump(err)
	}

	out := make(chan Record)
	p := marcparser{file: m.marcfile, rules: rules, languageCodes: languageCodes}
	go p.parse(out)
	return out
}

func (m *marcparser) parse(out chan Record) {
	mr := fml.NewMarcIterator(m.file)

	for mr.Next() {
		record := mr.Value()

		r, err := marcToRecord(record, m.rules, m.languageCodes)
		if err != nil {
			log.Println(err)
		} else {
			out <- r
		}
	}
	close(out)
}

func marcToRecord(fmlRecord fml.Record, rules []*Rule, languageCodes map[string]string) (r Record, err error) {
	err = nil
	r = Record{}

	r.Identifier = fmlRecord.ControlNum()

	r.Source = "MIT Aleph"
	r.SourceLink = "https://library.mit.edu/item/" + r.Identifier
	r.OclcNumber = applyRule(fmlRecord, rules, "oclc_number")

	lccn := applyRule(fmlRecord, rules, "lccn")
	if lccn != nil {
		r.Lccn = strings.TrimSpace(lccn[0])
	}

	title := applyRule(fmlRecord, rules, "title")
	if title != nil {
		r.Title = title[0]
	} else {
		err = fmt.Errorf("Record %s has no title, check validity", r.Identifier)
		return r, err
	}

	r.AlternateTitles = applyRule(fmlRecord, rules, "alternate_titles")
	r.Creator = applyRule(fmlRecord, rules, "creators")
	r.Contributor = getContributors(fmlRecord, rules, "contributors")

	r.RelatedPlace = applyRule(fmlRecord, rules, "related_place")
	r.RelatedItems = getRelatedItems(fmlRecord, rules, "related_items")

	r.InBibliography = applyRule(fmlRecord, rules, "in_bibliography")

	r.Subject = applyRule(fmlRecord, rules, "subjects")

	r.Isbn = applyRule(fmlRecord, rules, "isbns")
	r.Issn = applyRule(fmlRecord, rules, "issns")
	r.Doi = applyRule(fmlRecord, rules, "dois")

	country := applyRule(fmlRecord, rules, "country_of_publication")
	if country != nil {
		r.Country = country[0]
	}

	// TODO: use lookup tables to translate returned codes to values
	r.Language = applyRule(fmlRecord, rules, "languages")
	r.Language = TranslateLanguageCodes(r.Language, languageCodes)

	r.CallNumber = applyRule(fmlRecord, rules, "call_numbers")

	edition := applyRule(fmlRecord, rules, "edition")
	if edition != nil {
		r.Edition = edition[0]
	}

	r.Imprint = applyRule(fmlRecord, rules, "imprint")

	description := applyRule(fmlRecord, rules, "physical_description")
	if description != nil {
		r.PhysicalDescription = description[0]
	}

	r.PublicationFrequency = applyRule(fmlRecord, rules, "publication_frequency")

	// publication year
	date := applyRule(fmlRecord, rules, "publication_date")
	if date != nil {
		r.PublicationDate = date[0]
	}

	numbering := applyRule(fmlRecord, rules, "numbering")
	if numbering != nil {
		r.Numbering = numbering[0]
	}

	r.Notes = applyRule(fmlRecord, rules, "notes")

	r.Contents = applyRule(fmlRecord, rules, "contents")

	r.Summary = applyRule(fmlRecord, rules, "summary")

	// TODO: use lookup tables to translate returned codes to values
	r.Format = applyRule(fmlRecord, rules, "format")

	// TODO: use lookup tables to translate returned codes to values
	// r.ContentType = contentType(fmlRecord.Leader.Type)

	lf := applyRule(fmlRecord, rules, "literary_form")
	r.LiteraryForm = literaryForm(lf)

	r.Links = getLinks(fmlRecord)

	return r, err
}

func applyRule(fmlRecord fml.Record, rules []*Rule, field string) []string {
	recordFieldRule := getRules(rules, field)

	res := extractData(recordFieldRule, fmlRecord)
	return res
}

// takes a supplied marc rule and fmlRecord returns an array of stringified subfields
func extractData(rule *Rule, fmlRecord fml.Record) []string {
	var field []string
	for _, r := range rule.Fields {
		f := filter(fmlRecord, r)
		for _, y := range f {
			field = append(field, y)
		}
	}
	return field
}

func filter(fmlRecord fml.Record, field *Field) []string {
	var stuff []string
	values := fmlRecord.Filter(field.Tag + field.Subfields)
	for _, f := range values {
		v := strings.Join(f, " ")
		if field.Bytes != "" {
			f := strings.Split(field.Bytes, ":")
			first, _ := strconv.Atoi(f[0])
			take, _ := strconv.Atoi(f[1])
			v = v[first:(first + take)]
		}
		stuff = append(stuff, v)
	}
	return stuff
}

// returns slice of contributors of marc fields taking into account the rules for which fields and subfields we care about as defined in marc_rules.json
func getContributors(fmlRecord fml.Record, rules []*Rule, field string) []*Contributor {
	recordFieldRule := getRules(rules, field)
	var c []*Contributor

	for _, r := range recordFieldRule.Fields {
		y := new(Contributor)
		y.Kind = r.Kind

		y.Value = filter(fmlRecord, r)

		if y.Value != nil {
			c = append(c, y)
		}
	}

	return c
}

// returns slice of related items of marc fields taking into account the rules for which fields and subfields we care about as defined in marc_rules.json
func getRelatedItems(fmlRecord fml.Record, rules []*Rule, field string) []*RelatedItem {
	recordFieldRule := getRules(rules, field)
	var c []*RelatedItem
	for _, r := range recordFieldRule.Fields {
		y := new(RelatedItem)
		y.Kind = r.Kind
		y.Value = filter(fmlRecord, r)
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

// RetrieveLanguageCodelist retrieves language codes for parsing MARC languages
func RetrieveLanguageCodelist() (map[string]string, error) {
	file, err := os.Open("config/languages.xml")
	if err != nil {
		log.Fatal(err)
	}
	// Language struct
	type Language struct {
		Name string `xml:"name"`
		Code string `xml:"code"`
	}

	decoder := xml.NewDecoder(file)
	languages := make(map[string]string)

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "language" {
				var l Language
				decoder.DecodeElement(&l, &se)
				languages[l.Code] = l.Name
			}
		}
	}
	return languages, err
}

// TranslateLanguageCodes takes an array of MARC language codes and returns the language names.
func TranslateLanguageCodes(recordCodes []string, languageCodes map[string]string) []string {
	var languages []string
	for _, l := range recordCodes {
		name := languageCodes[l]
		if name != "" {
			languages = append(languages, name)
		} else {
			languages = append(languages, l)
		}
	}
	return languages
}

// getLinks take a MARC record and eturns an array of Link objects from the 856 field data.
func getLinks(fmlRecord fml.Record) []Link {
	var links []Link
	marc856 := fmlRecord.DataField("856")
	if len(marc856) == 0 {
		return nil
	}
	for _, f := range marc856 {
		ind1 := string(f.Indicator1)
		ind2 := string(f.Indicator2)

		if ind1 == "4" && (ind2 == "0" || ind2 == "1") {
			link := Link{
				Kind:         subfieldValue(f.SubFields, "3"),
				URL:          subfieldValue(f.SubFields, "u"),
				Text:         subfieldValue(f.SubFields, "y"),
				Restrictions: subfieldValue(f.SubFields, "z")}
			links = append(links, link)
		}
	}
	return links
}

func subfieldValue(subs []fml.SubField, code string) string {
	for _, x := range subs {
		if x.Code == code {
			return x.Value
		}
	}
	return ""
}
