package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/olivere/elastic"
)

// Consumer defines an interface to be implemented by various consumers
type Consumer interface {
	Consume(<-chan Record, chan<- bool)
}

type ESConsumer struct {
	Index string
	RType string
	p     *elastic.BulkProcessor
}

func (es *ESConsumer) Consume(recs <-chan Record, done chan<- bool) {
	for r := range recs {
		d := elastic.NewBulkIndexRequest().Index(es.Index).Type(es.RType).Doc(r)
		es.p.Add(d)
	}
	done <- true
}

type JSONConsumer struct {
	out io.Writer
}

func (js *JSONConsumer) Consume(recs <-chan Record, done chan<- bool) {
	fmt.Fprintln(js.out, "[")
	var i int
	for r := range recs {
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
	done <- true
}

type TitleConsumer struct {
	out io.Writer
}

func (ti *TitleConsumer) Consume(recs <-chan Record, done chan<- bool) {
	for r := range recs {
		fmt.Fprintln(ti.out, r.Title)
	}

	done <- true
}
