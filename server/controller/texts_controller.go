package controller

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TextsController(c *gin.Context) {
	// 1. 获取到文本
	// 2. 获取安装路径
	// 3. 生产一个文件，使用 hash 命名
	// 4. 返回下载链接
	var json struct {
		Row string `json:"raw"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		// 获取执行文件的路径
		exePath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		// 获取执行文件的目录路径
		dirPath := filepath.Dir(exePath)
		// 常见 uploads 文件夹
		uploads := filepath.Join(dirPath, "uploads")
		errMkdir := os.MkdirAll(uploads, os.ModePerm)
		if errMkdir != nil {
			log.Fatal(err)
		}

		// 拼接文件的绝对路径，不包含 exe 目录
		filename := uuid.New().String()
		fullPath := filepath.Join("uploads", filename+".txt")
		errWriteFile := ioutil.WriteFile(filepath.Join(dirPath, fullPath), []byte(json.Row), 0644)
		if errWriteFile != nil {
			log.Fatal(err)
		}

		// 返回文件路径
		c.JSON(http.StatusOK, gin.H{"url": "/" + fullPath})
	}
}