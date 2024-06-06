package constValue

const (
	RootDir              = "~"
	KubeConfig           = RootDir + "/.kube/config" //kubectl config view
	DefaultNameSpace     = "default"                 //kubectl config view
	DefaultCrVersion     = "v1"
	DefaultCrUpdateField = "spec"
)
