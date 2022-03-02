package ingester

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/mitlibraries/mario/pkg/client"
	"github.com/mitlibraries/mario/pkg/consumer"
	"github.com/mitlibraries/mario/pkg/generator"
	"github.com/mitlibraries/mario/pkg/pipeline"
	"github.com/mitlibraries/mario/pkg/transformer"
)

// Config is a structure for passing a set of configuration parameters to
// an Ingester.
type Config struct {
	Filename string
	Source   string
	Consumer string
	Index    string
	NewIndex bool
	Promote  bool
}

// NewStream returns an io.ReadCloser from a path string. The path can be
// either a local directory path or a URL for an S3 object.
func NewStream(filename string) (io.ReadCloser, error) {
	parts, err := url.Parse(filename)
	if err != nil {
		return nil, err
	}
	if parts.Scheme == "s3" {
		return client.GetS3Obj(parts.Host, parts.Path)
	}
	return os.Open(filename)
}

// Ingester does the work of ingesting a data stream.
type Ingester struct {
	Stream    io.ReadCloser
	config    Config
	generator pipeline.Generator
	consumer  pipeline.Consumer
	Client    client.Indexer
}

// Configure an Ingester. This should be called before Ingest.
func (i *Ingester) Configure(config Config) error {
	var err error
	// Configure generator
	if config.Source == "alma" {
		i.generator = &generator.MarcGenerator{Marcfile: i.Stream}
	} else if config.Source == "aspace" {
		i.generator = &generator.ArchivesGenerator{Archivefile: i.Stream}
	} else if config.Source == "dspace" {
		i.generator = &generator.DspaceGenerator{Dspacefile: i.Stream}
	} else if config.Source == "mario" {
		i.generator = &generator.JSONGenerator{File: i.Stream}
	} else {
		return errors.New("Unknown source data")
	}

	// Configure consumer
	if config.Consumer == "es" {
		if config.NewIndex == true {
			now := time.Now().UTC()
			config.Index = fmt.Sprintf("%s-%s", config.Source, now.Format("2006-01-02t15-04-05z"))
		} else {
			current, err := i.Client.Current(config.Source)
			if err != nil || current == "" {
				e := fmt.Errorf("No existing production index for source '%s'. Either promote an existing %s index or add the 'new' flag to the ingest command to create a new index.", config.Source, config.Source)
				return e
			}
			log.Printf("Ingesting into current production index: %s", current)
			config.Index = current
			config.Promote = false
		}

		err = i.Client.Create(config.Index)
		if err != nil {
			return err
		}
		i.consumer = &consumer.ESConsumer{
			Index:  config.Index,
			RType:  "Record",
			Client: i.Client,
		}

		log.Printf("Configured OpenSearch consumer using source: %s, index: %s, and promote: %s", config.Source, config.Index, strconv.FormatBool(config.Promote))

	} else if config.Consumer == "json" {
		i.consumer = &consumer.JSONConsumer{Out: os.Stdout}
	} else if config.Consumer == "title" {
		i.consumer = &consumer.TitleConsumer{Out: os.Stdout}
	} else if config.Consumer == "silent" {
		i.consumer = &consumer.SilentConsumer{Out: os.Stdout}
	} else {
		return errors.New("Unknown consumer")
	}

	i.config = config
	return nil
}

// Ingest the configured data stream. The Ingester should have been
// configured before calling this method. It will return the number of
// ingested documents.
func (i *Ingester) Ingest() (int, error) {
	var err error
	p := pipeline.Pipeline{
		Generator: i.generator,
		Consumer:  i.consumer,
	}
	ctr := &transformer.Counter{}
	p.Next(ctr)
	if i.config.Consumer == "es" {
		err = i.Client.Start()
		if err != nil {
			return 0, err
		}
		defer i.Client.Stop()
	}
	out := p.Run()
	<-out
	if i.config.Promote {
		log.Printf("Automatic promotion is happening")
		err = i.Client.Promote(i.config.Index)
	}
	return ctr.Count, err
}
