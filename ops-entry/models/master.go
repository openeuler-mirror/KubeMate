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

package models

// KubernetesMasterNode 表示一个Kubernetes master节点
type KubernetesMasterNode struct {
	// Name 是master节点的名称
	Name string `json:"name,omitempty"`
	// IP 是master节点的IP地址
	IP string `json:"ip,omitempty"`
	// Port 是Kubernetes API服务的端口
	Port int `json:"port,omitempty"`
	// LibvirtConfig 包含了libvirt相关的配置
	LibvirtConfig LibvirtConfig `json:"libvirtConfig,omitempty"`
	// OpenStackConfig 包含了OpenStack相关的配置
	OpenStack OpenStackConfig `json:"open_stack,omitempty"`
	// KubernetesConfig 包含了Kubernetes相关的配置
	KubernetesConfig KubernetesConfig `json:"kubernetesConfig,omitempty"`
}
