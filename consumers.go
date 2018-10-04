package main

import (
	"github.com/olivere/elastic"
)

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
