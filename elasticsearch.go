package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/olivere/elastic"
	aws "github.com/olivere/elastic/aws/v4"
)

var (
	prodAlias = "timdex-prod"
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
	log.Println("Index created:", index)
	return nil
}

func previous(client *elastic.Client, prefix string) ([]string, error) {
	// retrieve all indexes linked to timdex-prod alias and filter by supplied prefix. These are the "old" indexes.
	aliases, err := aliases(client)
	if err != nil {
		return nil, err
	}

	var indexes []string

	for _, a := range aliases {
		if a.Alias == prodAlias && strings.HasPrefix(a.Index, prefix) {
			indexes = append(indexes, a.Index)
		}
	}

	if len(aliases) == 0 {
		log.Printf("No aliases found. Nothing to demote.")
	}

	log.Println("Previous indexes:", indexes)

	return indexes, nil
}

func aliases(client *elastic.Client) (elastic.CatAliasesResponse, error) {
	ctx := context.Background()
	aliases, err := client.CatAliases().Do(ctx)
	if err != nil {
		return nil, err
	}

	return aliases, nil
}

func delete(client *elastic.Client, index string) error {
	ctx := context.Background()
	_, err := client.DeleteIndex(index).Do(ctx)
	if err != nil {
		return err
	}
	log.Printf("Index deleted")
	return nil
}

func demote(client *elastic.Client, index string) error {
	ctx := context.Background()
	_, err := client.Alias().Remove(index, prodAlias).Do(ctx)
	if err != nil {
		return err
	}
	log.Println("Index demoted:", index)
	return nil
}

func indexes(client *elastic.Client) error {
	ctx := context.Background()
	indexes, err := client.CatIndices().Do(ctx)
	if err != nil {
		return err
	}
	for _, i := range indexes {
		fmt.Printf("Name: %s \n"+
			"  DocsCount: %d \n"+
			"  Health: %s \n"+
			"  Status: %s \n"+
			"  UUID: %s \n"+
			"  StoreSize: %s \n\n",
			i.Index, i.DocsCount, i.Health, i.Status, i.UUID, i.StoreSize)
	}
	if len(indexes) == 0 {
		log.Printf("No indexes found.")
	}
	return nil
}

func ping(client *elastic.Client, url string) error {
	ctx := context.Background()
	ping, code, err := client.Ping(url).Do(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Response code: %d \n"+
		"Name: %s \n"+
		"Cluster Name: %s \n"+
		"Tag line: %s \n"+
		"Version: %s \n"+
		"BuildHash: %s \n"+
		"LuceneVersion: %s \n",
		code, ping.Name, ping.ClusterName, ping.TagLine, ping.Version.Number,
		ping.Version.BuildHash, ping.Version.LuceneVersion)
	return nil
}

func promote(client *elastic.Client, index string) error {
	ctx := context.Background()
	_, err := client.Alias().Add(index, prodAlias).Do(ctx)
	if err != nil {
		return err
	}
	log.Println("Index promoted:", index)
	return nil
}

func after(exID int64, requests []elastic.BulkableRequest, resp *elastic.BulkResponse, err error) {
	if resp.Errors == true {
		fmt.Printf("Request ID: %d -- Errors: %t\n", exID, resp.Errors)
		errs := resp.Failed()
		for _, e := range errs {
			fmt.Println(e.Error)
		}
	}
}

func reindex(client *elastic.Client, source string, destination string) error {
	ctx := context.Background()
	resp, err := client.Reindex().SourceIndex(source).DestinationIndex(destination).Do(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Reindexed %d docs from %s into %s\n", resp.Total, source, destination)
	return nil
}
