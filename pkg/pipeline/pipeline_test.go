package pipeline

import (
	"github.com/mitlibraries/mario/pkg/record"
	"testing"
)

type Fooer struct{}

func (f *Fooer) Transform(in <-chan record.Record) <-chan record.Record {
	out := make(chan record.Record)
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

func (g *RecordGenerator) Generate() <-chan record.Record {
	out := make(chan record.Record)
	go func() {
		out <- record.Record{Title: "Bar"}
		out <- record.Record{Title: "Gaz"}
		close(out)
	}()
	return out
}

type RecordConsumer struct {
	records []record.Record
}

func (c *RecordConsumer) Consume(in <-chan record.Record) <-chan bool {
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
		Generator: &RecordGenerator{},
		Consumer:  c,
	}
	p.Next(&Fooer{})
	out := p.Run()
	<-out
	if c.records[0].Title != "BarFOO" {
		t.Error("Expected match, got", c.records[0].Title)
	}
}
