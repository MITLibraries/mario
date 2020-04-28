package client

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mitlibraries/mario/pkg/record"
	"github.com/olivere/elastic"
	aws "github.com/olivere/elastic/aws/v4"
	"io/ioutil"
	"net/http"
)

// Primary alias
const primary = "timdex-prod"

// ESClient wraps an olivere/elastic client. Create a new client with the
// NewESClient function.
type ESClient struct {
	client *elastic.Client
	bulker *elastic.BulkProcessor
}

// Current returns the name of the current index for the given prefix. A
// current index is defined as one which is linked to the primary alias. An
// error is returned if there is more than one matching index. An empty
// string indicates there were no matching indexes.
func (c ESClient) Current(prefix string) (string, error) {
	res, err := c.client.Aliases().Index(prefix + "*").Do(context.Background())
	if err != nil {
		return "", err
	}
	aliases := res.IndicesByAlias(primary)
	if len(aliases) == 0 {
		return "", nil
	} else if len(aliases) > 0 {
		return "", errors.New("Could not determine current index")
	} else {
		return aliases[0], nil
	}
}

// Create the new index.
func (c ESClient) Create(index string) error {
	mappings, err := ioutil.ReadFile("config/es_record_mappings.json")
	if err != nil {
		return err
	}
	_, err = c.client.
		CreateIndex(index).
		BodyString(string(mappings)).
		Do(context.Background())
	return err
}

// Start the bulk processor.
func (c *ESClient) Start() error {
	bulker, err := c.client.
		BulkProcessor().
		Name("BulkProcessor").
		Workers(2).
		Do(context.Background())
	c.bulker = bulker
	return err
}

// Stop the bulk processor.
func (c *ESClient) Stop() error {
	return c.bulker.Stop()
}

// Add a record using a bulk processor.
func (c *ESClient) Add(record record.Record, index string) {
	d := elastic.NewBulkIndexRequest().
		Index(index).
		Id(record.Identifier).
		Doc(record)
	c.bulker.Add(d)
}

// Promote will add the given index to the primary alias. If there is an
// existing index matching the prefix linked to the primary alias it will
// be removed from the alias. This action is atomic.
func (c ESClient) Promote(index string, prefix string) error {
	svc := c.client.Alias().Add(index, primary)
	current, err := c.Current(prefix)
	if err != nil {
		return err
	}
	if current != "" {
		svc.Remove(current, primary)
	}
	_, err = svc.Do(context.Background())
	return err
}

// Delete an index.
func (c ESClient) Delete(index string) error {
	_, err := c.client.DeleteIndex(index).Do(context.Background())
	return err
}

func (c ESClient) Indexes() (elastic.CatIndicesResponse, error) {
	return c.client.
		CatIndices().
		Columns("idx", "dc", "h", "s", "id", "ss").
		Do(context.Background())
}

func (c ESClient) Aliases() (elastic.CatAliasesResponse, error) {
	return c.client.CatAliases().Do(context.Background())
}

func (c ESClient) Ping(url string) (*elastic.PingResult, error) {
	res, _, err := c.client.Ping(url).Do(context.Background())
	return res, err
}

// Reindex the source index to the destination index. Returns the number
// of documents reindexed.
func (c ESClient) Reindex(source string, dest string) (int64, error) {
	resp, err := c.client.
		Reindex().
		SourceIndex(source).
		DestinationIndex(dest).
		Do(context.Background())
	if err != nil {
		return 0, err
	}
	return resp.Total, nil
}

// NewESClient creates a new Elasticsearch client.
func NewESClient(url string, v4 bool) (ESClient, error) {
	var client *http.Client
	if v4 {
		sess := session.Must(session.NewSession())
		creds := credentials.NewChainCredentials([]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
		})
		client = aws.NewV4SigningClient(creds, "us-east-1")
	} else {
		client = http.DefaultClient
	}
	es, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetHttpClient(client),
	)
	return ESClient{client: es}, err
}
