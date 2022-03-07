package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)
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

