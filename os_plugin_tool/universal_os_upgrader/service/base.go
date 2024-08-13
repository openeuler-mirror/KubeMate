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

package service

import (
	"github.com/sirupsen/logrus"
	"universal_os_upgrader/model"
)

func InitCmd() {
	univeralOSUpgradeCmd := model.NewUniversalOS()
	topCmd := univeralOSUpgradeCmd.RegisterEntryCmd()
	subCmdList := univeralOSUpgradeCmd.GetSubCmd()
	if len(subCmdList) < 1 {
		logrus.Error("empty subCmdList")
		return
	}
	for _, subCmd := range subCmdList {
		topCmd.AddCommand(subCmd)
	}

	if err := topCmd.Execute(); err != nil {
		return
	}
}