package constValue

type OperatorType string

const (
	Upgrade  OperatorType = "Upgrade"  //os 升级
	Rollback OperatorType = "Rollback" //os 回滚
)
