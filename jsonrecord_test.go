package main

import (
	"os"
	"testing"
)

func TestJsonParser(t *testing.T) {
	jsonfile, err := os.Open("fixtures/mit_test_records.json")
	if err != nil {
		t.Error(err)
	}

	out := make(chan Record)

	p := jsonparser{file: jsonfile}
	go p.parse(out)

	var chanLength int
	for range out {
		chanLength++
	}

	if chanLength != 1962 {
		t.Error("Expected match, got", chanLength)
	}
}

func TestJsonProcess(t *testing.T) {
	jsonfile, err := os.Open("fixtures/mit_test_records.json")
	if err != nil {
		t.Error(err)
	}

	var i int
	p := JSONGenerator{file: jsonfile}
	for range p.Generate() {
		i++
	}

	if i != 1962 {
		t.Error("Expected match, got", i)
	}
}
