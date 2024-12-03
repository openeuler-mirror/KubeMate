/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * KubeMate licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: liukuo <liukuo@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */

package util

const (
	ErrorCodeSuccess      = 0
	ErrorCodeRepeatReq    = 10  //	接到多个相同请求，只接受一个，其他返回该错误码
	ErrorCodeInvalidParam = 100 // 参数错误
	ErrorCodeUserNotValid = 101 // 用户非法
	ErrorCodeFail         = 200 // 失败
	ErrorCodeDbFail       = 201 // DB请求失败
	ErrorCodeExecFail     = 203 // 执行失败
	ErrorCodeTryAgain     = 300 // 失败，需要重试
	ErrorCodeUnknown      = 999 // 未知失败
)
