package main

import (
	"os"
	"testing"

	"github.com/MITLibraries/fml"
	"github.com/davecgh/go-spew/spew"
)

func TestMarcToRecord(t *testing.T) {
	file, err := os.Open("fixtures/record1.mrc")
	if err != nil {
		t.Error(err)
	}
	records := fml.NewMarcIterator(file)
	_ = records.Next()
	record, err := records.Value()

	if err != nil {
		t.Error(err)
	}

	rules, err := RetrieveRules("config/marc_rules.json")
	if err != nil {
		spew.Dump(err)
		return
	}

	languageCodes, err := RetrieveLanguageCodelist()
	if err != nil {
		spew.Dump(err)
		return
	}

	item, _ := marcToRecord(record, rules, languageCodes)

	if item.Creator[0] != "Sandburg, Carl, 1878-1967." {
		t.Error("Expected match, got", item.Creator)
	}

	if item.Identifier != "92005291" {
		t.Error("Expected match, got", item.Identifier)
	}

	if item.Title != "Arithmetic /" {
		t.Error("Expected match, got", item.Title)
	}

	if item.Contributor[0].Value[0] != "Rand, Ted, ill." {
		t.Error("Expected match, got", item.Contributor[0].Value[0])
	}

	if item.Subject[0] != "Arithmetic Juvenile poetry." {
		t.Error("Expected match, got", item.Subject[0])
	}

	if item.PublicationDate != "1993" {
		t.Error("Expected match, got", item.PublicationDate)
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
	languageCodes, err := RetrieveLanguageCodelist()
	if err != nil {
		spew.Dump(err)
		return
	}

	in := []string{"abk", "ach", "afa", "aaa", ""}
	out := []string{"Abkhaz", "Acoli", "Afroasiatic (Other)", "aaa", ""}
	langs := TranslateLanguageCodes(in, languageCodes)

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
	rules, err := RetrieveRules("config/marc_rules.json")

	if err != nil {
		spew.Dump(err)
		return
	}

	marcfile, err := os.Open("fixtures/test.mrc")
	if err != nil {
		t.Error(err)
	}

	out := make(chan Record)

	p := marcparser{file: marcfile, rules: rules}
	go p.parse(out)

	var chanLength int
	for range out {
		chanLength++
	}

	if chanLength != 85 {
		t.Error("Expected match, got", chanLength)
	}
}

func TestMarcProcess(t *testing.T) {
	marcfile, err := os.Open("fixtures/test.mrc")
	if err != nil {
		t.Error(err)
	}
	p := MarcGenerator{marcfile: marcfile, rulesfile: "config/marc_rules.json"}
	out := p.Generate()
	var i int
	for range out {
		i++
	}
	if i != 85 {
		t.Error("Expected match, got", i)
	}
}
