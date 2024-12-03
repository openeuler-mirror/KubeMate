/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * KubeMate licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: liukuo <liukuo@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"ops-entry/constValue"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
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

func ParseLabels(labelData string) (map[string]string, error) {
	var labels map[string]string
	if err := json.Unmarshal([]byte(labelData), &labels); err != nil {
		logrus.Errorf("error unmarshal lables data:%v", labelData)
		return nil, errors.New("error unmarshal lables data:" + err.Error())
	}
	return labels, nil
}

// 获取集群配置文件存储名称
func GetSaveFilename(labelData string, clusterID string) (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		logrus.Errorf("Error getting current user:%v", err)
		return "", err
	}

	savePath := filepath.Join(currentUser.HomeDir, constValue.ClusterConfigSavePath, clusterID)
	err = CreatePath(savePath)
	if err != nil {
		logrus.Errorf("Error creating path:%v", err)
		return "", err
	}

	filename := clusterID
	if labelData != "" {
		labels, err := ParseLabels(labelData)
		if err != nil {
			return "", err
		}

		for key, value := range labels {
			filename = filename + "-" + key + "-" + value
		}
	}

	dst := filepath.Join(savePath, fmt.Sprintf("%s%s", filename, ".yaml"))
	return dst, nil
}
