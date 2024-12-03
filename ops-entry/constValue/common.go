/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * KubeMate licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: weihuanhuan <weihuanhuan@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package constValue

type OperatorType int

const (
	Add OperatorType = iota + 1
	Delete
	Update
	Get
)

const MaxCopyFileSize = 10 * 1024 * 1024 // 10MB =

const (
	ListenIP   = "0.0.0.0"
	ListenPort = 9090
)

const NameSpace = "kubemate"
