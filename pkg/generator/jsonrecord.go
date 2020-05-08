package generator

import (
	"encoding/json"
	"github.com/mitlibraries/mario/pkg/record"
	"io"
	"log"
)

type jsonparser struct {
	file io.Reader
}

//JSONGenerator parses JSON records.
type JSONGenerator struct {
	File io.Reader
}

func (j *jsonparser) parse(out chan record.Record) {
	decoder := json.NewDecoder(j.file)

	// read open bracket
	_, err := decoder.Token()
	if err != nil {
		log.Fatal(err)
	}

	for decoder.More() {
		var r record.Record
		err = decoder.Decode(&r)
		if err != nil {
			log.Fatal(err)
		}
		out <- r
	}

	// read closing bracket
	_, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}

	close(out)
}

//Generate creates a channel of Records.
func (j *JSONGenerator) Generate() <-chan record.Record {
	out := make(chan record.Record)
	p := jsonparser{file: j.File}
	go p.parse(out)
	return out
}
