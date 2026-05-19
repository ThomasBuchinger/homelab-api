package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"go.uber.org/zap"
)

var ApiLogger *zap.SugaredLogger = common.GetServerConfig().RootLogger.Named("API")

func SetupDefaultRouter() *gin.Engine {
	serverConfig := common.GetServerConfig()
	logger := serverConfig.RootLogger.Desugar().Named("access")

	gin.SetMode(serverConfig.GinMode)
	router := gin.New()
	router.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{TimeFormat: time.RFC3339, UTC: true,
		SkipPaths: []string{
			"/api/livez",
			"/api/readyz",
		},
	}))
	router.Use(ginzap.RecoveryWithZap(logger, false))
	router.Use(requestid.New())

	router = setupCommonApiEndpoints(router)
	// router = SetupFrontendApiEndpoints(router)
	return router
}

func setupCommonApiEndpoints(r *gin.Engine) *gin.Engine {

	r.GET("/api/livez", handlePing)
	r.GET("/api/readyz", handlePing)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Publicly accessible API endpoints
	r.GET("/api/public/ping", handlePing)
	r.GET("/api/zib2", scrape_zib2)

	return r
}

func SetupStaticFileServing(r *gin.Engine) *gin.Engine {
	r.Use(static.Serve("/", static.LocalFile("./ui/out", true)))
	r.Use(static.Serve("/geoip", static.LocalFile("/geoip", true)))
	return r
}

func SetupDummyApiEndpointsForRunningInPublicMode(r *gin.Engine) *gin.Engine {
	r.GET("/", handlePing)
	return r
}

func handlePing(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

type ZibEpisode struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

func scrape_zib2(c *gin.Context) {
	targetURL := "https://on.orf.at/sendereihe/1211/zib-2"
	episodes := []ZibEpisode{}

	resp, err := http.Get(targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch page", "message": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "HTTP error: " + resp.Status})
		return
	}

	html, _ := goquery.NewDocumentFromReader(resp.Body)
	html.Find(".video-list a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		title := s.Find("h2.title").Text()

		// ORF returns relative URLs (/video/...). Convert them to absolute.
		if !strings.HasPrefix(href, "http") {
			href = "https://on.orf.at" + href
		}

		episodes = append(episodes, ZibEpisode{
			Url:   href,
			Title: strings.TrimSpace(title),
		})
	})

	c.JSON(http.StatusOK, gin.H{
		"episodes": episodes,
	})
}

// Authorized group (uses gin.BasicAuth() middleware)
// Same than:
// authorized := r.Group("/")
// authorized.Use(gin.BasicAuth(gin.Credentials{
//	  "foo":  "bar",
//	  "manu": "123",
//}))
// authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
// 	"foo":  "bar", // user:foo password:bar
// 	"manu": "123", // user:manu password:123
// }))

// authorized.POST("admin", func(c *gin.Context) {
// 	user := c.MustGet(gin.AuthUserKey).(string)
// 	// Parse JSON
// 	var json struct {
// 		Value string `json:"value" binding:"required"`
// 	}

// 	if c.Bind(&json) == nil {
// 		db[user] = json.Value
// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
// 	}
// })
