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
	ListenPort = 8080
)

const NameSpace = "kubemate"
