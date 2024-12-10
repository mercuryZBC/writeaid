package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"sync"
	"yuqueppbackend/service-base/config"
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
		mapping := map[string]interface{}{
			"settings": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 1,
			},
			"mappings": map[string]interface{}{
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type": "text",
					},
					"created_at": map[string]interface{}{
						"type": "date",
					},
				},
			},
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
