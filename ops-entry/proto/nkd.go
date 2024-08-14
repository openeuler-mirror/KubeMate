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
	ClusterID string `json:"cluster_id" form:"cluster_id" example:"k8s-cluster" description:"k8s cluster name"`
	Labels    string `json:"labels" form:"labels" example:"{\"version\":\"v0.1\",\"environment\":\"prod\"}" description:"A JSON string representing labels for the k8s cluster"`
}
