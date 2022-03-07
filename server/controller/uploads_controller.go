package controller

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func getUploadsDir() (uploads string) {
	// 获取执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	// 获取执行文件的目录路径
	dirPath := filepath.Dir(exePath)
	// 获取 uploads 目录路径
	uploads = filepath.Join(dirPath, "uploads")
	return
}

func UploadsController(c *gin.Context) {
	// 将网络路径 :path 变成本地路径
	// 获取本地文件，写到 HTTP 响应里
	if path := c.Param("path"); path != "" {
		target := filepath.Join(getUploadsDir(), path)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+path)
		c.Header("Content-Type", "application/octet-stream")
		c.File(target)
	} else {
		c.Status(http.StatusNotFound)
	}
}