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
	"path/filepath"
	"reflect"

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

func LoadConfig[T any](filePath string) (*T, error) {
	err := validateConfigFile[T](filePath)
	if err != nil {
		logrus.Errorf("failed to validate config file: %s", err)
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("failed to read config file: %s", err)
		return nil, err
	}

	var config T
	if err := yaml.Unmarshal(data, &config); err != nil {
		logrus.Errorf("failed to unmarshal config: %s", err)
		return nil, err
	}

	return &config, nil
}

func validateConfigFile[T any](filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logrus.Errorf("failed to create directory %s: %s", dir, err)
			return err
		}
		logrus.Infof("Directory created: %s\n", dir)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			logrus.Errorf("failed to create file %s: %s", filePath, err)
			return err
		}
		defer file.Close()

		// create config file
		var defaultConfig T
		v := reflect.ValueOf(&defaultConfig).Elem()
		if v.Kind() == reflect.Struct {
			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)
				if field.Kind() == reflect.String {
					field.SetString("")
				}
			}
		} else {
			return errors.New("failed to create config file")
		}

		data, err := yaml.Marshal(&defaultConfig)
		if err != nil {
			logrus.Errorf("failed to marshal config: %s", err)
			return err
		}

		_, err = file.Write(data)
		if err != nil {
			logrus.Errorf("failed to write config file %s: %s", filePath, err)
			return err
		}

		logrus.Infof("Config file created: %s\n", filePath)
	}

	return nil
}
