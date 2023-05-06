package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
)

// generateHash 生成数据hash
func GenerateHash(message []byte) string {
	bytes := sha256.Sum256(message)
	return hex.EncodeToString(bytes[:])
}

// generateFileHash 生成文件hash
func GenerateFileHash(filepath string) (string, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	bytes := sha256.Sum256(content)
	code := hex.EncodeToString(bytes[:])
	return code, nil
}
