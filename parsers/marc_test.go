package marc

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestContains(t *testing.T) {
	a := []byte{'a', 'v', 'x', 'y', 'z'}

	shouldContain := Contains(a, 'a')
	if !shouldContain {
		t.Error("Expected true, got ", shouldContain)
	}

	shouldNotContain := Contains(a, 'b')
	if shouldNotContain {
		t.Error("Expected true, got ", shouldNotContain)
	}
}

func TestConcatSubfields(t *testing.T) {
	file, err := os.Open("../fixtures/record1.mrc")
	if err != nil {
		t.Error(err)
	}
	record, err := marc21.ReadRecord(file)
	if err != nil {
		t.Error(err)
	}

	var subfields []string

	subfields = concatSubfields("245", []byte{'a'}, record)
	if subfields[0] != "Arithmetic /" {
		t.Error("Expected match got", subfields[0])
	}

	subfields = concatSubfields("245", []byte{'a', 'c'}, record)
	if subfields[0] != "Arithmetic / Carl Sandburg ; illustrated as an anamorphic adventure by Ted Rand." {
		t.Error("Expected match got", subfields[0])
	}

	subfields = concatSubfields("650", []byte{'a', 'x'}, record)
	if len(subfields) != 5 {
		t.Error("Expected 5 got", len(subfields))
	}

	if subfields[0] != "Arithmetic Juvenile poetry." {
		t.Error("Expected match got", subfields[0])
	}

	if subfields[4] != "Visual perception." {
		t.Error("Expected match got", subfields[0])
	}
}

func TestStringifySelectSubfields(t *testing.T) {
	file, err := os.Open("../fixtures/record1.mrc")
	if err != nil {
		t.Error(err)
	}
	record, err := marc21.ReadRecord(file)
	if err != nil {
		t.Error(err)
	}

	x := record.GetFields("245")

	subs := []byte{'a'}
	stringified := stringifySelectSubfields(x[0], subs)
	if stringified != "Arithmetic /" {
		t.Error("Expected match, got", stringified)
	}

	subs = []byte{'a', 'c'}
	stringified = stringifySelectSubfields(x[0], subs)
	if stringified != "Arithmetic / Carl Sandburg ; illustrated as an anamorphic adventure by Ted Rand." {
		t.Error("Expected match, got", stringified)
	}

	subs = []byte{'c'}
	stringified = stringifySelectSubfields(x[0], subs)
	if stringified != "Carl Sandburg ; illustrated as an anamorphic adventure by Ted Rand." {
		t.Error("Expected match, got", stringified)
	}
}

func TestMarcToRecord(t *testing.T) {
	file, err := os.Open("../fixtures/record1.mrc")
	if err != nil {
		t.Error(err)
	}
	record, err := marc21.ReadRecord(file)
	if err != nil {
		t.Error(err)
	}

	rules, err := RetrieveRules("../fixtures/marc_rules.json")
	if err != nil {
		spew.Dump(err)
		return
	}

	item := marcToRecord(record, rules)

	if item.author[0] != "Sandburg, Carl, 1878-1967." {
		t.Error("Expected match, got", item.author)
	}

	// yeah, this should be fixed
	if item.identifier != "   92005291 " {
		t.Error("Expected match, got", item.identifier)
	}

	if item.title != "Arithmetic /" {
		t.Error("Expected match, got", item.title)
	}

	if item.contributor[0] != "Rand, Ted, ill." {
		t.Error("Expected match, got", item.contributor[0])
	}

	if item.subject[0] != "Arithmetic Juvenile poetry." {
		t.Error("Expected match, got", item.subject[0])
	}

	if item.url != nil {
		t.Error("Expected no matches, got", item.url)
	}
}
