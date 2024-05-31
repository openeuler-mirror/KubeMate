package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
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

	router.GET("/", RootDirHandler)

	return router
}

func RootDirHandler(c *gin.Context) {
	logrus.Infof("Not found pid=%d ppid=%d", os.Getpid(), os.Getppid())
	errResult := ErrResult{10001, "Not Found Data", nil}
	c.JSON(http.StatusOK, errResult)
	return
}
