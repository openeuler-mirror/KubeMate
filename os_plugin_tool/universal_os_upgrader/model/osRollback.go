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
	"io"
	"strings"
	"universal_os_upgrader/constValue"
	"universal_os_upgrader/model/command"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

type OSRollbackImpl struct {
	OSRollbackConfig
}

type OSRollbackConfig struct {
	NfsServer    string `yaml:"nfs_server"`
	NfsPath      string `yaml:"nfs_path"`
	Hostname     string `yaml:"hostname"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	IpxeServer   string `yaml:"ipxe_server"`
	IpxeRootPath string `yaml:"ipxe_root_path"`
	SshPort      string `yaml:"ssh_port"`
}

func NewOSRollback() (*OSRollbackImpl, error) {
	config, err := command.LoadConfig[OSRollbackConfig](constValue.RollbackConfig)
	if err != nil {
		logrus.Errorf("failed to load os rollback config: %s", err)
		return nil, err
	}

	return &OSRollbackImpl{
		OSRollbackConfig: *config,
	}, nil
}

func (o *OSRollbackImpl) RegisterSubCmd() *cobra.Command {
	rollbackCmd := &cobra.Command{
		Use:   string(constValue.Rollback),
		Short: "os rollback",
		RunE:  o.RunRollbackCmd,
	}

	return rollbackCmd
}

func (o *OSRollbackImpl) RunRollbackCmd(cmd *cobra.Command, args []string) error {
	if err := o.validateParams(); err != nil {
		return err
	}

	// os backup
	osbackup, err := NewOSBackup()
	if err != nil {
		logrus.Errorf("failed to execute OS backup: %v", err)
		return err
	}
	if err := osbackup.CopyData(); err != nil {
		return err
	}

	ipxeCfg := fmt.Sprintf(`#!ipxe
dhcp

set nfs_server %s
set hostname %s
set nfs_root %s${hostname}

kernel nfs://${nfs_server}${nfs_root}/${hostname}.kernel boot=live netboot=nfs nfsroot=${nfs_server}:${nfs_root}
initrd nfs://${nfs_server}${nfs_root}/${hostname}.initrd.cgz

boot
`, o.NfsServer, o.Hostname, o.NfsPath)

	client, err := o.connectToHost()
	if err != nil {
		fmt.Printf("Error connecting to host: %v\n", err)
		return nil
	}
	defer client.Close()

	if err := o.uploadToHost(client, ipxeCfg); err != nil {
		fmt.Printf("Error uploading content: %v\n", err)
	} else {
		fmt.Println("Content uploaded successfully!")
	}

	return nil
}

func (o *OSRollbackImpl) validateParams() error {
	if o.NfsServer == "" {
		return fmt.Errorf("failed to get param, nfs_server")
	}

	if o.NfsPath == "" {
		return fmt.Errorf("failed to get param, nfs_path")
	}
	if !strings.HasSuffix(o.NfsPath, "/") {
		o.NfsPath += "/"
	}

	if o.Hostname == "" {
		return fmt.Errorf("failed to get  param, hostname")
	}
	if o.User == "" {
		return fmt.Errorf("failed to get param, user")
	}
	if o.Password == "" {
		return fmt.Errorf("failed to get param, password")
	}
	if o.IpxeServer == "" {
		return fmt.Errorf("failed to get param, ipxe_server")
	}
	if o.IpxeRootPath == "" {
		return fmt.Errorf("failed to get param, ipxe_root_path")
	}
	if o.SshPort == "" {
		return fmt.Errorf("failed to get param, ssh_port")
	}
	return nil
}

func (o *OSRollbackImpl) connectToHost() (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: o.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(o.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf(o.IpxeServer + ":" + o.SshPort)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	return client, nil
}

func (o *OSRollbackImpl) uploadToHost(client *ssh.Client, fileContent string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	fileSize := len(fileContent)
	fileName := "ipxe.cfg"

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fmt.Fprintf(w, "C0644 %d %s\n", fileSize, fileName)
		io.Copy(w, strings.NewReader(fileContent))
		fmt.Fprint(w, "\x00")
	}()

	if !strings.HasSuffix(o.IpxeRootPath, "/") {
		o.IpxeRootPath += "/"
	}

	if err := session.Run(fmt.Sprintf("scp -t %s", o.IpxeRootPath+"ipxe.cfg")); err != nil {
		return fmt.Errorf("failed to run scp command: %v", err)
	}

	return nil
}
