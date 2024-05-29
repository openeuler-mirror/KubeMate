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

// LibvirtConfig 表示libvirt的配置
type LibvirtConfig struct {
	// URI 是libvirt的URI，用于连接到libvirtd守护进程
	URI      string           `json:"uri,omitempty"`
	Networks []LibvirtNetwork `json:"networks,omitempty"`
	// StoragePools 包含了libvirt管理的存储池列表
	StoragePools []LibvirtStoragePool `json:"storage_pools,omitempty"`
	OsImage      string               `json:"os_image,omitempty"`
	Cidr         string               `json:"cidr,omitempty"`
	Gateway      string               `json:"gateway,omitempty"`
}

// LibvirtNetwork 表示一个libvirt管理的网络
type LibvirtNetwork struct {
	Name string `json:"name,omitempty"`
}

// LibvirtStoragePool 表示一个libvirt管理的存储池
type LibvirtStoragePool struct {
	Name string `json:"name,omitempty"`
}
