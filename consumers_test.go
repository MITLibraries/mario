package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestTitleConsumerConsume(t *testing.T) {
	var b bytes.Buffer
	in := make(chan Record)
	c := TitleConsumer{out: &b}
	out := c.Consume(in)
	in <- Record{Title: "Hatsopoulos Microfluids"}
	close(in)
	<-out
	s := strings.TrimSpace(b.String())
	if s != "Hatsopoulos Microfluids" {
		t.Error("Expected match, got", s)
	}
}

func TestTitleJsonConsume(t *testing.T) {
	var b bytes.Buffer
	in := make(chan Record)
	c := JSONConsumer{out: &b}
	out := c.Consume(in)
	in <- Record{Title: "Hatsopoulos Microfluids"}
	close(in)
	<-out

	var records []*Record
	json.NewDecoder(&b).Decode(&records)

	if records[0].Title != "Hatsopoulos Microfluids" {
		t.Error("Expected match, got", records[0].Title)
	}
}
