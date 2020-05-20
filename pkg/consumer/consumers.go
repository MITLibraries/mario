package consumer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/mitlibraries/mario/pkg/client"
	"github.com/mitlibraries/mario/pkg/record"
)

//ESConsumer adds Records to ElasticSearch.
type ESConsumer struct {
	Index  string
	RType  string
	Client client.Indexer
}

//Consume the records.
func (es *ESConsumer) Consume(in <-chan record.Record) <-chan bool {
	out := make(chan bool)
	go func() {
		for r := range in {
			es.Client.Add(r, es.Index, es.RType)
		}
		close(out)
	}()
	return out
}

//JSONConsumer outputs Records as JSON. The Records will be written
//to JSONConsumer.out.
type JSONConsumer struct {
	Out io.Writer
}

//Consume the records.
func (js *JSONConsumer) Consume(in <-chan record.Record) <-chan bool {
	out := make(chan bool)
	go func() {
		fmt.Fprintln(js.Out, "[")
		var i int
		for r := range in {
			b, err := json.MarshalIndent(r, "", "    ")
			if err != nil {
				log.Println(err)
			}
			if i != 0 {
				fmt.Fprintln(js.Out, ",")
			}
			fmt.Fprintln(js.Out, string(b))
			i++
		}
		fmt.Fprintln(js.Out, "]")
		close(out)
	}()
	return out
}

//TitleConsumer just outputs the title of Records. The titles will be
//written to TitleConsumer.out.
type TitleConsumer struct {
	Out io.Writer
}

//Consume the records.
func (t *TitleConsumer) Consume(in <-chan record.Record) <-chan bool {
	out := make(chan bool)
	go func() {
		for r := range in {
			fmt.Fprintln(t.Out, r.Title)
		}
		close(out)
	}()
	return out
}

//SilentConsumer is useful for debugging sometimes
type SilentConsumer struct {
	Out io.Writer
}

//Consume the records and close the channel when done. No processing is done.
func (s *SilentConsumer) Consume(in <-chan record.Record) <-chan bool {
	out := make(chan bool)
	go func() {
		for range in {
			continue
		}
		close(out)
	}()
	return out
}
