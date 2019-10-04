package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/olivere/elastic"
)

//ESConsumer adds Records to ElasticSearch.
type ESConsumer struct {
	Index string
	RType string
	p     *elastic.BulkProcessor
}

//Consume the records.
func (es *ESConsumer) Consume(in <-chan Record) <-chan bool {
	out := make(chan bool)
	go func() {
		for r := range in {
			d := elastic.NewBulkIndexRequest().
				Index(es.Index).
				Id(r.Identifier).
				Type(es.RType).
				Doc(r)
			es.p.Add(d)
		}
		close(out)
	}()
	return out
}

//JSONConsumer outputs Records as JSON. The Records will be written
//to JSONConsumer.out.
type JSONConsumer struct {
	out io.Writer
}

//Consume the records.
func (js *JSONConsumer) Consume(in <-chan Record) <-chan bool {
	out := make(chan bool)
	go func() {
		fmt.Fprintln(js.out, "[")
		var i int
		for r := range in {
			b, err := json.MarshalIndent(r, "", "    ")
			if err != nil {
				log.Println(err)
			}
			if i != 0 {
				fmt.Fprintln(js.out, ",")
			}
			fmt.Fprintln(js.out, string(b))
			i++
		}
		fmt.Fprintln(js.out, "]")
		close(out)
	}()
	return out
}

//TitleConsumer just outputs the title of Records. The titles will be
//written to TitleConsumer.out.
type TitleConsumer struct {
	out io.Writer
}

//Consume the records.
func (t *TitleConsumer) Consume(in <-chan Record) <-chan bool {
	out := make(chan bool)
	go func() {
		for r := range in {
			fmt.Fprintln(t.out, r.Title)
		}
		close(out)
	}()
	return out
}

//SilentConsumer is useful for debugging sometimes
type SilentConsumer struct {
	out io.Writer
}

//Consume the records and close the channel when done. No processing is done.
func (s *SilentConsumer) Consume(in <-chan Record) <-chan bool {
	out := make(chan bool)
	go func() {
		for range in {
			continue
		}
		close(out)
	}()
	return out
}
