package backend

import (
	"crypto/tls"
	"log"
	"net/http"

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
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetRedirectPolicy(resty.NoRedirectPolicy()).
		R()
	if host_header != "na" {
		req = req.SetHeader("Host", host_header)
	}

	resp, err := req.Get(url)
	log.Println("GET ", url, " | ", resp.StatusCode())

	if resp.StatusCode() == 0 {
		log.Println("Network Error (", url, "):", err.Error())
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
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		R()

	resp, err := req.Get(url)
	log.Println("GET ", url, " | ", resp.StatusCode())

	if resp.StatusCode() == 0 {
		log.Println("Network Error (", url, "):", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}
	c.String(resp.StatusCode(), resp.String())
}

const (
	HeaderRealClientIP = "x-forwarded-for"
)

func handleAuthSimple(c *gin.Context) {
	real_ip := c.GetHeader(HeaderRealClientIP)

	if AuthByGeoip(real_ip, authConfig) || AllowInternalIp(real_ip) {
		log.Println("Success: Auth Simple")
		c.JSON(http.StatusOK, struct{}{})
	} else {
		log.Println("Failed: Auth Simple")
		c.JSON(http.StatusForbidden, struct{}{})
	}
}

func handleAuthWithCred(c *gin.Context) {
	real_ip := c.GetHeader(HeaderRealClientIP)
	user, pass, hasCredentials := c.Request.BasicAuth()

	if AuthByGeoip(real_ip, authConfig) && hasCredentials && AuthByCredentials(user, pass) && AuthByTempAllow("", "") {
		log.Println("Success: Auth Login")
		c.JSON(http.StatusOK, struct{}{})
	} else {
		log.Println("Failed: Auth Login")
		c.JSON(http.StatusForbidden, struct{}{})
	}
}
