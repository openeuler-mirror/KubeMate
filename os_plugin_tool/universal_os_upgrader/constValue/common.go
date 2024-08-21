package constValue

type OperatorType string

const (
	Backup   OperatorType = "Backup"   // os 备份
	Upgrade  OperatorType = "Upgrade"  //os 升级
	Rollback OperatorType = "Rollback" //os 回滚

	BackupConfig = "/opt/kubemate/backup.yaml"
)
