package main

import (
	"encoding/json"
	"io"
	"log"
)

type JSONParser struct {
	file io.Reader
}

type JSONProcessor struct {
	file     io.Reader
	consumer Consumer
	out      chan Record
	done     chan bool
}

func (j *JSONParser) Parse(out chan Record) {
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

func (j *JSONProcessor) Process() {
	p := JSONParser{file: j.file}
	go p.Parse(j.out)
	go j.consumer.Consume(j.out, j.done)

	// wait until the Consume routine reports `done` channel
	<-j.done

	log.Println("Ingested ", ingested, "records")
}
