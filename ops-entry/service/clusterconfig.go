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
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"ops-entry/common/util"
	"ops-entry/constValue"
	"ops-entry/db/configManager/config"
	"ops-entry/proto"
	"os"
	"os/user"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func UploadClusterConfigFile(c util.Context, param *proto.FileUploadParam) error {
	var labels map[string]string

	allowedExts := []string{constValue.YamlExt, constValue.YmlExt}
	ext := filepath.Ext(param.File.Filename)

	// 检查文件扩展名
	if len(ext) > 0 && !util.IsAllowedExtension(ext, allowedExts) {
		logrus.Errorf(c.P()+"Not allowed file extension:%v", param.File.Filename)
		return errors.New("Not allowed file extension" + param.File.Filename)
	}

	open, err := param.File.Open()
	if err != nil {
		logrus.Errorf(c.P()+"Error opening file:%v", err)
		return err
	}
	defer open.Close()

	content, err := io.ReadAll(open)
	if err != nil {
		logrus.Errorf(c.P()+"Error reading file:%v", err)
		return err
	}

	isValid := validYamlConfig(c, content)
	if !isValid {
		return errors.New("error: invalid yaml config file")
	}

	currentUser, err := user.Current()
	if err != nil {
		logrus.Errorf(c.P()+"Error getting current user:%v", err)
		return err
	}

	clusterConfigSavePath := filepath.Join(currentUser.HomeDir, constValue.ClusterConfigSavePath, param.ClusterId)
	err = util.CreatePath(clusterConfigSavePath)
	if err != nil {
		logrus.Errorf(c.P()+"Error creating path:%v", err)
		return err
	}

	str := param.ClusterId
	if param.Labels != "" {
		labels, err = parseLabels(param.Labels)
		if err != nil {
			return err
		}

		for key, value := range labels {
			str = str + "-" + key + "-" + value
		}
	}
	dst := filepath.Join(clusterConfigSavePath, fmt.Sprintf("%s%s", str, ".yaml"))
	outFile, err := os.Create(dst)
	if err != nil {
		logrus.Errorf(c.P()+"Error creating file:%v", err)
		return errors.New("Error creating file:" + err.Error())
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, bytes.NewReader(content)); err != nil {
		logrus.Errorf(c.P()+"Error copying file:%v", err)
		return errors.New("Error copying file:" + err.Error())
	}

	if param.Type == proto.FileTypeCR {
		if err := applyCRResource(content, labels); err != nil {
			return err
		}
	}

	return saveClusterConfig2Secret(c, param.ClusterId, content, labels)
}

func DeleteClusterConfigFile(c util.Context, clusterID string, labels string) error {
	str := clusterID
	secretName := clusterID + constValue.ClusterconfigPrefix
	if labels != "" {
		deleteLabels, err := parseLabels(labels)
		if err != nil {
			return err
		}

		for key, value := range deleteLabels {
			secretName = secretName + "-" + key + "-" + value
			str = str + "-" + key + "-" + value
		}
	}
	currentUser, err := user.Current()
	if err != nil {
		logrus.Errorf(c.P()+"Error getting current user:%v", err)
		return err
	}

	clusterConfigSavePath := filepath.Join(currentUser.HomeDir, constValue.ClusterConfigSavePath, clusterID)
	dst := filepath.Join(clusterConfigSavePath, fmt.Sprintf("%s%s", str, ".yaml"))
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		// 如果文件存在，则尝试删除
		if err := os.Remove(dst); err != nil {
			logrus.Errorf(c.P()+"Error deleting file:%s\n", err)
		}
	}

	if !util.IsValidResourceName(secretName) {
		logrus.Errorf(c.P()+"invalid secret name: %s\n", secretName)
		return errors.New("invalid secret name")
	}
	sr := config.NewSecretImpl(constValue.NameSpace, secretName, nil, "")
	return sr.Delete(context.TODO(), metav1.DeleteOptions{})
}

func QueryClusterConfigFile(c util.Context, clusterID string, labels string) (*corev1.Secret, error) {
	secretName := clusterID + constValue.ClusterconfigPrefix
	if labels != "" {
		deleteLabels, err := parseLabels(labels)
		if err != nil {
			return nil, err
		}

		for key, value := range deleteLabels {
			secretName = secretName + "-" + key + "-" + value
		}
	}

	if !util.IsValidResourceName(secretName) {
		logrus.Errorf(c.P()+"invalid secret name: %s\n", secretName)
		return nil, errors.New("invalid secret name")
	}
	sr := config.NewSecretImpl(constValue.NameSpace, secretName, nil, "")
	return sr.Get(context.TODO(), metav1.GetOptions{})
}

func UpdateClusterConfigFile(c util.Context, param *proto.FileUpdateParam) error {
	var updateLabels map[string]string
	allowedExts := []string{constValue.YamlExt, constValue.YmlExt}
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
	defer open.Close()

	content, err := io.ReadAll(open)
	if err != nil {
		logrus.Errorf(c.P()+"Error reading file:%v", err)
		return err
	}

	isValid := validYamlConfig(c, content)
	if !isValid {
		return errors.New("error: invalid yaml config file")
	}

	currentUser, err := user.Current()
	if err != nil {
		logrus.Errorf(c.P()+"Error getting current user:%v", err)
		return err
	}

	savePath := filepath.Join(currentUser.HomeDir, constValue.ClusterConfigSavePath, param.ClusterId)
	err = util.CreatePath(savePath)
	if err != nil {
		logrus.Errorf(c.P()+"Error creating path:%v", err)
		return err
	}

	str := param.ClusterId
	if param.Labels != "" {
		updateLabels, err = parseLabels(param.Labels)
		if err != nil {
			return err
		}

		for key, value := range updateLabels {
			str = str + "-" + key + "-" + value
		}
	}
	dst := filepath.Join(savePath, fmt.Sprintf("%s%s", str, ".yaml"))

	outFile, err := os.Create(dst)
	if err != nil {
		logrus.Errorf(c.P()+"Error creating file:%v", err)
		return errors.New("Error creating file:" + err.Error())
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, bytes.NewReader(content)); err != nil {
		logrus.Errorf(c.P()+"Error copying file:%v", err)
		return errors.New("Error copying file:" + err.Error())
	}

	return updateClusterConfig2Secret(c, param.ClusterId, content, updateLabels)
}

// validYamlConfig 检查给定的内容是否是有效的 YAML 配置文件
func validYamlConfig(c util.Context, content []byte) bool {
	var data interface{}
	err := yaml.Unmarshal(content, &data)
	if err != nil {
		logrus.Errorf(c.P()+"invalid yaml: %v\n", err)
		return false
	}
	return true
}

func parseLabels(labelData string) (map[string]string, error) {
	var labels map[string]string
	if err := json.Unmarshal([]byte(labelData), &labels); err != nil {
		logrus.Errorf("error unmarshal lables data:%v", labelData)
		return nil, errors.New("error unmarshal lables data:" + err.Error())
	}
	return labels, nil
}

func saveClusterConfig2Secret(c util.Context, clusterID string, configBytes []byte, labelData map[string]string) error {
	var secretName string
	encodedConfig := base64.StdEncoding.EncodeToString(configBytes)
	clusterConfigData := map[string]string{
		constValue.Clusterconfig: string(encodedConfig),
	}
	secretName = clusterID + constValue.ClusterconfigPrefix
	if len(labelData) > 0 {
		for key, value := range labelData {
			secretName = secretName + "-" + key + "-" + value
		}
	}
	if !util.IsValidResourceName(secretName) {
		logrus.Errorf(c.P()+"invalid secret name: %s\n", secretName)
		return errors.New("invalid secret name")
	}
	sr := config.NewSecretImpl(constValue.NameSpace, secretName, labelData, "")
	return sr.Create(context.TODO(), metav1.CreateOptions{}, clusterConfigData)
}

func updateClusterConfig2Secret(c util.Context, clusterID string, configBytes []byte, labelData map[string]string) error {
	var secretName string
	encodedConfig := base64.StdEncoding.EncodeToString(configBytes)
	clusterConfigData := map[string]string{
		constValue.Clusterconfig: string(encodedConfig),
	}
	secretName = clusterID + constValue.ClusterconfigPrefix
	if len(labelData) > 0 {
		for key, value := range labelData {
			secretName = secretName + "-" + key + "-" + value
		}
	}
	if !util.IsValidResourceName(secretName) {
		logrus.Errorf(c.P()+"invalid secret name: %s\n", secretName)
		return errors.New("invalid secret name")
	}
	sr := config.NewSecretImpl(constValue.NameSpace, secretName, labelData, "")
	return sr.Update(context.TODO(), metav1.UpdateOptions{}, clusterConfigData)
}

// 应用Cr资源
func applyCRResource(clusterconfigBytes []byte, labelData map[string]string) error {
	cr := config.NewCrImpl(
		"group",
		constValue.DefaultCrVersion,
		"kind",
		"resource",
		constValue.NameSpace,
		"crName",
		constValue.DefaultCrUpdateField,
		labelData,
	)
	err := cr.Create(context.TODO(), metav1.CreateOptions{}, clusterconfigBytes)
	if err != nil {
		logrus.Errorf("failed to apply CR resource: %v", err)
		return errors.New("failed to apply CR resource:" + err.Error())
	}

	return nil
}
