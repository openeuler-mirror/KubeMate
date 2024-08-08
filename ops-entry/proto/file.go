package proto

import (
	"mime/multipart"
)

// FileResult 接口响应数据
type FileResult struct {
	BaseResult
}

// KubeConfigResult kubeconfig文件数据

type KubeConfigResult struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// ClusterConfigResult clusterconfig文件数据

type ClusterConfigResult struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type FileType string

const (
	FileTypeCR   FileType = "crfile"
	FileTypeFile FileType = "configfile"
)

// swagger:proto FileUploadParam
type FileUploadParam struct {
	ClusterId string                `json:"cluster_id" form:"cluster_id" example:"k8s-001" description:"The name of k8s"`
	Labels    string                `json:"labels" form:"labels" example:"{\"version\":\"v0.1\",\"environment\":\"prod\"}" description:"A JSON string representing labels for the uploaded file"`
	File      *multipart.FileHeader `json:"-" form:"file" swagger:"file" description:"The file to upload"`
	Type      FileType              `json:"type" form:"type" example:"cr" description:"The type of the uploaded file"`
}
