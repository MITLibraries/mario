package transformer

import "github.com/mitlibraries/mario/pkg/record"

//Counter transformer records the number of records handled.
type Counter struct {
	Count int
}

//Transform counts the records.
func (c *Counter) Transform(in <-chan record.Record) <-chan record.Record {
	out := make(chan record.Record)
	go func() {
		for r := range in {
			c.Count++
			out <- r
		}
		close(out)
	}()
	return out
}
