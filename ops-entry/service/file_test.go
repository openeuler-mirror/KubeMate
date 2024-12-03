/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * KubeMate licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: lijian <lijian@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package service

import (
	"fmt"
	"mime/multipart"
	"ops-entry/common/util"
	"ops-entry/constValue"
	"ops-entry/proto"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	//var c util.Context
	param := &proto.FileUploadParam{}
	param.ClusterId = "cn.beijing.data02.k8s001"
	data := map[string][]string{
		"Accept-Encoding": []string{"gzip", "deflate", "br"},
	}
	param.File = &multipart.FileHeader{
		Filename: "config",
		Size:     1000,
		Header:   data,
	}
	allowedExts := []string{constValue.YamlExt, constValue.YmlExt, constValue.Kubeconfig} // 允许的文件扩展名

	//curl -F "file=@/path/to/your/kubeconfig" http://localhost:9090/upload
	t.Run("UploadFileExtension", func(t *testing.T) {
		ext := filepath.Ext(param.File.Filename)
		if len(ext) == 0 {
			t.Logf("success filename: %s", param.File.Filename)
		}
		param.File.Filename = "config.kubeconfig"
		ext = filepath.Ext(param.File.Filename)
		status := util.IsAllowedExtension(ext, allowedExts)
		if !status {
			t.Errorf("UploadFileFailed file: %s", param.File.Filename)
		} else {
			t.Logf("success filename: %s", param.File.Filename)
		}
		param.File.Filename = "config.cnf"
		ext = filepath.Ext(param.File.Filename)
		status = util.IsAllowedExtension(ext, allowedExts)
		if !status {
			t.Errorf("UploadFileFailed file: %s", param.File.Filename)
		} else {
			t.Logf("success filename: %s", param.File.Filename)
		}
	})

	t.Run("UploadFileSavePath", func(t *testing.T) {
		param.File.Filename = "config.kubeconfig"
		patches := gomonkey.ApplyFunc(validKubeConfig, func(c util.Context, open multipart.File) bool {
			return false
		})
		defer patches.Reset() // 在测试结束后恢复原始函数
		currentUser, err := user.Current()
		if err != nil {
			t.Logf("fail to get current user: %v", err)
		}

		KubeConfigSavePath := filepath.Join(currentUser.HomeDir, constValue.KubeconfigSavePath)
		dst := filepath.Join(KubeConfigSavePath, fmt.Sprintf("%s-%s", param.ClusterId, param.File.Filename))

		expectPath := filepath.Join(currentUser.HomeDir, constValue.KubeconfigSavePath, param.ClusterId+".config")
		assert.Equal(t, expectPath, dst)
	})
}
