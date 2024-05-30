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

// KubernetesConfig 表示Kubernetes的配置
type KubernetesConfig struct {
	// Version 是Kubernetes的版本
	Version           string            `json:"version,omitempty"`
	ApiVersion        string            `json:"api_version,omitempty"`
	ApiserverEndpoint string            `json:"apiserver_endpoint,omitempty"`
	ImageRegistry     string            `json:"image_registry,omitempty"`
	PauseImage        string            `json:"pause_image,omitempty"`
	Token             string            `json:"token,omitempty"`
	Adminkubeconfig   string            `json:"adminkubeconfig,omitempty"`
	Certificatekey    string            `json:"certificatekey,omitempty"`
	Network           KubernetesNetwork `json:"network,omitempty"`
	// ControlPlaneComponents 包含了控制平面组件的配置（如API服务器、调度器等）
	ControlPlaneComponents ControlPlaneComponents `json:"controlPlaneComponents,omitempty"`
}

type KubernetesNetwork struct {
	ServiceSubnet string `json:"service_subnet,omitempty"`
	PodSubnet     string `json:"pod_subnet,omitempty"`
	Plugin        string `json:"plugin,omitempty"`
}

// ControlPlaneComponents 表示控制平面组件的配置
type ControlPlaneComponents struct {
	APIServer APIServerConfig `json:"api_server,omitempty"`
}

// APIServerConfig 表示API服务器的配置
type APIServerConfig struct {
}
