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
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// ClusterconfigFileUploadHandler 上传集群配置文件
//
//	@Summary		Upload a cluster config file
//	@Description	Upload a file with optional description
//	@Tags			集群配置文件
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file		formData	file	true	"The cluster config file to upload"
//	@Param			cluster_id	formData	string	true	"k8s name"
//	@Param 			type		formData 	string 	true 	"The type of the uploaded file (e.g., 'crfile', 'configfile')"
//	@Param 			labels	    formData    string  false	"The JSON string containing labels for the uploaded file"
//	@Success		200			{object}	proto.FileResult
//	@Router			/clusterconfig/upload [POST]
func ClusterconfigFileUploadHandler(gc *gin.Context) {
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

	if param.Type != proto.FileTypeCR && param.Type != proto.FileTypeFile {
		logrus.Errorf("Type field is required and cannot be empty")
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = "Type field is incorrect or empty"
		gc.JSON(http.StatusOK, result)
		return
	}

	param.File, err = gc.FormFile("file")
	if err != nil {
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = "Invalid parameters: file is empty"
		logrus.Errorf(c.P()+"Invalid parameters: file is empty: %v", err)
		gc.JSON(http.StatusOK, result)
		return
	}

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

	err = service.UploadClusterConfigFile(c, param)
	if err != nil {
		logrus.Errorf(c.P()+"UploadClusterConfigFile failed: %s", err.Error())
		result.Code = util.ErrorCodeFail
		result.Msg = err.Error()
		gc.JSON(http.StatusOK, result)
		return
	}

	gc.JSON(http.StatusOK, result)
	return
}

// ClusterconfigFileDeleteHandler 删除集群配置文件
//
//	@Summary		Delete a cluster config file
//	@Description	Delete a cluster config file with optional description
//	@Tags			集群配置文件
//	@Param			cluster_id	path	string	true	"k8s name"
//	@Success		204			"No Content - Indicates successful deletion"
//	@Router			/clusterconfig/{cluster_id} [DELETE]
func ClusterconfigFileDeleteHandler(gc *gin.Context) {
	requestId := gc.Param("cluster_id")
	c := util.CreateContext(requestId)
	if len(requestId) == 0 {
		gc.Request.Header.Set("Request-Id", c.RequestId)
	}
	service.DeleteClusterConfigFile(c, requestId)

	gc.Status(http.StatusNoContent)
	return
}

// ClusterconfigFileQueryHandler 查询集群配置文件
// @Summary 	Query a clusterconfig file
// @Description Query a clusterconfig file by cluster ID
// @Tags 		集群配置文件
// @Param 		cluster_id 	path 		string true "k8s cluster ID"
// @Success		200			{object}	proto.FileResult
// @Router /clusterconfig/{cluster_id} [GET]
func ClusterconfigFileQueryHandler(gc *gin.Context) {
	var (
		clusterConfigInfo proto.ClusterConfigResult
		result            proto.FileResult
	)
	requestId := gc.Param("cluster_id")
	c := util.CreateContext(requestId)
	if len(requestId) == 0 {
		gc.Request.Header.Set("Request-Id", c.RequestId)
	}
	result.Code = 0
	result.Msg = "success"
	result.RequestId = c.RequestId

	secret, err := service.QueryClusterConfigFile(c, requestId)
	if err != nil && !k8serrors.IsNotFound(err) {
		logrus.Errorf("failed to get secret to cluster_id: %s: %v", requestId, err)
		result.Code = util.ErrorCodeFail
		result.Msg = err.Error()
		gc.JSON(http.StatusOK, result)
		return
	}

	clusterConfigInfo.Name = secret.Name
	clusterConfigInfo.Data = string(secret.Data[constValue.Clusterconfig])
	gc.JSON(http.StatusOK, clusterConfigInfo)

	return
}
