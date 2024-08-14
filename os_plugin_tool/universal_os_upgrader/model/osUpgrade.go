/*
 *
 * Copyright 2024 KylinSoft  Co., Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * /
 */

package model

import (
	"github.com/spf13/cobra"
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
	return &cobra.Command{
		Use:   "upgrade",
		Short: "os upgrade",
		RunE:  RunUpgradeCmd,
	}
}

func RunUpgradeCmd(cmd *cobra.Command, args []string) error {
	osbackup := NewOSBackup()
	if err := osbackup.CopyData(); err != nil {
		return err
	}

	return nil
}
