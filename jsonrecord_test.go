package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestJsonParser(t *testing.T) {
	jsonfile, err := os.Open("fixtures/mit_test_records.json")
	if err != nil {
		t.Error(err)
	}

	out := make(chan Record)

	p := JSONParser{file: jsonfile}
	go p.Parse(out)

	var chanLength int
	for _ = range out {
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
	var buf bytes.Buffer
	log.SetOutput(&buf)
	tmp := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	out := make(chan Record)
	done := make(chan bool, 1)

	consumer := &TitleConsumer{out: os.Stdout}
	p := JSONProcessor{file: jsonfile, consumer: consumer, out: out, done: done}
	p.Process()

	log.SetOutput(os.Stderr)
	os.Stdout = tmp
	if !strings.Contains(buf.String(), "Ingested  1962 records") {
		t.Error("Expected match, got", buf.String())
	}
}
