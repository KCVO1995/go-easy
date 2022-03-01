package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/webview/webview"
	"io/fs"
	"log"
	"net/http"
	"os"
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
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
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
