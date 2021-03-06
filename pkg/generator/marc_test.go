package generator

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/mitlibraries/fml"
	"github.com/mitlibraries/mario/pkg/record"
)

func TestMarcToRecord(t *testing.T) {
	file, err := os.Open("../../fixtures/alma_samples.mrc")
	if err != nil {
		t.Error(err)
	}
	records := fml.NewMarcIterator(file)
	_ = records.Next()
	record, err := records.Value()

	if err != nil {
		t.Error(err)
	}

	rules, err := RetrieveRules("/config/marc_rules.json")
	if err != nil {
		spew.Dump(err)
		return
	}

	languageCodes, err := RetrieveCodelist("language", "/config/languages.xml")
	if err != nil {
		spew.Dump(err)
		return
	}

	countryCodes, err := RetrieveCodelist("country", "/config/countries.xml")
	if err != nil {
		spew.Dump(err)
		return
	}

	item, _ := marcToRecord(record, rules, languageCodes, countryCodes)

	if item.Contributor[0].Value != "D'Rivera, Paquito, 1948-" {
		t.Error("Expected match, got", item.Contributor[0].Value)
	}

	if item.Contributor[0].Kind != "author" {
		t.Error("Expected match, got", item.Contributor[0].Kind)
	}

	if item.Identifier != "990026671500206761" {
		t.Error("Expected match, got", item.Identifier)
	}

	if item.Title != "Spice it up! the best of Paquito D'Rivera." {
		t.Error("Expected match, got", item.Title)
	}

	if item.Subject[0] != "Jazz." {
		t.Error("Expected match, got", item.Subject[0])
	}

	if item.PublicationDate != "2008" {
		t.Error("Expected match, got", item.PublicationDate)
	}
}

func TestMarcHoldings(t *testing.T) {
	file, err := os.Open("../../fixtures/alma_holdings_test_records.mrc")
	if err != nil {
		t.Error(err)
	}
	records := fml.NewMarcIterator(file)
	_ = records.Next()
	record, err := records.Value()

	if err != nil {
		t.Error(err)
	}

	rules, err := RetrieveRules("/config/marc_rules.json")
	if err != nil {
		spew.Dump(err)
		return
	}

	languageCodes, err := RetrieveCodelist("language", "/config/languages.xml")
	if err != nil {
		spew.Dump(err)
		return
	}

	countryCodes, err := RetrieveCodelist("country", "/config/countries.xml")
	if err != nil {
		spew.Dump(err)
		return
	}

	// This record has an 852, but no 866
	item, _ := marcToRecord(record, rules, languageCodes, countryCodes)

	h := item.Holdings[0]
	if h.Location != "Hayden Library" {
		t.Error("Expected match, got", h.Location)
	}
	if h.Collection != "Graphic Novel Collection" {
		t.Error("Expected match, got", h.Collection)
	}

	if h.Format != "Print volume" {
		t.Error("Expected match, got", h.Format)
	}

	// This record has an 866 field and 852. We use 852 and append the 866 summary.
	_ = records.Next()
	record, _ = records.Value()
	item, _ = marcToRecord(record, rules, languageCodes, countryCodes)
	h = item.Holdings[0]
	if h.Location != "Barker Library" {
		t.Error("Expected match, got", h.Location)
	}
	if h.Collection != "Stacks" {
		t.Error("Expected match, got", h.Collection)
	}
	if h.Summary != "MCM 1995 and updates" {
		t.Error("Expected match, got", h.Summary)
	}
	if h.Format != "Print volume" {
		t.Error("Expected match, got", h.Format)
	}
}

var contenttypetests = []struct {
	in  byte
	out string
}{
	{'a', "Text"},
	{'b', "Text"},
	{'c', "Musical score"},
	{'d', "Musical score"},
	{'e', "Cartographic material"},
	{'f', "Cartographic material"},
	{'g', "Moving image"},
	{'h', "Text"},
	{'i', "Sound recording"},
	{'j', "Sound recording"},
	{'k', "Still image"},
	{'l', "Text"},
	{'m', "Computer file"},
	{'n', "Text"},
	{'o', "Kit"},
	{'p', "Mixed materials"},
	{'q', "Text"},
	{'r', "Object"},
	{'s', "Text"},
	{'t', "Text"},
	{'u', "Text"},
	{'v', "Text"},
	{'w', "Text"},
	{'x', "Text"},
	{'y', "Text"},
	{'z', "Text"},
}

func TestContentType(t *testing.T) {
	for _, ct := range contenttypetests {
		t.Run(string(ct.in), func(t *testing.T) {
			ctCase := contentType(ct.in)
			if ctCase != ct.out {
				t.Errorf("got %q, want %q", ctCase, ct.out)
			}
		})
	}
}

func TestTranslateLanguageCodes(t *testing.T) {
	languageCodes, err := RetrieveCodelist("language", "/config/languages.xml")
	if err != nil {
		spew.Dump(err)
		return
	}

	in := []string{"abk", "ach", "afa", "aaa", ""}
	out := []string{"Abkhaz", "Acoli", "Afroasiatic (Other)", "aaa", ""}
	langs := TranslateCodes(in, languageCodes)

	if len(langs) != len(out) {
		t.Errorf("got %q items, want %q", len(langs), len(out))
		return
	}

	for i, x := range langs {
		if x != out[i] {
			t.Errorf("got %q, want %q", x, out[i])
		}
	}
}

func TestMarcParser(t *testing.T) {
	rules, err := RetrieveRules("/config/marc_rules.json")

	if err != nil {
		spew.Dump(err)
		return
	}

	marcfile, err := os.Open("../../fixtures/alma_samples.mrc")
	if err != nil {
		t.Error(err)
	}

	out := make(chan record.Record)

	p := marcparser{file: marcfile, rules: rules}
	go p.parse(out)

	var chanLength int
	for range out {
		chanLength++
	}

	if chanLength != 358 {
		t.Error("Expected match, got", chanLength)
	}
}

func TestMarcProcess(t *testing.T) {
	marcfile, err := os.Open("../../fixtures/alma_samples.mrc")
	if err != nil {
		t.Error(err)
	}
	p := MarcGenerator{Marcfile: marcfile}
	out := p.Generate()
	var i int
	for range out {
		i++
	}
	if i != 358 {
		t.Error("Expected match, got", i)
	}
}

func TestStringInSlice(t *testing.T) {
	l := []string{"hello", "goodbye"}
	r := stringInSlice("hello", l)
	if r != true {
		t.Error("Expected true, got", r)
	}
}

func TestOclcs(t *testing.T) {
	file, err := os.Open("../../fixtures/alma_samples.mrc")
	if err != nil {
		t.Error(err)
	}
	records := fml.NewMarcIterator(file)
	_ = records.Next()
	record, err := records.Value()

	if err != nil {
		t.Error(err)
	}

	rules, err := RetrieveRules("/config/marc_rules.json")
	if err != nil {
		spew.Dump(err)
		return
	}

	languageCodes, err := RetrieveCodelist("language", "/config/languages.xml")
	if err != nil {
		spew.Dump(err)
		return
	}

	countryCodes, err := RetrieveCodelist("country", "/config/countries.xml")
	if err != nil {
		spew.Dump(err)
		return
	}

	item, _ := marcToRecord(record, rules, languageCodes, countryCodes)

	// Confirm oclc prefix is removed
	if item.OclcNumber[0] != "811549562" {
		t.Error("Expected match, got", item.OclcNumber)
	}

	// Confirm old system numbers are not included.
	if len(item.OclcNumber) != 1 {
		t.Error("Expected 1, got", len(item.OclcNumber))
	}
}
