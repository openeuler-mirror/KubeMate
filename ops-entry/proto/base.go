/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * KubeMate licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: lijian <lijian@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package proto

type BaseResult struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestId string `json:"request_id"`
}
