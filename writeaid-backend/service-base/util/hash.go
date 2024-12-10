package util

import (
	"crypto/sha256"
	"fmt"
	"os"
)

// 哈希函数，使用 SHA256 算法
func HashDocumentContent(filePath string) (string, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// 计算文件内容的 SHA256 哈希值
	hash := sha256.New()
	hash.Write(data)
	hashBytes := hash.Sum(nil)

	// 返回哈希值的十六进制字符串表示
	return fmt.Sprintf("%x", hashBytes), nil
}
