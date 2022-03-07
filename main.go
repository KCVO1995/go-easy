package main

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/webview/webview"
	"github.com/skip2/go-qrcode"
)

//go:embed frontend/dist/*
var FS embed.FS

func ginFunc() {
	router := gin.Default()

	router.GET("/api/v1/addresses", AddressesController)
	router.GET("/uploads/:path", UploadsController)
	router.GET("/api/v1/qrcodes", QrcodesController)
	router.POST("/api/v1/texts", TextsController)

	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	router.StaticFS("/static", http.FS(staticFiles))

	router.NoRoute(func(context *gin.Context) {
		path := context.Request.URL.Path
		if strings.HasPrefix(path, "/static") {
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
	w.Navigate("http://127.0.0.1:8080/static")
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

func AddressesController(c *gin.Context) {
	// 获取电脑端所有 ip 地址
	// 通过 json 返回给前端
	addrs, _ := net.InterfaceAddrs()
	var result []string
	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && (ipNet.IP.IsLoopback()) {
			result = append(result, ipNet.IP.String())
		}
	}
	c.JSON(http.StatusOK, gin.H{"addresses": result})
}

func GetUploadsDir() (uploads string) {
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
		target := filepath.Join(GetUploadsDir(), path)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+path)
		c.Header("Content-Type", "application/octet-stream")
		c.File(target)
	} else {
		c.Status(http.StatusNotFound)
	}
}

func QrcodesController (c * gin.Context) {
	// 获取文本内容
	// 将文本转为图片
	// 将图片写入 HTTP 响应
	if content := c.Query("content"); content != "" {
		qrcode, err := qrcode.Encode(content, qrcode.Medium, 256)
		if err != nil {
			log.Fatal(err)
		}
		c.Data(http.StatusOK, "image/png", qrcode)
	} else {
		c.Status(http.StatusBadRequest)
	}
}