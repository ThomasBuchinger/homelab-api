package backend

import (
	"testing"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func Test_AuthConfig_should_CountryList_uses_UpperCase_only(t *testing.T) {
	conf := CreateAuthConfig("DE, At, ch", "", "")
	expected := []string{"DE", "AT", "CH"}

	for index, v := range conf.AllowedCountries {
		if expected[index] != v {
			t.Fatalf("Expected Countries DE and AT ti be in the list")
		}
	}
}
func Test_authGeoip_should_authenticate_AT_ips(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()

	if !AuthByGeoip("91.115.30.180", CreateAuthConfig("DE, AT", "", "")) {
		t.Fatalf("Should allow IPs from AT")
	}
}
func Test_authGeoip_should_forbid_empty_string(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()

	if AuthByGeoip("", CreateAuthConfig("DE, AT", "", "")) {
		t.Fatalf("Empty String is not allowed")
	}
}
func Test_authGeoip_should_forbid_random_ips(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()

	if AuthByGeoip("1.2.3.4", CreateAuthConfig("DE, AT", "", "")) {
		t.Fatalf("Random IPs are not allowed")
	}
}
func Test_authGeoip_should_forbid_non_ips(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()

	if AuthByGeoip("host.example.com", CreateAuthConfig("DE, AT", "", "")) {
		t.Fatalf("Non IPs are not allowed")
	}
}
