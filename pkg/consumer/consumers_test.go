package consumer

import (
	"bytes"
	"encoding/json"
	"github.com/mitlibraries/mario/pkg/record"
	"strings"
	"testing"
)

func TestTitleConsumerConsume(t *testing.T) {
	var b bytes.Buffer
	in := make(chan record.Record)
	c := TitleConsumer{Out: &b}
	out := c.Consume(in)
	in <- record.Record{Title: "Hatsopoulos Microfluids"}
	close(in)
	<-out
	s := strings.TrimSpace(b.String())
	if s != "Hatsopoulos Microfluids" {
		t.Error("Expected match, got", s)
	}
}

func TestTitleJsonConsume(t *testing.T) {
	var b bytes.Buffer
	in := make(chan record.Record)
	c := JSONConsumer{Out: &b}
	out := c.Consume(in)
	in <- record.Record{Title: "Hatsopoulos Microfluids"}
	close(in)
	<-out

	var records []*record.Record
	json.NewDecoder(&b).Decode(&records)

	if records[0].Title != "Hatsopoulos Microfluids" {
		t.Error("Expected match, got", records[0].Title)
	}
}
