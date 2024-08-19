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

package proto

type NKDResult struct {
	BaseResult
}

type NKDParam struct {
	ClusterID string `json:"cluster_id" binding:"required" form:"cluster_id" example:"cluster" description:"Name of a kubernetes cluster"`
}

type NKDDeployParam struct {
	NKDParam
	Labels string `json:"labels" form:"labels" example:"{\"version\":\"v0.1\"}" description:"A JSON string representing labels for the kubernetes cluster"`
}

type NKDDestroyParam struct {
	NKDParam
}

type NKDExtendParam struct {
	NKDParam
	Num string `json:"num" binding:"required" form:"num" example:"1" description:"Number of nodes to extend"`
}
