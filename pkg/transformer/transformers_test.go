package transformer

import (
	"github.com/mitlibraries/mario/pkg/record"
	"testing"
)

func TestCounterTransform(t *testing.T) {
	in := make(chan record.Record, 2)
	in <- record.Record{Title: "Foo"}
	in <- record.Record{Title: "Bar"}
	close(in)
	c := Counter{}
	out := c.Transform(in)
	<-out
	if c.Count != 2 {
		t.Error("Expected match, got", c.Count)
	}
}
