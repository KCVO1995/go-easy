package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/KCVO1995/go-easy/server/controller"
	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var FS embed.FS

var Port = "27149"

func Run () {
	router := gin.Default()

	router.GET("/api/v1/addresses", controller.AddressesController)
	router.GET("/uploads/:path", controller.UploadsController)
	router.GET("/api/v1/qrcodes", controller.QrcodesController)
	router.POST("/api/v1/texts", controller.TextsController)
	router.POST("/api/v1/files", controller.FilesController)

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
	router.Run(":" + Port)
}








