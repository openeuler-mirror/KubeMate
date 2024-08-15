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
package runner

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Runner struct {
}

func (r *Runner) RunCommand(command string) error {
	var out bytes.Buffer
	cmd := exec.Command("sh", "-c", command)
	multiWriter := io.MultiWriter(&out, os.Stdout)
	cmd.Stdout = multiWriter

	err := cmd.Run()
	if err != nil {
		logrus.Errorf("failed to run command %s: %v", command, err)
		return err
	}

	return nil
}

func (r *Runner) RunShell(shell string) error {
	tempFile, err := os.CreateTemp("/tmp/", "rear.sh")
	if err != nil {
		logrus.Errorf("failed to create temp file: %v", err)
		return err
	}
	defer os.Remove(tempFile.Name())

	// 将shell脚本写入临时文件
	if _, err := tempFile.WriteString(shell); err != nil {
		logrus.Errorf("failed to write to temp file: %v", err)
		return err
	}
	if err := tempFile.Close(); err != nil {
		logrus.Errorf("failed to close temp file: %v", err)
		return err
	}

	if err := os.Chmod(tempFile.Name(), 0555); err != nil {
		logrus.Errorf("failed to set execute permission on temp file: %v", err)
		return err
	}

	if err = r.RunCommand(tempFile.Name()); err != nil {
		return err
	}
	logrus.Info("backup os executed successfully!")

	return nil
}
