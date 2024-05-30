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

type OpenStackConfig struct {
	UserName         string `json:"user_name,omitempty"`
	Password         string `json:"password,omitempty"`
	TenantName       string `json:"tenant_name,omitempty"`
	AuthURL          string `json:"auth_url,omitempty"`
	Region           string `json:"region,omitempty"`
	InternalNetwork  string `json:"internal_network,omitempty"`
	ExternalNetwork  string `json:"external_network,omitempty"`
	GlanceName       string `json:"glance_name,omitempty"`
	AvailabilityNone string `json:"availability_none,omitempty"`
}
