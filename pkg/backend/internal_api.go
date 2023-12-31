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
	ApiLogger.Debugf("Called API: GET %s | %v", url, resp.StatusCode)

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
	ApiLogger.Debugf("Called API: GET %s | %v", url, resp.StatusCode)

	if resp.StatusCode() == 0 {
		log.Println("Network Error (", url, "):", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}
	c.String(resp.StatusCode(), resp.String())
}

// func handleAuthSimple(c *gin.Context) {
// 	ip := c.GetHeader(HeaderRealClientIP)
// 	ApiLogger.Debugf("Authentication Request: Remote: %s | Host: %v | Path: %v | Headers: %v", c.RemoteIP(), c.Request.Host, c.Request.URL.Path, c.Request.Header)

// 	if AuthByGeoip(ip, authConfig) || AllowInternalIp(ip) {
// 		ApiLogger.Info("Authorized: %s to %s/%v\n", ip, c.Request.Host, c.Request.URL.Path)
// 		c.JSON(http.StatusOK, struct{}{})
// 	} else {
// 		ApiLogger.Info("Denied: %s to %s/%v\n", ip, c.Request.Host, c.Request.URL.Path)
// 		c.JSON(http.StatusForbidden, struct{}{})
// 	}
// }

// func handleAuth(c *gin.Context) {
// 	real_ip := c.GetHeader(HeaderRealClientIP)
// 	user, pass, hasCredentials := c.Request.BasicAuth()

// 	if AuthByGeoip(real_ip, authConfig) && hasCredentials && AuthByCredentials(user, pass) && AuthByTempAllow("", "") {
// 		c.JSON(http.StatusOK, struct{}{})
// 	} else {
// 		c.JSON(http.StatusForbidden, struct{}{})
// 	}
// }
