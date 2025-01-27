package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type File struct {
	ID        int       `db:"id" json:"id"`
	FileName  string    `db:"file_name" json:"file_name"`
	Category  string    `db:"category" json:"category"`
	UserID    uint      `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
type FileMessage struct {
	Type      string `json:"type"`
	FileName  string `json:"fileName"`
	FileData  string `json:"fileData"`
	Recipient string `json:"recipient"`
	Sender    string `json:"sender"`
}

func (f *File) AddToIndex() error {

	if ESClient == nil {
		log.Fatalf("Elasticsearch client is not initialized")
		return fmt.Errorf("Elasticsearch client is not initialized")
	}

	document := struct {
		FileName string `json:"file_name"`
		Category string `json:"category"`
	}{f.FileName, f.Category}

	log.Printf("Document to be indexed: %+v\n", document)

	data, err := json.Marshal(document)
	if err != nil {
		log.Printf("Error marshaling the document: %s", err)
		return err
	}

	log.Println("Marshaled")

	req := esapi.IndexRequest{
		Index:      SearchIndex,
		DocumentID: strconv.Itoa(int(f.ID)),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	log.Println("Created index request")

	res, err := req.Do(context.Background(), ESClient)
	if err != nil {
		log.Printf("Error sending request to Elasticsearch: %s", err)
		return err
	}

	if res.IsError() {
		log.Printf("Error indexing document: %s", res.String())
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	defer res.Body.Close()

	log.Printf("Elasticsearch response: %+v\n", res)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %s", err)
	} else {
		log.Printf("Response body: %s\n", string(body))
	}

	log.Printf("Indexed document %s to index %s\n", res.String(), SearchIndex)
	return nil
}

func FileSearch(searchQuery string) []uint {
	var buf bytes.Buffer

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  searchQuery,
				"fields": []string{"file_name", "category"},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("error encoding: %s", err)
	}

	res, err := ESClient.Search(
		ESClient.Search.WithIndex(SearchIndex),
		ESClient.Search.WithBody(&buf),
	)

	defer res.Body.Close()

	if err != nil || res.IsError() {
		return nil
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil
	}

	var ids []uint

	if hits, ok := r["hits"].(map[string]interface{}); ok {
		if hitsHits, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsHits {
				if hitMap, ok := hit.(map[string]interface{}); ok {
					if idStr, ok := hitMap["_id"].(string); ok {
						id, _ := strconv.Atoi(idStr)
						ids = append(ids, uint(id))
					}
				}
			}
		}
	}

	return ids
}
