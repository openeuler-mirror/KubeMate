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
	"universal_os_upgrader/pkg/common"
	"universal_os_upgrader/pkg/utils/runner"

	"github.com/sirupsen/logrus"
)

type OSBackupImpl struct {
	OSBackupConfig
}

type OSBackupConfig struct {
}

func NewOSBackup() *OSBackupImpl {
	return &OSBackupImpl{}
}

func (o *OSBackupImpl) CopyData() error {
	shell, err := common.GetRearShell(common.HandleBackup)
	if err != nil {
		logrus.Errorf("error to get backup shell file:%v", err)
		return err
	}
	r := &runner.Runner{}
	_, err = r.RunCommand(shell)
	if err != nil {
		return err
	}

	return nil
}
