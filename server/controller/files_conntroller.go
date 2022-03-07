package controller

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)
func FilesController (c * gin.Context) {
	// 获取 go 执行文件所在目录
	// 在该目录常见 uploasd 目录
	// 将上传文件保存为另一个文件
	// 返回后者的下载路径
	// file = c.form
	file, err := c.FormFile("raw")
	if err != nil {
			log.Fatal(err)
	}
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
	fullPath := filepath.Join("uploads",  filename + filepath.Ext(file.Filename))
	errWriteFile := c.SaveUploadedFile(file, filepath.Join(dirPath, fullPath))
	if errWriteFile != nil {
		log.Fatal(err)
	}

	// 返回文件路径
	c.JSON(http.StatusOK, gin.H{"url": "/" + fullPath})	
}