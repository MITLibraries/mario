package main

import (
	"encoding/json"
	"io"
	"log"
)

type jsonparser struct {
	file io.Reader
}

//JSONGenerator parses JSON-formatted MARC records.
type JSONGenerator struct {
	file io.Reader
}

func (j *jsonparser) parse(out chan Record) {
	ingested = 0
	decoder := json.NewDecoder(j.file)

	// read open bracket
	_, err := decoder.Token()
	if err != nil {
		log.Fatal(err)
	}

	for decoder.More() {
		var r Record
		err = decoder.Decode(&r)
		if err != nil {
			log.Fatal(err)
		}
		ingested++
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
func (j *JSONGenerator) Generate() <-chan Record {
	out := make(chan Record)
	p := jsonparser{file: j.file}
	go p.parse(out)
	return out
}
