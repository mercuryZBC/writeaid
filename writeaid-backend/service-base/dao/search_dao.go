package dao

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"yuqueppbackend/service-base/util"
)

type SearchDao struct {
	esClient *elasticsearch.Client
}

func NewSearchDao(esClient *elasticsearch.Client) *SearchDao {
	return &SearchDao{esClient: esClient}
}

func (dao *SearchDao) GetDocIDByTitleFuzzy(title string) ([]string, error) {
	// 构造查询条件，使用 match 查询并设置 fuzziness 实现模糊匹配
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": map[string]interface{}{
					"query":     title,  // 用户输入的标题
					"fuzziness": "AUTO", // 启用模糊匹配
				},
			},
		},
	}

	// 将查询条件转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// 执行搜索
	res, err := util.GetElasticSearchClient().Search(
		dao.esClient.Search.WithIndex("document"),
		dao.esClient.Search.WithBody(&buf),
		dao.esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 解析结果
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 提取文档 ID
	var ids []string
	if hits, ok := result["hits"].(map[string]interface{}); ok {
		for _, hit := range hits["hits"].([]interface{}) {
			if doc, ok := hit.(map[string]interface{}); ok {
				ids = append(ids, doc["_id"].(string))
			}
		}
	}

	return ids, nil
}

func (dao *SearchDao) GetDocIDByContent(content string) ([]string, error) {
	// 构造查询条件
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"content": content,
			},
		},
	}

	// 将查询条件转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// 执行搜索
	res, err := dao.esClient.Search(
		dao.esClient.Search.WithIndex("document"),
		dao.esClient.Search.WithBody(&buf),
		dao.esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 解析结果
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 提取文档 ID
	var ids []string
	if hits, ok := result["hits"].(map[string]interface{}); ok {
		for _, hit := range hits["hits"].([]interface{}) {
			if doc, ok := hit.(map[string]interface{}); ok {
				ids = append(ids, doc["_id"].(string))
			}
		}
	}

	return ids, nil
}

func (dao *SearchDao) GetKnowledgeBaseIDByName(name string) ([]string, error) {
	// 构造查询条件
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"kb_name": map[string]interface{}{
					"query":     name,   // 用户输入的标题
					"fuzziness": "AUTO", // 启用模糊匹配
				},
			},
		},
	}

	// 将查询条件转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	// 执行搜索
	res, err := dao.esClient.Search(
		dao.esClient.Search.WithIndex("knowledgebase"),
		dao.esClient.Search.WithBody(&buf),
		dao.esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("error executing search: %w", err)
	}
	defer res.Body.Close()

	// 检查响应状态
	if res.IsError() {
		return nil, fmt.Errorf("error in response: %s", res.String())
	}

	// 解析结果
	var result struct {
		Hits struct {
			Hits []struct {
				ID string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// 提取文档 ID
	var ids []string
	for _, hit := range result.Hits.Hits {
		ids = append(ids, hit.ID)
	}

	return ids, nil
}
