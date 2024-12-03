/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * KubeMate licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: lijian <lijian@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package middlewares

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		addr := c.Request.Header.Get("X-Real-IP")
		if addr == "" {
			addr = c.Request.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = c.Request.RemoteAddr
			}
		}

		lbAddr := c.Request.Header.Get("X-Original-To")
		requestAgent := c.Request.Header.Get("User-Agent")
		XProto := c.Request.Header.Get("X-Proto")

		c.Next()

		requestId := c.Request.Header.Get("Request-Id")
		logrus.Infof("[Logger and Request INFO]: %s, %s, %s, %s, %s, %s, %s, %v, %s, %v", requestId, XProto, c.Request.Method, c.Request.URL.Path, addr, lbAddr, requestAgent, c.Writer.Status(), http.StatusText(c.Writer.Status()), time.Since(start))
		//logrus.Flush()
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("[FATAL_PANIC] snow request[%s] panic: %s \n %s", c.Request.Header.Get("Request-Id"), err, debug.Stack())
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}
