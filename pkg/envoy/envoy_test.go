package envoy_test

import (
	"testing"
	"time"

	"github.com/thomasbuchinger/homelab-api/pkg/envoy"
)

func getDefaultVirtualHost() envoy.EnvoyVirtualHost {
	return envoy.NewVirtualHost("host.public.com", "host.internal.com")
}

func defautDuration() time.Duration {
	d, _ := time.ParseDuration("1h")
	return d
}

func Test_envoyhost_should_have_no_routes_by_default(t *testing.T) {
	e := getDefaultVirtualHost()
	if len(e.Routing) != 0 {
		t.Fatalf("Expected Host to have no routes configured!")
	}
}

func Test_envoyhost_should_default_config_is_always_expired(t *testing.T) {
	e := getDefaultVirtualHost()
	if !e.IsExpired() {
		t.Fatalf("Expected host be expired")
	}
}

func Test_envoyhost_should_be_able_to_configure_expireation(t *testing.T) {
	e := getDefaultVirtualHost().WithExpireationIn(defautDuration())
	if e.IsExpired() {
		t.Fatalf("Host should not be expired")
	}
}
func Test_envoyhost_negative_expiration_means_permanent(t *testing.T) {
	e := getDefaultVirtualHost().WithExpireationIn(time.Duration(-1))
	if e.IsExpired() {
		t.Fatalf("Host should not be expired")
	}
	if !e.Permanent {
		t.Fatalf("Host should not be permanent")
	}
}
func Test_envoyhost_be_able_to_chain_methods(t *testing.T) {
	e := getDefaultVirtualHost().
		WithExpireationIn(time.Duration(defautDuration())).
		WithEndpoint("test")

	if e.IsExpired() {
		t.Fatalf("Host should not be expired")
	}
	if e.Endpoint != "test" {
		t.Fatalf("Expected endpoint to be 'test'. Got: %v", e.Endpoint)
	}
}

func Test_envoyhost_be_able_to_define_routes(t *testing.T) {
	e := getDefaultVirtualHost().
		WithRoute("/", envoy.DenyRoute()).
		WithRoute("/api", envoy.DenyRoute())

	if len(e.Routing) != 2 {
		t.Fatalf("Expected: 2 Routes  Got: %v", len(e.Routing))
	}
}
func Test_envoyhost_can_bulk_update_routes(t *testing.T) {
	allowed := []string{"AT", "DE"}
	mapperFunc := func(r envoy.RoutingRule) envoy.RoutingRule {
		return r.WithCountries(allowed...)
	}

	e := getDefaultVirtualHost().
		WithRoute("/", envoy.DenyRoute()).
		WithRoute("/v1", envoy.DenyRoute()).
		ChangeRoutes(mapperFunc)

	countries := e.Routing["/"].AllowedCountries
	if len(countries) != 2 {
		t.Fatalf("Expected: 2 Allowed Country on '/'.  Got: %v", countries)
	}
	countries = e.Routing["/v1"].AllowedCountries
	if len(countries) != 2 {
		t.Fatalf("Expected: 2 Allowed Country on '/v1'.  Got: %v", countries)
	}
}

// ===== Routes ===================================================================================
func Test_routes_can_be_georestricted_to_one_or_more_countries(t *testing.T) {
	e := getDefaultVirtualHost().
		WithRoute("/", envoy.DenyRoute().WithCountries("AT", "DE"))
	countries := e.Routing["/"].AllowedCountries
	if len(countries) != 2 {
		t.Fatalf("Expected: 2 Allowed Country  Got: %v", countries)
	}
}
func Test_routes_should_allow_only_certain_users(t *testing.T) {
	e := getDefaultVirtualHost().
		WithRoute("/", envoy.DenyRoute().WithUsers("someone", "anotherone"))
	users := e.Routing["/"].AllowedUsers
	if len(users) != 2 {
		t.Fatalf("Expected: 2 Allowed Users  Got: %v", users)
	}
}
func Test_routes_should_allow_only_certain_http_verbs(t *testing.T) {
	e := getDefaultVirtualHost().
		WithRoute("/", envoy.DenyRoute().WithVerbs("GET"))
	verbs := e.Routing["/"].AllowedVerbs
	if len(verbs) != 1 {
		t.Fatalf("Expected: 1 Allowed Verb  Got: %v", verbs)
	}
}
func Test_routes_can_enforce_certain_headers(t *testing.T) {
	e := getDefaultVirtualHost().
		WithRoute("/", envoy.DenyRoute().WithVerbs("GET"))
	verbs := e.Routing["/"].AllowedVerbs
	if len(verbs) != 1 {
		t.Fatalf("Expected: 1 Allowed Verb  Got: %v", verbs)
	}
}
func Test_routes_use_regex_matchers(t *testing.T) {
	e := getDefaultVirtualHost().
		WithRoute("/", envoy.DenyRoute().WithMatchType(envoy.MatchTypeRegex))
	matchtype := e.Routing["/"].Type
	if matchtype != envoy.MatchTypeRegex {
		t.Fatalf("Expected: 'Regex' Matchtype  Got: %v", matchtype)
	}
}
