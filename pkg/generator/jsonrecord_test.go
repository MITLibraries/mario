package generator

import (
	"github.com/mitlibraries/mario/pkg/record"
	"os"
	"testing"
)

func TestJsonParser(t *testing.T) {
	jsonfile, err := os.Open("../../fixtures/timdex_record_samples.json")
	if err != nil {
		t.Error(err)
	}

	out := make(chan record.Record)

	p := jsonparser{file: jsonfile}
	go p.parse(out)

	var chanLength int
	for range out {
		chanLength++
	}

	if chanLength != 6 {
		t.Error("Expected match, got", chanLength)
	}
}

func TestJsonProcess(t *testing.T) {
	jsonfile, err := os.Open("../../fixtures/timdex_record_samples.json")
	if err != nil {
		t.Error(err)
	}

	var i int
	p := JSONGenerator{File: jsonfile}
	for range p.Generate() {
		i++
	}

	if i != 6 {
		t.Error("Expected match, got", i)
	}
}
