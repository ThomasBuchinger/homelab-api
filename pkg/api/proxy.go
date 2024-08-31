package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func ProxyWithBasicAuth(proxy_url, user, pass string, c *gin.Context) {
	req, err := http.NewRequestWithContext(c.Request.Context(), "GET", proxy_url, http.NoBody)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	req.SetBasicAuth(user, pass)

	response, err := http.DefaultClient.Do(req)
	if err != nil || response.StatusCode != http.StatusOK {
		common.GetServerConfig().RootLogger.Errorln("Failed calling Proxy URL: ", response.StatusCode, err.Error())
		c.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, map[string]string{})
}
