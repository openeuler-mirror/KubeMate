/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * KubeMate licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: liukuo <liukuo@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package util

import (
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
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
