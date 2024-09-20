package envoy

import (
	"strings"
	"time"
)

type EnvoyVirtualHost struct {
	Domain    string
	Endpoint  string
	Permanent bool
	ExpiresAt time.Time

	Routing          map[string]RoutingRule
	AllowedCountries []string

	UpstreamUrl string
}
type RoutingRule struct {
	Type             MatchType
	AllowedVerbs     []string
	AllowedUsers     []string
	AllowedCountries []string
	RequiredHeaders  map[string]string
}
type MatchType string

const (
	MatchTypePrefix MatchType = "Prefix"
	MatchTypeRegex  MatchType = "Regex"
)

// ===== Builder Functions ========================================================================
func NewVirtualHost(from, to string) EnvoyVirtualHost {
	return EnvoyVirtualHost{
		Domain:    from,
		Endpoint:  "https",
		Permanent: false,
		ExpiresAt: time.Time{},

		Routing: map[string]RoutingRule{},

		UpstreamUrl: to,
	}
}
func (e EnvoyVirtualHost) WithExpireationIn(dur time.Duration) EnvoyVirtualHost {
	if dur < 0 {
		e.Permanent = true

	}
	e.ExpiresAt = time.Now().Add(dur)
	return e
}
func (e EnvoyVirtualHost) WithEndpoint(endpontName string) EnvoyVirtualHost {
	e.Endpoint = endpontName
	return e
}
func (e EnvoyVirtualHost) WithRoute(path string, rule RoutingRule) EnvoyVirtualHost {
	e.Routing[path] = rule
	return e
}
func (e EnvoyVirtualHost) ChangeRoutes(mapperFunc func(r RoutingRule) RoutingRule) EnvoyVirtualHost {
	for k, v := range e.Routing {
		e.Routing[k] = mapperFunc(v)
	}
	return e
}

func StrictRoute() RoutingRule {
	return RoutingRule{
		Type:             MatchTypePrefix,
		AllowedVerbs:     []string{},
		AllowedUsers:     []string{},
		AllowedCountries: []string{},
		RequiredHeaders:  map[string]string{},
	}
}
func AllowEverythngRoute() RoutingRule {
	return RoutingRule{
		Type:             MatchTypePrefix,
		AllowedVerbs:     []string{"*"},
		AllowedUsers:     []string{"*"},
		AllowedCountries: []string{"*"},
		RequiredHeaders:  map[string]string{},
	}
}

func (r RoutingRule) WithCountries(allowlist ...string) RoutingRule {
	r.AllowedCountries = listToUpperCase(allowlist)
	return r
}
func (r RoutingRule) WithUsers(allowlist ...string) RoutingRule {
	r.AllowedUsers = allowlist
	return r
}
func (r RoutingRule) WithVerbs(allowlist ...string) RoutingRule {
	r.AllowedVerbs = listToUpperCase(allowlist)
	return r
}
func (r RoutingRule) WithHeader(key, value string) RoutingRule {
	r.RequiredHeaders[key] = value
	return r
}

func (r RoutingRule) WithMatchType(matchype MatchType) RoutingRule {
	r.Type = matchype
	return r
}

func listToUpperCase(list []string) []string {
	upper := make([]string, len(list))
	for i, element := range list {
		upper[i] = strings.ToUpper(element)
	}
	return upper
}

// ===== Methods ==================================================================================
func (e EnvoyVirtualHost) IsExpired() bool {
	return time.Now().After(e.ExpiresAt) && !e.Permanent
}

// func (e EnvoyVirtualHost) IsCountryAllowed() bool {
// 	ret := false
// 	for _, v := range e.AllowedCountries()
// }
