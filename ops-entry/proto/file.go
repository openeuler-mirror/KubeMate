package proto

import (
	"mime/multipart"
)

// FileUploadResult kubeconfig上传接口响应数据
type FileUploadResult struct {
	BaseResult
}

// swagger:proto FileUploadParam
type FileUploadParam struct {
	ClusterId string                `json:"cluster_id" form:"cluster_id" example:"k8s-001" description:"The name of k8s"`
	File      *multipart.FileHeader `json:"-" form:"file" swagger:"file" description:"The file to upload"`
}
