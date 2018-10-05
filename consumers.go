package main

import (
	"encoding/json"
	"fmt"
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
	Output string
}

func (js *JSONConsumer) Consume(recs <-chan Record, done chan<- bool) {
	fmt.Println("[")
	var i int
	for r := range recs {
		b, err := json.MarshalIndent(r, "", "    ")
		if err != nil {
			log.Println(err)
		}
		if i != 0 {
			fmt.Println(",")
		}
		fmt.Println(string(b))
		i++
	}
	fmt.Println("]")
	done <- true
}

type TitleConsumer struct {
}

func (ti *TitleConsumer) Consume(recs <-chan Record, done chan<- bool) {
	for r := range recs {
		fmt.Println(r.Title)
	}

	done <- true
}
