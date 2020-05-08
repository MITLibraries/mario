package pipeline

import "github.com/mitlibraries/mario/pkg/record"

//A Pipeline builds and runs a data pipeline for process Records. A
//Pipeline consists of exactly one Generator, one Consumer and zero or
//more Transformers.
type Pipeline struct {
	Generator    Generator
	Transformers []Transformer
	Consumer     Consumer
}

//The Transformer interface can be used to create an intermediate stage
//in a Pipeline.
type Transformer interface {
	Transform(<-chan record.Record) <-chan record.Record
}

//The Generator interface should be used to create the initial stage of
//a Pipeline.
type Generator interface {
	Generate() <-chan record.Record
}

//The Consumer interface should be used to create the last stage of a
//Pipeline.
type Consumer interface {
	Consume(<-chan record.Record) <-chan bool
}

//Next adds one or more Transformers to the Pipeline. Next can be called
//multiple times. All Transformers will be run in the order added.
func (p *Pipeline) Next(t ...Transformer) {
	p.Transformers = append(p.Transformers, t...)
}

//Run the Pipeline. Be sure to read from the empty channel that's returned
//as that signals the Pipeline has finished running.
func (p *Pipeline) Run() <-chan bool {
	out := p.Generator.Generate()
	for _, t := range p.Transformers {
		out = t.Transform(out)
	}
	return p.Consumer.Consume(out)
}
