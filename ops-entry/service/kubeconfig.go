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
package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"ops-entry/common/util"
	"ops-entry/constValue"
	"ops-entry/db/configManager"
	"ops-entry/db/configManager/config"
	"ops-entry/proto"
	"os"
	"os/user"
	"path/filepath"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func UploadKubeconfigFile(c util.Context, param *proto.FileUploadParam) error {
	allowedExts := []string{constValue.YamlExt, constValue.YmlExt, constValue.KubeconfigExt} // 允许的文件扩展名
	ext := filepath.Ext(param.File.Filename)
	if len(ext) > 0 && !util.IsAllowedExtension(ext, allowedExts) {
		logrus.Errorf(c.P()+"Not allowed file extension:%v", param.File.Filename)
		return errors.New("Not allowed file extension" + param.File.Filename)
	}

	open, err := param.File.Open()
	if err != nil {
		logrus.Errorf(c.P()+"Error opening file:%v", err)
		return err
	}

	content, err := io.ReadAll(open)
	if err != nil {
		logrus.Errorf(c.P()+"Error reading file:%v", err)
		return err
	}

	isValid := validKubeConfig(c, content)
	if !isValid {
		logrus.Errorf(c.P() + "Invalid KubeConfig file")
		return errors.New("Invalid KubeConfig file")
	}

	currentUser, err := user.Current()
	if err != nil {
		logrus.Errorf(c.P()+"Error getting current user:%v", err)
		return err
	}

	KubeConfigSavePath := filepath.Join(currentUser.HomeDir, constValue.KubeconfigSavePath)
	err = util.CreatePath(KubeConfigSavePath)
	if err != nil {
		logrus.Errorf(c.P()+"Error creating path:%v", err)
		return err
	}

	dst := filepath.Join(KubeConfigSavePath, fmt.Sprintf("%s-%s", param.ClusterId, param.File.Filename))

	outFile, err := os.Create(dst)
	if err != nil {
		logrus.Errorf(c.P()+"Error creating file:%v", err)
		return errors.New("Error creating file:" + err.Error())
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, open); err != nil {
		logrus.Errorf(c.P()+"Error copying file:%v", err)
		return errors.New("Error copying file:" + err.Error())
	}

	return saveKubeConfig2Secret(c, param.ClusterId, content, nil)
}

// validKubeConfig  校验kubeconfig是否合法
func validKubeConfig(c util.Context, content []byte) bool {
	clientSet, err := configManager.GetOuterClientSet(content)
	if err != nil {
		logrus.Errorf(c.P()+"GetOuterClientSet failed: %v\n", err)
		return false
	}

	_, err = clientSet.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		logrus.Errorf(c.P()+"Error checking Kubernetes connectivity: %v\n", err)
		return false
	}
	return true
}

func DeleteKubeconfigFile(c util.Context, clusterID string) error {
	secretName := clusterID + constValue.KubeconfigPrefix
	if !util.IsValidResourceName(secretName) {
		logrus.Errorf(c.P()+"invalid secret name: %s\n", secretName)
		return errors.New("invalid secret name")
	}
	sr := config.NewSecretImpl(constValue.NameSpace, secretName, nil, "")
	return sr.Delete(context.TODO(), metav1.DeleteOptions{})
}

func QueryKubeconfigFile(c util.Context, clusterID string) (*corev1.Secret, error) {
	secretName := clusterID + constValue.KubeconfigPrefix
	if !util.IsValidResourceName(secretName) {
		logrus.Errorf(c.P()+"invalid secret name: %s\n", secretName)
		return nil, errors.New("invalid secret name")
	}
	sr := config.NewSecretImpl(constValue.NameSpace, secretName, nil, "")
	return sr.Get(context.TODO(), metav1.GetOptions{})
}

// saveKubeConfig2Secret 保存kubeconfig到secret
func saveKubeConfig2Secret(c util.Context, clusterID string, kubeconfigBytes []byte, labelData map[string]string) error {
	encodedConfig := base64.StdEncoding.EncodeToString(kubeconfigBytes)
	kubeconfigData := map[string]string{
		constValue.Kubeconfig: string(encodedConfig),
	}
	secretName := clusterID + constValue.KubeconfigPrefix
	if !util.IsValidResourceName(secretName) {
		logrus.Errorf(c.P()+"invalid secret name: %s\n", secretName)
		return errors.New("invalid secret name")
	}
	sr := config.NewSecretImpl(constValue.NameSpace, secretName, labelData, corev1.ServiceAccountKubeconfigKey)
	return sr.Create(context.TODO(), metav1.CreateOptions{}, kubeconfigData)
}
