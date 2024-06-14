package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
	"time"
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
