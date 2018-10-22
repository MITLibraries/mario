package main

//A Pipeline builds and runs a data pipeline for process Records. A
//Pipeline consists of exactly one Generator, one Consumer and zero or
//more Transformers.
type Pipeline struct {
	generator    Generator
	transformers []Transformer
	consumer     Consumer
}

//The Transformer interface can be used to create an intermediate stage
//in a Pipeline.
type Transformer interface {
	Transform(<-chan Record) <-chan Record
}

//The Generator interface should be used to create the initial stage of
//a Pipeline.
type Generator interface {
	Generate() <-chan Record
}

//The Consumer interface should be used to create the last stage of a
//Pipeline.
type Consumer interface {
	Consume(<-chan Record) <-chan bool
}

//Next adds one or more Transformers to the Pipeline. Next can be called
//multiple times. All Transformers will be run in the order added.
func (p *Pipeline) Next(t ...Transformer) {
	p.transformers = append(p.transformers, t...)
}

//Run the Pipeline. Be sure to read from the empty channel that's returned
//as that signals the Pipeline has finished running.
func (p *Pipeline) Run() <-chan bool {
	out := p.generator.Generate()
	for _, t := range p.transformers {
		out = t.Transform(out)
	}
	return p.consumer.Consume(out)
}
