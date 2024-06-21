package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"net/http"
	"ops-entry/controllers"
	_ "ops-entry/docs"
	"ops-entry/middlewares"
	"os"
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
	fileRouter := router.Group("/file")
	{
		fileRouter.POST("/upload/kubeconfig", controllers.KubeconfigFileUploadHandler)
	}
	return router
}

func RootDirHandler(c *gin.Context) {
	logrus.Infof("Not found pid=%d ppid=%d", os.Getpid(), os.Getppid())
	errResult := ErrResult{10001, "Not Found Data", nil}
	c.JSON(http.StatusOK, errResult)
	return
}
