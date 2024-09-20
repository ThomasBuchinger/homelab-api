package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thomasbuchinger/homelab-api/pkg/envoy"
	"github.com/thomasbuchinger/homelab-api/pkg/geoip"
)

func SetupAuthApiEndpoints(r *gin.Engine) *gin.Engine {
	r.Any("/auth/*authpath", handleAuth)
	return r
}

const (
	HeaderSourceIP = "X-Envoy-Forwarded-For"
)

func allKnownVirtualHosts() []envoy.EnvoyVirtualHost {
	return []envoy.EnvoyVirtualHost{
		envoy.AuthPolicyHomlapApiPublic(),
		envoy.AuthPolicyDefault(),
	}
}

func handleAuth(c *gin.Context) {
	ip := c.GetHeader(HeaderSourceIP)
	authPolicy := ResolveAuthPolicyForRequest(c.Request.Host, allKnownVirtualHosts())

	loc, err := geoip.LookupIP(ip)
	if err == nil {
		LogMetricGeoip(loc)
	}

	ApiLogger.Debugf("Authentication Request: Remote: %s | Host: %v | Path: %v | Headers: %v", c.RemoteIP(), c.Request.Host, c.Request.URL.Path, c.Request.Header)

	if EvaluateAuthPolicyAgainstRequest(*c.Request, authPolicy).Passed {
		ApiLogger.Info("Authorized: %s to %s/%v\n", ip, c.Request.Host, c.Request.URL.Path)
		c.JSON(http.StatusOK, struct{}{})
	} else {
		ApiLogger.Info("Denied: %s to %s/%v\n", ip, c.Request.Host, c.Request.URL.Path)
		c.JSON(http.StatusForbidden, struct{}{})
	}
}

func matchWildcardToHostname(pattern, hostname string) (bool, int) {
	if strings.HasPrefix(pattern, "*") {
		pattern = strings.TrimLeft(pattern, "*")
		return strings.HasSuffix(hostname, pattern), len(pattern)
	}
	return false, 0
}

func normalizeRoutePath(path string) string {
	path = strings.TrimPrefix(path, "/auth")
	if path == "" {
		path = "/"
	}
	return path
}

func ResolveAuthPolicyForRequest(hostname string, list []envoy.EnvoyVirtualHost) envoy.EnvoyVirtualHost {
	best_match := envoy.AuthPolicyDefault()
	match_prio := -1
	hostname = strings.Split(hostname, ":")[0]

	for _, virtualhost_config := range list {
		if hostname == virtualhost_config.Domain && match_prio < 99 {
			best_match = virtualhost_config
		}
		match_wildcard, match_len := matchWildcardToHostname(virtualhost_config.Domain, hostname)
		if match_wildcard && match_len > match_prio {
			best_match = virtualhost_config
			match_prio = match_len
		}
	}

	return best_match
}

type PolicyMatcher struct {
	Passed   bool
	Messages []string

	Route string
	Rule  envoy.RoutingRule
}

func (p PolicyMatcher) Ok() PolicyMatcher {
	p.Passed = p.Passed && true
	return p
}

func (p PolicyMatcher) Fail(message string) PolicyMatcher {
	p.Messages = append(p.Messages, message)
	p.Passed = p.Passed && false
	return p
}

func (p PolicyMatcher) MatchPath(path string, e envoy.EnvoyVirtualHost) PolicyMatcher {
	matchFunc := func(path string, route string, matchtype envoy.MatchType) bool {
		if matchtype == envoy.MatchTypePrefix {
			return strings.HasPrefix(path, route)
		}
		if matchtype == envoy.MatchTypeRegex {
			matched, _ := regexp.MatchString("^"+route, path)
			return matched
		}
		return false
	}

	for route, ruleset := range e.Routing {
		if matchFunc(path, route, ruleset.Type) {
			p.Route = route
			return p.Ok()
		}
	}
	return p.Fail("No matching Path")
}

func (p PolicyMatcher) MatchIpToAllowedContries(ip string, vhost envoy.EnvoyVirtualHost) PolicyMatcher {
	allowed := vhost.Routing[p.Route].AllowedCountries
	if len(allowed) == 0 {
		return p.Fail("Geoblocked every Country")
	}
	if len(allowed) == 1 && allowed[0] == "*" {
		return p.Ok()
	}

	if !geoip.FeatureGeoipEnabled() {
		p.Fail("Geopip Feature required, but not Enabled! Did you configure the Database?")
	}

	result, err := geoip.LookupIP(ip)
	if err != nil {
		return p.Fail("Geoip Lookup failed: " + err.Error())
	}

	for _, c := range allowed {
		if c == result.Country {
			return p.Ok()
		}
	}

	return p.Fail(fmt.Sprintf("Country '%s' is not allowed. Must be: %s", result.Country, strings.Join(allowed, ", ")))
}

func (p PolicyMatcher) AuthByCredentials(user string, pass string, ok bool, vhost envoy.EnvoyVirtualHost) PolicyMatcher {

	return p.Fail("adasd")
}

func EvaluateAuthPolicyAgainstRequest(req http.Request, vhost envoy.EnvoyVirtualHost) PolicyMatcher {
	matcher := PolicyMatcher{Passed: true}
	real_path := normalizeRoutePath(req.URL.Path)
	real_ip := req.Header.Get(HeaderSourceIP)
	// user, pass, ok := req.BasicAuth()

	matcher = matcher.MatchPath(real_path, vhost).
		MatchIpToAllowedContries(real_ip, vhost)

	return matcher
}

// func AllowInternalIp(ip string) bool {
// 	return common.IsIpAddressInternal(ip)
// }

// func AuthByCredentials(user string, password string) bool {
// 	auth_user := common.GetEnvWithDefault("AUTH_USER", "")
// 	auth_pass := common.GetEnvWithDefault("AUTH_PASS", "")

// 	return (auth_user != "" && auth_user == user) && (auth_pass != "" && auth_pass == password)
// }

// func AuthByTempAllow(host string, path string) bool {
// 	return host == path && false
// }
