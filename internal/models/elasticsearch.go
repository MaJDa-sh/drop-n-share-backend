package models

import (
	"context"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var ESClient *elasticsearch.Client

const SearchIndex = "files"

func init() {
	var err error
	ESClient, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			os.Getenv("ELASTICSEARCH_URL"),
		},
	})
	if err != nil {
		log.Fatalf("Error initializing Elasticsearch client: %s", err)
	}
}

func ESCcreateIndexIfNotExists() {
	_, err := esapi.IndicesExistsAliasRequest{
		Index: []string{SearchIndex},
	}.Do(context.Background(), ESClient)

	if err != nil {
		ESClient.Indices.Create(SearchIndex)
	}
}
