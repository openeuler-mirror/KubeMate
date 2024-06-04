package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"ops-entry/common/util"
	"ops-entry/constValue"
	"ops-entry/proto"
	"ops-entry/service"
)

// KubeconfigFileUploadHandler
// @Summary Upload a file
// @Description Upload a file with optional description
// @Tags kubeconfig文件上传
// @Accept  multipart/form-data
// @Produce  json
// @Param file formData file true "The file to upload"
// @Param cluster_id formData string true "k8s name"
// @Success 200 {object} proto.FileUploadResult
// @Router /file/upload/kubeconfig [POST]
func KubeconfigFileUploadHandler(gc *gin.Context) {
	requestId := gc.GetHeader("Request-Id")
	c := util.CreateContext(requestId)
	if len(requestId) == 0 {
		gc.Request.Header.Set("Request-Id", c.RequestId)
	}
	var result proto.FileUploadResult
	result.Code = 0
	result.Msg = "success"
	result.RequestId = c.RequestId

	param := new(proto.FileUploadParam)
	err := gc.Bind(param)
	if err != nil {
		logrus.Errorf(c.P()+"Invalid params: %s", err.Error())
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = err.Error()
		gc.JSON(http.StatusOK, result)
		return
	}

	param.File, err = gc.FormFile("file")

	// 将每次通过接口访问的参数记录下来，便于排查 问题。
	byteParam, _ := json.Marshal(param)
	logrus.Infof(c.P()+"FileUploadHandler param: %s", string(byteParam))

	if len(param.ClusterId) == 0 || !util.IsValidResourceName(param.ClusterId) {
		logrus.Errorf(c.P() + "Invalid params:Empty ClusterId or is not a valid resource")
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = "Invalid params:Empty ClusterId or is not a valid resource"
		gc.JSON(http.StatusOK, result)
		return
	}

	if param.File == nil {
		logrus.Errorf(c.P() + "Invalid params:Empty File")
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = "Invalid params:Empty File"
		gc.JSON(http.StatusOK, result)
		return
	}

	if param.File.Size > constValue.MaxFileSize || param.File.Size < 0 {
		logrus.Errorf(c.P() + "Invalid params:Big File Size or Zero File Size")
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = "Invalid params:Big File Size or Zero File Size"
		gc.JSON(http.StatusOK, result)
		return
	}

	err = service.UploadFile(c, param)
	if err != nil {
		logrus.Errorf(c.P()+"UploadFile failed: %s", err.Error())
		result.Code = util.ErrorCodeFail
		result.Msg = err.Error()
		gc.JSON(http.StatusOK, result)
		return
	}

	gc.JSON(http.StatusOK, result)
	return
}
