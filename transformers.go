package main

//Counter transformer records the number of records handled.
type Counter struct {
	Count int
}

//Transform counts the records.
func (c *Counter) Transform(in <-chan Record) <-chan Record {
	out := make(chan Record)
	go func() {
		for r := range in {
			c.Count++
			out <- r
		}
		close(out)
	}()
	return out
}
