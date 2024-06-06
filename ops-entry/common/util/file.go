package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IsAllowedExtension 检查扩展名是否在给定的允许列表中
func IsAllowedExtension(ext string, allowedExts []string) bool {
	ext = strings.ToLower(ext)
	for _, a := range allowedExts {
		if strings.ToLower(a) == ext {
			return true
		}
	}
	return false
}

// CreatePath 创建路径（如果它不存在）。
func CreatePath(path string) error {
	if !isValidPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	// 使用filepath.Dir获取路径的父目录
	parentDir := filepath.Dir(path)

	// 如果父目录不存在，则递归创建它
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return err
		}
	}

	return nil
}
