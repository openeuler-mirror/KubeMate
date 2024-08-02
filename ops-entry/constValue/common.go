package constValue

type OperatorType int

const (
	Add OperatorType = iota + 1
	Delete
	Update
	Get
)

const MaxCopyFileSize = 10 * 1024 * 1024 // 10MB =

const (
	ListenIP   = "0.0.0.0"
	ListenPort = 9090
)

const NameSpace = "kubemate"
