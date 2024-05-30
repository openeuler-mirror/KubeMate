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
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ops-entry/constValue"
	"ops-entry/log"
	router2 "ops-entry/router"
)

// @title 统一运维入口
// @version 1.0
// @description 接受运维管理平台的请求，进行os以云原生的方式进行升级
// @termsOfService http://swagger.io/terms/

// @contact.name https://gitee.com/openeuler/KubeMate
// @contact.url https://gitee.com/openeuler/KubeMate
// @contact.email https://gitee.com/openeuler/KubeMate

// @license.name Mulan PSL v2
// @license.url  http://license.coscl.org.cn/MulanPSL2

// @host 0.0.0.0:8080
// @BasePath /swagger/index.html

func main() {
	log.InitLog()
	router := router2.NewRouter()
	listen := fmt.Sprintf("%s:%d", constValue.ListenIP, constValue.ListenPort)
	logrus.Infof("Logger and gin inited, GinIsDebug[%v], listenIp[%s] listenPort[%d]", gin.IsDebugging(), constValue.ListenIP, constValue.ListenPort)

	err := router.Run(listen) // listen and serve on 0.0.0.0:8080
	if err != nil {
		logrus.Errorf("listen err: %s", err.Error())
	}
}
