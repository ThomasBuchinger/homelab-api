package common_test

import (
	"testing"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func Test_FeatureGeoip_should_be_enabled_manualy(t *testing.T) {
	cfg := common.GetServerConfig()
	if cfg.EnableGeoip {
		t.Fatalf("Geoip not initialized. EnableGeoip should be false")
	}
}
func Test_FeatureGeoip_should_be_can_be_enabled(t *testing.T) {
	ok := common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()
	if !ok {
		t.Fatalf("Geoip database should beavailable in tests")
	}

	cfg := common.GetServerConfig()
	if !cfg.EnableGeoip {
		t.Fatalf("Geoip should be enabled if database is available")
	}
}
func Test_FeatureGeoip_should_resolve_country_codes(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()

	code, _ := common.LookupIP("91.115.30.180")
	if code != "AT" {
		t.Fatalf("Country code: Expected: AT, Actual: %s", code)
	}
}
func Test_FeatureGeoip_should_resolve_country_codes_in_ipv6(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()

	code, _ := common.LookupIP("2001:871:22b:6c4c:a63b:11c4:2e23:e71a")
	if code != "AT" {
		t.Fatalf("Country code: Expected: AT, Actual: %s", code)
	}
}
