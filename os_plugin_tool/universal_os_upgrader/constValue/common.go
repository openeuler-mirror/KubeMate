/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * KubeMate licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: lijian <lijian@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package constValue

type OperatorType string

const (
	Backup   OperatorType = "Backup"   // os 备份
	Upgrade  OperatorType = "Upgrade"  //os 升级
	Rollback OperatorType = "Rollback" //os 回滚

	BackupConfig   = "/opt/kubemate/config/backup.yaml"
	UpgradeConfig  = "/opt/kubemate/config/upgrade.yaml"
	RollbackConfig = "/opt/kubemate/config/rollback.yaml"
)
