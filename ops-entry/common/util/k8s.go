package util

import (
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

// IsValidResourceName 检查Kubernetes资源名称是否合法
func IsValidResourceName(name string) bool {
	if len(name) == 0 || len(name) > 153 {
		return false
	}

	validNameRegex := `^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	matched, err := regexp.MatchString(validNameRegex, name)
	if err != nil {
		logrus.Errorf("Error compiling regex: %v\n", err)
		return false
	}

	if !matched || strings.Contains(name, "_") {
		return false
	}
	return true
}
