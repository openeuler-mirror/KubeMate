package util

import (
	"fmt"
	"os"
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

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	return nil
}
