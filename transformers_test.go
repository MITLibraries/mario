package main

import (
	"testing"
)

func TestCounterTransform(t *testing.T) {
	in := make(chan Record, 2)
	in <- Record{Title: "Foo"}
	in <- Record{Title: "Bar"}
	close(in)
	c := Counter{}
	out := c.Transform(in)
	<-out
	if c.Count != 2 {
		t.Error("Expected match, got", c.Count)
	}
}
