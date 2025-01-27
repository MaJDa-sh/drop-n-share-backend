package config

// import (
// 	"log"
// 	"os"

// 	"github.com/elastic/go-elasticsearch/v8"
// )

// var ESClient *elasticsearch.Client

// func InitElasticsearch() {
// 	var err error
// 	ESClient, err = elasticsearch.NewClient(elasticsearch.Config{
// 		Addresses: []string{
// 			os.Getenv("ELASTICSEARCH_URL"),
// 		},
// 	})
// 	if err != nil {
// 		log.Fatalf("Error initializing Elasticsearch client: %s", err)
// 	}
// }
