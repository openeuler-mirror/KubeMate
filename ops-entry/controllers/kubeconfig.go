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
package controllers

import (
	"encoding/json"
	"net/http"
	"ops-entry/common/util"
	"ops-entry/constValue"
	"ops-entry/proto"
	"ops-entry/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// KubeconfigFileUploadHandler 上传kubeconfig文件
//
//	@Summary		Upload a kubeconfig file
//	@Description	Upload a file with optional description
//	@Tags			kubeconfig文件
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file		formData	file	true	"The kubeconfig file to upload"
//	@Param			cluster_id	formData	string	true	"k8s name"
//	@Success		200			{object}	proto.FileResult "Successful file upload"
//	@Router			/kubeconfig/upload [POST]
func KubeconfigFileUploadHandler(gc *gin.Context) {
	requestId := gc.GetHeader("Request-Id")
	c := util.CreateContext(requestId)
	if len(requestId) == 0 {
		gc.Request.Header.Set("Request-Id", c.RequestId)
	}
	var result proto.FileResult
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

	err = service.UploadKubeconfigFile(c, param)
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

// KubeconfigFileDeleteHandler 删除kubeconfig文件
//
//	@Summary		Delete a kubeconfig file
//	@Description	Delete a file with optional description
//	@Tags			kubeconfig文件
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path	string	true	"k8s name"
//	@Success		204			"No Content - Indicates successful deletion"
//	@Router			/kubeconfig/{cluster_id} [DELETE]
func KubeconfigFileDeleteHandler(gc *gin.Context) {
	requestId := gc.Param("cluster_id")
	c := util.CreateContext(requestId)
	if len(requestId) == 0 {
		gc.Request.Header.Set("Request-Id", c.RequestId)
	}
	service.DeleteKubeconfigFile(c, requestId)

	gc.Status(http.StatusNoContent)
	return
}
