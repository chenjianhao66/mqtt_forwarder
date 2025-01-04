package ui

import (
	"embed"
	"github.com/gin-gonic/gin"
	"net/http"
)

//go:embed dist/*
var staticFiles embed.FS

func RegisterRoutes(router *gin.Engine) {
	// 创建一个文件服务器，使用嵌入的静态文件
	fileServer := http.FileServer(http.FS(staticFiles))

	// 提供嵌入的文件服务
	router.NoRoute(func(c *gin.Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
