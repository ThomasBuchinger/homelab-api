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

	code, _ := common.LookupCountryCode("91.115.30.180")
	if code != "AT" {
		t.Fatalf("Country code: Expected: AT, Actual: %s", code)
	}
}
func Test_FeatureGeoip_should_resolve_country_codes_in_ipv6(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()

	code, _ := common.LookupCountryCode("2001:871:22b:6c4c:a63b:11c4:2e23:e71a")
	if code != "AT" {
		t.Fatalf("Country code: Expected: AT, Actual: %s", code)
	}
}

func Test_FeatureGeoip_should_resolve_cities_and_location(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()

	loc, _ := common.LookupIP("91.115.30.180")
	if loc.Country != "AT" {
		t.Fatalf("Country code: Expected: AT, Actual: %s", loc.Country)
	}
	if loc.CityName != "Bad Fischau" {
		t.Fatalf("City Name: Expected: AT, Actual: %s", loc.CityName)
	}
	if loc.Latitude == 0 || loc.Longitude == 0 {
		t.Fatalf("Coordinates: Coordinates should not be %v/%v", loc.Latitude, loc.Longitude)
	}
}
