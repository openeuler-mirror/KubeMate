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
package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"universal_os_upgrader/constValue"
	"universal_os_upgrader/pkg/common"
	"universal_os_upgrader/pkg/utils/runner"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type OSBackupImpl struct {
	OSBackupConfig
}

type OSBackupConfig struct {
	NfsServer string `yaml:"nfs_server"`
	NfsPath   string `yaml:"nfs_path"`
}

func NewOSBackup() (*OSBackupImpl, error) {
	config, err := loadConfig(constValue.BackupConfig)
	if err != nil {
		logrus.Errorf("failed to load os backup config: %s", err)
		return nil, err
	}

	return &OSBackupImpl{
		OSBackupConfig: *config,
	}, nil
}

func loadConfig(filePath string) (*OSBackupConfig, error) {
	err := valideConfigFile(filePath)
	if err != nil {
		logrus.Errorf("failed to get os bakcup config: %s", err)
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("failed to read config file: %s", err)
		return nil, err
	}

	var config OSBackupConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		logrus.Errorf("failed to unmarshal config: %s", err)
		return nil, err
	}

	return &config, nil
}

func valideConfigFile(filePath string) error {
	// valide dir
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logrus.Errorf("failed to create directory %s: %s", dir, err)
			return err
		}
		logrus.Infof("Directory created: %s\n", dir)
	}

	// valide file
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			logrus.Errorf("failed to create file %s: %s", filePath, err)
			return err
		}
		defer file.Close()

		defaultConfig := OSBackupConfig{
			NfsServer: "",
			NfsPath:   "",
		}
		data, err := yaml.Marshal(&defaultConfig)
		if err != nil {
			logrus.Errorf("failed to marshal default config: %s", err)
			return err
		}

		_, err = file.Write(data)
		if err != nil {
			logrus.Errorf("failed to write default config to file %s: %s", filePath, err)
			return err
		}

		logrus.Infof("File created with default config: %s\n", filePath)
	}

	return nil
}

func (o *OSBackupImpl) RegisterSubCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "os backup",
		RunE:  o.RunBackupCmd,
	}
}

func (o *OSBackupImpl) RunBackupCmd(cmd *cobra.Command, args []string) error {
	if err := o.CopyData(); err != nil {
		return err
	}

	return nil
}

func (o *OSBackupImpl) CopyData() error {
	if err := o.validateParams(); err != nil {
		logrus.Errorf("Invalid param, %s", err.Error())
		return err
	}

	r := runner.Runner{}

	shell, err := common.GetRearShell("backup", o.NfsServer, o.NfsPath)
	if err != nil {
		logrus.Errorf("error to get backup shell file:%s", err)
		return err
	}

	if err := r.RunShell(shell); err != nil {
		return err
	}

	return nil
}

func (o *OSBackupImpl) validateParams() error {
	if o.NfsServer == "" {
		return fmt.Errorf("failed to get config param, nfs_server")
	}

	if o.NfsPath == "" {
		return fmt.Errorf("failed to get config param, nfs_path")
	}
	if !strings.HasSuffix(o.NfsPath, "/") {
		o.NfsPath += "/"
	}
	return nil
}
