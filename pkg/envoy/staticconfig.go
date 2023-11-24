package envoy

import "time"

func AuthPolicyHomlapApiPublic(countries []string) EnvoyVirtualHost {
	return NewVirtualHost("api.cloud.buc.sh", "external-homelab-api.external-homelabapi.svc:8080").
		WithExpireationIn(time.Duration(-1)).
		WithRoute("/$", DenyRoute().WithMatchType(MatchTypeRegex).WithVerbs("GET").WithCountries(countries...)).
		WithRoute("/index.html$", DenyRoute().WithMatchType(MatchTypeRegex).WithVerbs("GET").WithCountries(countries...)).
		WithRoute("/404.html$", DenyRoute().WithMatchType(MatchTypeRegex).WithVerbs("GET").WithCountries(countries...)).
		WithRoute("/favicon.ico$", DenyRoute().WithMatchType(MatchTypeRegex).WithVerbs("GET").WithCountries(countries...)).
		WithRoute("/_next", DenyRoute().WithVerbs("GET").WithCountries(countries...)).
		WithRoute("/api/public", DenyRoute().WithVerbs("GET").WithCountries(countries...)).
		WithRoute("/icons", DenyRoute().WithVerbs("GET").WithCountries(countries...))
}

func AuthPolicyMinioApi(countries []string) EnvoyVirtualHost {
	return NewVirtualHost("min.io", "min.io:9000").
		WithEndpoint("hidden_api").
		WithExpireationIn(time.Duration(-1)).
		WithRoute("/", DenyRoute().WithVerbs("*").WithCountries(countries...))
}

func AuthPolicyDefault(countries []string) EnvoyVirtualHost {
	return NewVirtualHost("*.cloud.buc.sh", "invalid.svc:8080").
		WithExpireationIn(time.Duration(-1)).
		WithRoute("/", DenyRoute().WithVerbs("GET").WithCountries(countries...))
}
