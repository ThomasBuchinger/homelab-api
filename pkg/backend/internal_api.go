package backend

import (
	"fmt"
	"net/http"
	"crypto/tls"
	"github.com/gin-gonic/gin"
	resty "github.com/go-resty/resty/v2"
)

func handleLegacyProxy(c *gin.Context) {
	url := c.Query("url")
	host_header := c.Query("host")

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "please specify a URL in the URL query parameter",
		})
  }

	req := resty.New().
		SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true }).
		SetRedirectPolicy(resty.NoRedirectPolicy()).
		R()
	if host_header != "na" {
		req = req.SetHeader("Host", host_header)
	}

	resp, err := req.Get(url)
	fmt.Println("GET ", url, " | ", resp.StatusCode())

	if resp.StatusCode() == 0 {
		fmt.Println("Network Error (", url,"):", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}
	c.Data(resp.StatusCode(), resp.Header().Get("content-type"), resp.Body())
}

func handleLegacyMetrics(c *gin.Context) {
	url := c.Query("url")

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "please specify a URL in the URL query parameter",
		})
  }

	req := resty.New().
		SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true }).
		R()

	resp, err := req.Get(url)
	fmt.Println("GET ", url, " | ", resp.StatusCode())

	if resp.StatusCode() == 0 {
		fmt.Println("Network Error (", url,"):", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}
	c.String(resp.StatusCode(), resp.String())
	
}
