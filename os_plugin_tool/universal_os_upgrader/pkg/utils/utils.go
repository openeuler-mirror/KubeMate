/*
 * Copyright 2024 KylinSoft  Co., Ltd.
 * KubeMate is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
 * PURPOSE.
 * See the Mulan PSL v2 for more details.
 */
package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	//yum源基础路径
	RepoPath    = "/etc/yum.repos.d"
	NewRepoName = "upgrade.repo"
)

func RenameRepoFiles() error {
	files, err := os.ReadDir(RepoPath)
	if err != nil {
		logrus.Errorf("failed to read directory: %v", err)
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".repo") {
			oldPath := filepath.Join(RepoPath, file.Name())
			newPath := oldPath + ".bu"

			// 重命名文件
			err := os.Rename(oldPath, newPath)
			if err != nil {
				logrus.Errorf("failed to rename file %s to %s: %v", oldPath, newPath, err)
				return err
			}
		}
	}

	return nil
}
