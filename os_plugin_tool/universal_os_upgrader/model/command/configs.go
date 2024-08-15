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
package command

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func ReadConfigFile(opts *OptionsList) (*PluginConfig, error) {
	if opts.File == "" {
		logrus.Error("plugin config file path is empty")
		return nil, errors.New("plugin config file path empty")
	}

	fileData, err := os.ReadFile(opts.File)
	if err != nil {
		logrus.Errorf("failed to read config file %s :%v", opts.File, err)
		return nil, err
	}
	configData := &PluginConfig{}
	if err := yaml.Unmarshal(fileData, configData); err != nil {
		logrus.Errorf("failed to unmarshal json file:%v", err)
		return nil, err
	}
	return configData, nil
}
