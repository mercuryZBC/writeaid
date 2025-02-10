package util

import (
	"fmt"
	"log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var client *oss.Client
var bucket *oss.Bucket

func init() {
	// 设置您的阿里云访问密钥和 OSS Endpoint

	accessKeyID := "LTAI5tLc6ntrhUU7k6VgyHeU"
	accessKeySecret := "y6oCunl7GAAMlgNgZhIPaZgRRS9aEp"
	endpoint := "https://oss-cn-beijing.aliyuncs.com" // 阿里云 OSS Endpoint，例如华东1（杭州）的 endpoint

	// 初始化 OSS 客户端
	var err error
	client, err = oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		log.Fatalf("Failed to create OSS client: %v", err)
	}

	// 设置目标 Bucket 名称
	bucketName := "bucket-2519"

	// 获取 Bucket 对象
	bucket, err = client.Bucket(bucketName)
	if err != nil {
		log.Fatalf("Failed to get bucket: %v", err)
	}
}

// 上传文件到 OSS
func uploadFile(filepath string, ossObjectKey string) error {
	localFile := filepath     // 本地文件路径
	ossObject := ossObjectKey // OSS 上的文件名

	// 上传本地文件到 OSS
	err := bucket.PutObjectFromFile(ossObject, localFile)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	fmt.Printf("Successfully uploaded %s to OSS\n", localFile)
	return nil
}

// 从 OSS 下载文件
func downloadFile(ossObjectKey string, storePath string) error {
	ossObject := ossObjectKey                   // OSS 上的文件名
	localFile := storePath + "/" + ossObjectKey // 本地文件路径

	// 下载 OSS 上的文件到本地
	err := bucket.GetObjectToFile(ossObject, localFile)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	fmt.Printf("Successfully downloaded %s from OSS\n", ossObject)
	return nil
}
