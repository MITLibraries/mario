package main

import (
	"testing"
)

type Fooer struct{}

func (f *Fooer) Transform(in <-chan Record) <-chan Record {
	out := make(chan Record)
	go func() {
		for r := range in {
			r.Title = r.Title + "FOO"
			out <- r
		}
		close(out)
	}()
	return out
}

type RecordGenerator struct{}

func (g *RecordGenerator) Generate() <-chan Record {
	out := make(chan Record)
	go func() {
		out <- Record{Title: "Bar"}
		out <- Record{Title: "Gaz"}
		close(out)
	}()
	return out
}

type RecordConsumer struct {
	records []Record
}

func (c *RecordConsumer) Consume(in <-chan Record) <-chan bool {
	out := make(chan bool)
	go func() {
		for r := range in {
			c.records = append(c.records, r)
		}
		close(out)
	}()
	return out
}

func TestRun(t *testing.T) {
	c := &RecordConsumer{}
	p := Pipeline{
		generator: &RecordGenerator{},
		consumer:  c,
	}
	p.Next(&Fooer{})
	out := p.Run()
	<-out
	if c.records[0].Title != "BarFOO" {
		t.Error("Expected match, got", c.records[0].Title)
	}
}
