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
	"universal_os_upgrader/pkg/common"
	"universal_os_upgrader/pkg/utils/runner"

	"github.com/sirupsen/logrus"
)

type OSBackupImpl struct {
	OSBackupConfig
	r runner.Runner
}

type OSBackupConfig struct {
}

func NewOSBackup(r runner.Runner) *OSBackupImpl {
	return &OSBackupImpl{
		r: r,
	}
}

func (o *OSBackupImpl) CopyData() error {
	shell, err := common.GetRearShell(common.HandleBackup)
	if err != nil {
		logrus.Errorf("error to get backup shell file:%v", err)
		return err
	}
	if err := o.r.RunShell(shell); err != nil {
		return err
	}

	return nil
}