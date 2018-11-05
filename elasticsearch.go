package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/olivere/elastic"
	aws "github.com/olivere/elastic/aws/v4"
)

// EsClient configures the elasticsearch client
func esClient(url string, index string, v4 bool) (*elastic.Client, error) {
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
	return elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetHttpClient(client),
	)
}

// CreateRecordIndex creates our Record index
func createRecordIndex(client *elastic.Client, index string) error {
	mappings, err := ioutil.ReadFile("config/es_record_mappings.json")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	_, err = client.CreateIndex(index).BodyString(string(mappings)).Do(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Index created")
	return nil
}
