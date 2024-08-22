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
	"fmt"

	"github.com/spf13/cobra"
)

type UniversalOS struct {
	FuncCmd       []*cobra.Command
	OSBackupImpl  *OSBackupImpl
	OSUpgradeImpl *OSUpgradeImpl
}

func NewUniversalOS() (*UniversalOS, error) {
	osBackupImpl, err := NewOSBackup()
	if err != nil {
		return nil, fmt.Errorf("failed to execute OS backup: %w", err)
	}

	return &UniversalOS{
		OSBackupImpl:  osBackupImpl,
		OSUpgradeImpl: NewOSUpgrade(),
	}, nil
}

func (uo *UniversalOS) RegisterEntryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "universal_os_upgrade",
		Short: "universal os upgrade tool",
	}

	uo.RegisterSubCmd(uo.OSBackupImpl.RegisterSubCmd())
	uo.RegisterSubCmd(uo.OSUpgradeImpl.RegisterSubCmd())
	return cmd
}

func (uo *UniversalOS) RegisterSubCmd(subCmd *cobra.Command) {
	uo.FuncCmd = append(uo.FuncCmd, subCmd)
}

func (uo *UniversalOS) GetSubCmd() []*cobra.Command {
	return uo.FuncCmd
}
