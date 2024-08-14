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
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// NKDDeployHandler
//
//	@Summary		Deploy a kubernetes cluster
//	@Description	Deploy a kubernetes cluster
//	@Tags			NKD管理Kubernetes集群
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			cluster_id	formData	string	true	"k8s cluster name"
//	@Param			labels      query		string	false	"The JSON string containing labels to filter the files to delete. Optional."
//	@Success		200			{object}	proto.NKDResult
//	@Router			/nkd/deploy	[POST]
func NKDDeployHandler(gc *gin.Context) {
	requestId := gc.GetHeader("Request-Id")
	c := util.CreateContext(requestId)
	if len(requestId) == 0 {
		gc.Request.Header.Set("Request-Id", c.RequestId)
	}

	var result proto.NKDResult
	result.Code = 0
	result.Msg = "success"
	result.RequestId = c.RequestId

	param := new(proto.NKDParam)
	err := gc.Bind(param)
	if err != nil {
		logrus.Errorf(c.P()+"Invalid params: %s", err.Error())
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = err.Error()
		gc.JSON(http.StatusOK, result)
		return
	}

	param.ClusterID = gc.PostForm("cluster_id")
	if param.ClusterID == "" {
		logrus.Errorf(c.P() + "Invalid param")
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = "Invalid param"
		gc.JSON(http.StatusOK, result)
		return
	}

	// 可选参数
	param.Labels = gc.Query("labels")

	byteParam, _ := json.Marshal(param)
	logrus.Infof(c.P()+"FileUploadHandler param: %s", string(byteParam))

	dst, err := util.GetSaveFilename(param.Labels, param.ClusterID)
	if err != nil {
		logrus.Errorf(c.P()+"Failed to get cluster config file: %s", err.Error())
		return
	}

	cmd := exec.Command(constValue.NkdPath, "deploy", "-f", dst)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf(c.P()+"Failed to deploy cluster: %s, output: %s", err.Error(), string(output))
		result.Code = util.ErrorCodeExecFail
		result.Msg = err.Error()
		gc.JSON(http.StatusOK, result)
		return
	}
	logrus.Infof(c.P()+"Cluster deployed successfully, output: %s", string(output))
	result.Msg = "Cluster deployed successfully"

	gc.JSON(http.StatusOK, result)
}

// NKDDeleteHandler
//
//	@Summary		Destroy a kubernetes cluster
//	@Description	Destroy a kubernetes cluster
//	@Tags			NKD管理Kubernetes集群
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			clusterID	formData	string	true	"kubernetes cluster name"
//	@Success		200			{object}	proto.NKDResult
//	@Router			/nkd/destroy [Delete]
func NKDDeleteHandler(gc *gin.Context) {
	requestId := gc.GetHeader("Request-Id")
	c := util.CreateContext(requestId)
	if len(requestId) == 0 {
		gc.Request.Header.Set("Request-Id", c.RequestId)
	}
	var result proto.NKDResult
	result.Code = 0
	result.Msg = "success"
	result.RequestId = c.RequestId

	param := new(proto.NKDParam)
	err := gc.Bind(param)
	if err != nil {
		logrus.Errorf(c.P()+"Invalid param: %s", err.Error())
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = err.Error()
		gc.JSON(http.StatusOK, result)
		return
	}

	param.ClusterID = gc.PostForm("clusterID")
	if param.ClusterID == "" {
		logrus.Errorf(c.P() + "Invalid param")
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = "Invalid param"
		gc.JSON(http.StatusOK, result)
		return
	}

	cmd := exec.Command(constValue.NkdPath, "destroy", "--cluster-id", param.ClusterID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf(c.P()+"Failed to destroy cluster: %s, output: %s", err.Error(), string(output))
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = err.Error()
		gc.JSON(http.StatusOK, result)
		return
	}
	logrus.Infof(c.P()+"Cluster destroyed successfully, output: %s", string(output))

	result.Msg = "Cluster destroyed successfully"
	gc.JSON(http.StatusOK, result)
}
