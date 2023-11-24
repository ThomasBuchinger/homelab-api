package backend

import (
	"github.com/gin-gonic/gin"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"github.com/thomasbuchinger/homelab-api/pkg/health"
)

func handleClientConfig(c *gin.Context) {
	real_ip := c.ClientIP()
	ip := c.GetHeader("x-forwarded-for")
	if ip == "" {
		ip = real_ip
	}
	country := "--"
	if common.GetServerConfig().EnableGeoip {
		country, _ = common.LookupCountryCode(real_ip)
	}

	ApiLogger.Debugw("Client-Config", "ip", ip, "country", country, "internal", common.IsIpAddressInternal(ip))
	c.JSON(200, gin.H{
		"ip":       ip,
		"real_ip":  real_ip,
		"internal": common.IsIpAddressInternal(ip),
		"country":  country,
	})
}

func handlePublicHealth(c *gin.Context) {
	target := c.Query("target")
	targets := map[string]func() health.ExternalHealthCheckResult{
		"Servers":      health.Ok,
		"Network":      health.Ok,
		"API":          health.Ok,
		"External API": health.CheckApiPublic,
	}
	res := targets[target]()
	messages := []string{}
	for _, r := range res.Results {
		if r.Message != "" {
			messages = append(messages, r.Message)
		}
	}
	ApiLogger.Debugf("health: %v | Messages: %v", res.Health, messages)
	c.JSON(200, gin.H{
		"healthy":  res.Health,
		"passed":   res.PassedChecks,
		"total":    res.TotalChecks,
		"messages": messages,
	})
}
