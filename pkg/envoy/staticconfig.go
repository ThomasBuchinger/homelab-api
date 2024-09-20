package envoy

import (
	"time"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func AuthPolicyHomlapApiPublic() EnvoyVirtualHost {
	return NewVirtualHost("api.cloud.buc.sh", "external-homelab-api.external-homelabapi.svc:8080").
		WithExpireationIn(time.Duration(-1)).
		WithRoute("/$", StrictRoute().WithMatchType(MatchTypeRegex).WithVerbs("GET").WithCountries(common.GetConfigEnvoyAuthAllowedCountries()...)).
		WithRoute("/index.html$", StrictRoute().WithMatchType(MatchTypeRegex).WithVerbs("GET").WithCountries(common.GetConfigEnvoyAuthAllowedCountries()...)).
		WithRoute("/404.html$", StrictRoute().WithMatchType(MatchTypeRegex).WithVerbs("GET").WithCountries(common.GetConfigEnvoyAuthAllowedCountries()...)).
		WithRoute("/favicon.ico$", StrictRoute().WithMatchType(MatchTypeRegex).WithVerbs("GET").WithCountries(common.GetConfigEnvoyAuthAllowedCountries()...)).
		WithRoute("/_next", StrictRoute().WithVerbs("GET").WithCountries(common.GetConfigEnvoyAuthAllowedCountries()...)).
		WithRoute("/api/public", StrictRoute().WithVerbs("GET").WithCountries(common.GetConfigEnvoyAuthAllowedCountries()...)).
		WithRoute("/icons", StrictRoute().WithVerbs("GET").WithCountries(common.GetConfigEnvoyAuthAllowedCountries()...))
}

func AuthPolicyDefault() EnvoyVirtualHost {
	return NewVirtualHost("*.cloud.buc.sh", "invalid.svc:8080").
		WithExpireationIn(time.Duration(-1)).
		WithRoute("/", StrictRoute().WithVerbs("GET").WithCountries(common.GetConfigEnvoyAuthAllowedCountries()...))
}
