package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"yuqueppbackend/service-base/config"

	"github.com/elastic/go-elasticsearch/v8"
)

var esClient *elasticsearch.Client
var esOnce sync.Once

func GetElasticSearchClient() *elasticsearch.Client {
	esOnce.Do(func() {
		cfg := elasticsearch.Config{
			Addresses: []string{
				config.GetElasticSearchAddress(),
			},
		}
		var err error
		esClient, err = elasticsearch.NewClient(cfg)
		if err != nil {
			log.Fatalf("Error creating the client: %s", err)
		}

		// 测试连接
		res, err := esClient.Info()
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()
		fmt.Println(res)
		if esClient == nil {
			panic("elasticSearch client not initialized")
		}
		err = checkAndCreateIndex(esClient, "document")
		if err != nil {
			log.Fatalf("Error creating document: %s", err)
		}
		err = checkAndCreateIndex(esClient, "knowledgebase")
		if err != nil {
			log.Fatalf("Error creating knowledgeBase: %s", err)
		}
	})
	return esClient
}

func checkAndCreateIndex(es *elasticsearch.Client, indexName string) error {
	// 1. 检查索引是否存在
	res, err := es.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("error checking index existence: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 { // 索引不存在
		// 2. 创建索引

		var mapping map[string]interface{}
		switch indexName {
		//"id":      docIdStr,
		//"title":   document.Title,
		//"content": content,
		case "document":
			mapping = map[string]interface{}{
				"mappings": map[string]interface{}{
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":  "text",
							"index": false,
						},
						"title": map[string]interface{}{
							"type":     "text",
							"index":    true,
							"analyzer": "standard",
						},
						"content": map[string]interface{}{
							"type":     "text",
							"index":    true,
							"analyzer": "standard",
						},
					},
				},
			}
		case "knowledgebase":
			//"kb_id":          strKbId,
			//"kb_name":        kb.Name,
			//"kb_description": kb.Description,
			mapping = map[string]interface{}{
				"mappings": map[string]interface{}{
					"properties": map[string]interface{}{
						"kb_id": map[string]interface{}{
							"type":  "text",
							"index": false,
						},
						"kb_name": map[string]interface{}{
							"type":     "text",
							"index":    true,
							"analyzer": "standard",
						},
						"kb_description": map[string]interface{}{
							"type":     "text",
							"index":    true,
							"analyzer": "standard",
						},
					},
				},
			}

		}
		mappingBytes, _ := json.Marshal(mapping)

		res, err := es.Indices.Create(
			indexName,
			es.Indices.Create.WithBody(bytes.NewReader(mappingBytes)),
		)
		if err != nil {
			return fmt.Errorf("error creating index: %w", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("failed to create index: %s", res.String())
		}
		fmt.Println("Index created successfully:", indexName)
	} else if res.StatusCode == 200 { // 索引已存在
		fmt.Println("Index already exists:", indexName)
	} else { // 处理其他可能的状态码
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}
