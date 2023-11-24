package backend_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/thomasbuchinger/homelab-api/pkg/backend"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"github.com/thomasbuchinger/homelab-api/pkg/envoy"
)

func GetDummyAuthPolicies() []envoy.EnvoyVirtualHost {
	return []envoy.EnvoyVirtualHost{
		envoy.NewVirtualHost("*.company.com", "wildcard"),
		envoy.NewVirtualHost("dummy1.example.com", "dummy1"),
	}
}
func GetDummyPolicyMatcher(rule envoy.RoutingRule) (backend.PolicyMatcher, envoy.EnvoyVirtualHost) {
	vhost := envoy.NewVirtualHost("test.example.com", "somewhere.svc").WithRoute("/", rule)
	matcher := backend.PolicyMatcher{Passed: true}.MatchPath("/", vhost)
	return matcher, vhost
}

func Test_ResolveAuthPolicyForRequest_finds_a_exact_match(t *testing.T) {
	e := backend.ResolveAuthPolicyForRequest("dummy1.example.com", GetDummyAuthPolicies())
	if e.UpstreamUrl != "dummy1" {
		t.Fatalf("Expected Host to be dummy1. Got: %s", e.UpstreamUrl)
	}
}
func Test_ResolveAuthPolicyForRequest_finds_wilcards_matches(t *testing.T) {
	e := backend.ResolveAuthPolicyForRequest("something.company.com", GetDummyAuthPolicies())
	if e.UpstreamUrl != "wildcard" {
		t.Fatalf("Expected Host to be wildcard. Got: %s", e.UpstreamUrl)
	}
}
func Test_ResolveAuthPolicyForRequest_works_with_catch_all(t *testing.T) {
	e := backend.ResolveAuthPolicyForRequest("hello.world", []envoy.EnvoyVirtualHost{envoy.NewVirtualHost("*", "everything.svc")})
	if e.UpstreamUrl != "everything.svc" {
		t.Fatalf("Expected Host to be everything.svc. Got: %s", e.UpstreamUrl)
	}
}
func Test_ResolveAuthPolicyForRequest_uses_the_longest_wildcard_match(t *testing.T) {
	e := backend.ResolveAuthPolicyForRequest("somewhere.example.com", []envoy.EnvoyVirtualHost{
		envoy.NewVirtualHost("*.example.com", "example"),
		envoy.NewVirtualHost("*.com", "com"),
	})
	if e.UpstreamUrl != "example" {
		t.Fatalf("Expected Host to be example. Got: %s", e.UpstreamUrl)
	}
}

func Test_ResolveAuthPolicyForRequest_works_with_port_numbers_in_hostname(t *testing.T) {
	e := backend.ResolveAuthPolicyForRequest("dummy1.example.com:80", GetDummyAuthPolicies())
	if e.UpstreamUrl != "dummy1" {
		t.Fatalf("Expected Host to be dummy1. Got: %s", e.UpstreamUrl)
	}
}
func Test_SplitCountryConfigString(t *testing.T) {
	countries := backend.SplitCountryConfigString("DE, At, ch")
	expected := []string{"DE", "AT", "CH"}

	for index, v := range countries {
		if expected[index] != v {
			t.Fatalf("Expected Countries DE and AT ti be in the list")
		}
	}
}

func Test_EvaluateAuthPolicyAgainstRequest_handles_the_simplest_request(t *testing.T) {
	e := envoy.NewVirtualHost("test.example.com", "somewhere.svc").WithRoute("/", envoy.AllowEverythngRoute())
	req, _ := http.NewRequest("GET", "https://test.eaxmple.com/hello", http.NoBody)

	if !backend.EvaluateAuthPolicyAgainstRequest(*req, e).Passed {
		t.Fatal("Expected Route /hello to me matched by /")
	}
}

func Test_EvaluateAuthPolicyAgainstRequest_handles_missing_path_in_url(t *testing.T) {
	e := envoy.NewVirtualHost("test.example.com", "somewhere.svc").WithRoute("/", envoy.AllowEverythngRoute())
	req, _ := http.NewRequest("GET", "https://test.eaxmple.com", http.NoBody)

	if !backend.EvaluateAuthPolicyAgainstRequest(*req, e).Passed {
		t.Fatal("Expected empty Route to me matched by /")
	}
}
func Test_EvaluateAuthPolicyAgainstRequest_handles_url_with_auth_prefix(t *testing.T) {
	e := envoy.NewVirtualHost("test.example.com", "somewhere.svc").WithRoute("/hello", envoy.AllowEverythngRoute())
	req, _ := http.NewRequest("GET", "https://test.eaxmple.com/auth/hello", http.NoBody)

	if !backend.EvaluateAuthPolicyAgainstRequest(*req, e).Passed {
		t.Fatal("Expected Route /auth/hello to me matched by /hello")
	}
}
func Test_EvaluateAuthPolicyAgainstRequest_no_routes_get_denied(t *testing.T) {
	e := envoy.NewVirtualHost("test.example.com", "somewhere.svc")
	req, _ := http.NewRequest("GET", "https://test.eaxmple.com/auth/hello", http.NoBody)

	if backend.EvaluateAuthPolicyAgainstRequest(*req, e).Passed {
		t.Fatal("Expected vHost without routes to deny everything")
	}
}
func Test_EvaluateAuthPolicyAgainstRequest_different_routes_get_denied(t *testing.T) {
	e := envoy.NewVirtualHost("test.example.com", "somewhere.svc").WithRoute("/hello", envoy.AllowEverythngRoute())
	req, _ := http.NewRequest("GET", "https://test.eaxmple.com/auth/world", http.NoBody)

	if backend.EvaluateAuthPolicyAgainstRequest(*req, e).Passed {
		t.Fatal("Expected vHost without matching routes to deny everything")
	}
}
func Test_AuthByGeoip_should_approve_AT_ips(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()
	matcher, vhost := GetDummyPolicyMatcher(envoy.AllowEverythngRoute().WithCountries("at"))

	res := matcher.AuthByGeoip("91.115.30.180", vhost)
	if !res.Passed {
		t.Fatalf("Should allow IPs from AT: " + strings.Join(res.Messages, " "))
	}
}
func Test_AuthByGeoip_should_deny_empty_string(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()
	matcher, vhost := GetDummyPolicyMatcher(envoy.AllowEverythngRoute().WithCountries("at"))

	res := matcher.AuthByGeoip("", vhost)
	if res.Passed {
		t.Fatalf("Should deny empty string")
	}
}

func Test_AuthByGeoip_should_deny_random_ip(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()
	matcher, vhost := GetDummyPolicyMatcher(envoy.AllowEverythngRoute().WithCountries("at"))

	res := matcher.AuthByGeoip("1.2.3.4", vhost)
	if res.Passed {
		t.Fatalf("IP 1.2.3.4 is not allowed")
	}
}
func Test_AuthByGeoip_should_deny_hostnames(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()
	matcher, vhost := GetDummyPolicyMatcher(envoy.AllowEverythngRoute().WithCountries("at"))

	res := matcher.AuthByGeoip("host.example.com", vhost)
	if res.Passed {
		t.Fatalf("'host.example.com' it not an ip address")
	}
}

// func Test_AuthBasic_user(t *testing.T) {
// 	matcher, vhost := GetDummyPolicyMatcher(envoy.AllowEverythngRoute().WithCountries("at"))

// 	res := matcher.AuthByCredentials("admin", "admin", true, vhost)
// 	if res.Passed {
// 		t.Fatalf("'host.example.com' it not an ip address")
// 	}
// }
