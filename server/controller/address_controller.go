package controller

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)
func AddressesController(c *gin.Context) {
	// 获取电脑端所有 ip 地址
	// 通过 json 返回给前端
	addrs, _ := net.InterfaceAddrs()
	var result []string
	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				result = append(result, ipNet.IP.String())
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"addresses": result})
}

