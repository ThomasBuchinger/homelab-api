package backend

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func handleServerConfig(c *gin.Context) {
	c.JSONP(http.StatusOK, common.GetServerConfig())
}

func handleClientConfig(c *gin.Context) {
	real_ip := c.ClientIP()
	ip := c.GetHeader("x-forwarded-for")
	if ip == "" {
		ip = real_ip
	}

	c.JSON(200, gin.H{
		"ip": ip,
		"real_ip": real_ip,
		"internal": common.IsIpAddressInternal(ip),
	})
}

func handlePublicHealth(c *gin.Context) {
	target := c.Query("target")
	static_data := map[string][]string{
		"Servers": []string{"No Issues!"},
		"Network": []string{"No Issues!"},
		"API": []string{"No Issues!"},
	}

	c.JSON(200, gin.H{
		"healthy": true,
		"messages": static_data[target],
	})
}
