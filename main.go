package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/webview/webview"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed frontend/dist/*
var FS embed.FS

func ginFunc() {
	router := gin.Default()

	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	router.StaticFS("/", http.FS(staticFiles))

	router.POST("/interrupt", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
		os.Exit(1)
	})

	router.POST("/api/v1/texts", TextsController)
	router.NoRoute(func(context *gin.Context) {
		path := context.Request.URL.Path
		if !strings.HasPrefix(path, "/api") {
			file, err := staticFiles.Open("index.html")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			stat, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}
			context.DataFromReader(http.StatusOK, stat.Size(), "text/html", file, nil)
		} else {
			context.Status(http.StatusNotFound)
		}
	})
	router.Run()
}

func main() {
	go ginFunc()
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("GoEasy")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://127.0.0.1:8080")
	w.Run()
}

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

func NoRouteController() {

}