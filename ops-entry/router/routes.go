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
package router

import (
	"net/http"
	"ops-entry/controllers"
	_ "ops-entry/docs"
	"ops-entry/middlewares"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

type ErrResult struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data []string `json:"data"`
}

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // use release gin mode
	router := gin.New()

	router.Use(middlewares.Logger())   // logger middlerware
	router.Use(middlewares.Recovery()) // panic in single request instead of whole project
	router.GET("/", RootDirHandler)
	router.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	// api for file
	kubeconfigRouter := router.Group("/kubeconfig")
	{
		kubeconfigRouter.POST("/upload", controllers.KubeconfigFileUploadHandler)
		kubeconfigRouter.DELETE("/:cluster_id", controllers.KubeconfigFileDeleteHandler)
		kubeconfigRouter.GET("/:cluster_id", controllers.KubeconfigFileQueryHandler)
	}

	clusterConfigRouter := router.Group("/clusterconfig")
	{
		clusterConfigRouter.POST("/upload", controllers.ClusterconfigFileUploadHandler)
		clusterConfigRouter.DELETE("/:cluster_id", controllers.ClusterconfigFileDeleteHandler)
		clusterConfigRouter.GET("/:cluster_id", controllers.ClusterconfigFileQueryHandler)
	}

	return router
}

func RootDirHandler(c *gin.Context) {
	logrus.Infof("Not found pid=%d ppid=%d", os.Getpid(), os.Getppid())
	errResult := ErrResult{10001, "Not Found Data", nil}
	c.JSON(http.StatusOK, errResult)
	return
}
