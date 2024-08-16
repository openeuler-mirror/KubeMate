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
	"os"
	"path/filepath"
	"universal_os_upgrader/model/command"
	"universal_os_upgrader/pkg/utils"
	"universal_os_upgrader/pkg/utils/runner"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	yumCleanall  = "yum clean all"
	yumMakecache = "yum makecache"
	yumUpdate    = "yum update -y"
)

type OSUpgradeImpl struct {
	OSUpgradeConfig
}

type OSUpgradeConfig struct {
	Image   string `json:"image"`
	Version string `json:"version"`
}

func NewOSUpgrade() *OSUpgradeImpl {
	return &OSUpgradeImpl{}
}

func (o *OSUpgradeImpl) RegisterSubCmd() *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:   "upgrade",
		Short: "os upgrade",
		RunE:  RunUpgradeCmd,
	}
	command.SetupUpgradeCmdOpts(upgradeCmd)
	return upgradeCmd
}

func RunUpgradeCmd(cmd *cobra.Command, args []string) error {
	r := runner.Runner{}
	//OS备份
	osbackup := NewOSBackup(r)
	if err := osbackup.CopyData(); err != nil {
		return err
	}

	configData, err := command.ReadConfigFile(&command.Opts)
	if err != nil {
		return err
	}
	//repo源文件备份
	if err := utils.RenameRepoFiles(); err != nil {
		return err
	}

	//repo源文件更新
	repoFile := filepath.Join(utils.RepoPath, utils.NewRepoName)
	err = os.WriteFile(repoFile, []byte(configData.Repo), 0644)
	if err != nil {
		logrus.Errorf("error writing repo file: %v\n", err)
		return err
	}

	logrus.Info("Starting upgrade...")
	if err := r.RunCommand(yumCleanall); err != nil {
		logrus.Errorf("failed to run command %s", yumCleanall)
		return err
	}

	if err := r.RunCommand(yumMakecache); err != nil {
		logrus.Errorf("failed to run command %s", yumMakecache)
		return err
	}

	if err := r.RunCommand(yumUpdate); err != nil {
		logrus.Errorf("failed to run command %s", yumUpdate)
		return err
	}
	logrus.Info("Upgrade successfully!")

	return nil
}
