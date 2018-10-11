package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestTitleConsumerConsume(t *testing.T) {
	rules, err := RetrieveRules("fixtures/marc_rules.json")

	if err != nil {
		t.Error(err)
		return
	}

	marcfile, err := os.Open("fixtures/record1.mrc")
	if err != nil {
		t.Error(err)
	}
	done := make(chan bool, 1)
	out := make(chan Record)

	p := MarcParser{file: marcfile, rules: rules, out: out}
	go p.Parse()

	var b bytes.Buffer
	consumer := &TitleConsumer{out: &b}
	go consumer.Consume(out, done)

	// wait until the ConsumeRecords routine reports it is done via `done` channel
	<-done

	if strings.TrimSpace(b.String()) != "Arithmetic /" {
		t.Error("Expected match, got", b.String())
	}
}

func TestTitleJsonConsume(t *testing.T) {
	rules, err := RetrieveRules("fixtures/marc_rules.json")

	if err != nil {
		t.Error(err)
		return
	}

	marcfile, err := os.Open("fixtures/mit_test_records.mrc")
	if err != nil {
		t.Error(err)
	}
	done := make(chan bool, 1)
	out := make(chan Record)

	p := MarcParser{file: marcfile, rules: rules, out: out}
	go p.Parse()

	var b bytes.Buffer
	consumer := &JSONConsumer{out: &b}
	go consumer.Consume(out, done)

	// wait until the ConsumeRecords routine reports it is done via `done` channel
	<-done

	var records []*Record
	json.NewDecoder(&b).Decode(&records)

	if records[0].Title != "Black Panther adventures." {
		t.Error("Expected match, got", records[0].Title)
	}
	if records[0].Identifier != "002621216" {
		t.Error("Expected match, got", records[0].Identifier)
	}
}
