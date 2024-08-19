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
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// NKDDeployHandler
//
//	@Summary		Deploy a kubernetes cluster
//	@Description	Deploy a kubernetes cluster
//	@Tags			Use NKD to manage a kubernetes cluster
//	@Accept			application/json
//	@Produce		json
//	@Param			deploy			body		proto.NKDDeployParam	true	"Deploy a kubernetes cluster"
//	@Success		200				{object}	proto.NKDResult
//	@Router			/nkd/deploy		[POST]
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

	var requestBody proto.NKDDeployParam
	if err := gc.ShouldBindJSON(&requestBody); err != nil {
		logrus.Errorf(c.P()+"Invalid param: %s", err.Error())
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = err.Error()
		gc.JSON(http.StatusOK, result)
		return
	}

	if requestBody.ClusterID == "" {
		logrus.Errorf(c.P() + "Invalid param: clusterID is missing")
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = "Invalid param: clusterID is missing"
		gc.JSON(http.StatusOK, result)
		return
	}

	byteParam, _ := json.Marshal(requestBody)
	logrus.Infof(c.P()+"FileUploadHandler param: %s", string(byteParam))

	dst, err := util.GetSaveFilename(requestBody.Labels, requestBody.ClusterID)
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
//	@Tags			Use NKD to manage a kubernetes cluster
//	@Accept			application/json
//	@Produce		json
//	@Param			destroy			body		proto.NKDDestroyParam	true	"Destroy a kubernetes cluster"
//	@Success		200				{object}	proto.NKDResult
//	@Router			/nkd/destroy 	[Delete]
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

	var requestBody proto.NKDDestroyParam
	if err := gc.ShouldBindJSON(&requestBody); err != nil {
		logrus.Errorf(c.P()+"Invalid param: %s", err.Error())
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = err.Error()
		gc.JSON(http.StatusOK, result)
		return
	}

	if requestBody.ClusterID == "" {
		logrus.Errorf(c.P() + "Invalid param")
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = "Invalid param"
		gc.JSON(http.StatusOK, result)
		return
	}

	cmd := exec.Command(constValue.NkdPath, "destroy", "--cluster-id", requestBody.ClusterID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf(c.P()+"Failed to destroy cluster: %s, output: %s", err.Error(), string(output))
		result.Code = util.ErrorCodeExecFail
		result.Msg = err.Error()
		gc.JSON(http.StatusOK, result)
		return
	}
	logrus.Infof(c.P()+"Cluster destroyed successfully, output: %s", string(output))

	result.Msg = "Cluster destroyed successfully"
	gc.JSON(http.StatusOK, result)
}

// NKDExtendHandler
//
//	@Summary		Extend a kubernetes cluster
//	@Description	Extend a kubernetes cluster
//	@Tags			Use NKD to manage a kubernetes cluster
//	@Accept			application/json
//	@Produce		json
//	@Param			extend		body		proto.NKDExtendParam	true	"Extend a kubernetes cluster"
//	@Success		200			{object}	proto.NKDResult
//	@Router			/nkd/extend [POST]
func NKDExtendHandler(gc *gin.Context) {
	requestId := gc.GetHeader("Request-Id")
	c := util.CreateContext(requestId)
	if len(requestId) == 0 {
		gc.Request.Header.Set("Request-Id", c.RequestId)
	}
	var result proto.NKDResult
	result.Code = 0
	result.Msg = "success"
	result.RequestId = c.RequestId

	var requestBody proto.NKDExtendParam
	if err := gc.ShouldBindJSON(&requestBody); err != nil {
		logrus.Errorf(c.P()+"Invalid param: %s", err.Error())
		result.Code = util.ErrorCodeFail
		result.Msg = err.Error()
		gc.JSON(http.StatusOK, result)
		return
	}

	if requestBody.ClusterID == "" {
		logrus.Errorf(c.P() + "Invalid param")
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = "Invalid param"
		gc.JSON(http.StatusOK, result)
		return
	}

	num, err := strconv.Atoi(requestBody.Num)
	if err != nil || num <= 0 {
		logrus.Errorf(c.P()+"Invalid num: %s", err.Error())
		result.Code = util.ErrorCodeInvalidParam
		result.Msg = "Invalid num"
		gc.JSON(http.StatusOK, result)
		return
	}

	cmd := exec.Command(constValue.NkdPath, "extend", "--cluster-id", requestBody.ClusterID, "-n", requestBody.Num)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf(c.P()+"Failed to extend cluster: %s, output: %s", err.Error(), string(output))
		result.Code = util.ErrorCodeExecFail
		result.Msg = err.Error()
		gc.JSON(http.StatusOK, result)
		return
	}
	logrus.Infof(c.P()+"Cluster extended successfully, output: %s", string(output))

	result.Msg = "Cluster extended successfully"
	gc.JSON(http.StatusOK, result)
}
