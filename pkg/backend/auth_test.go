package backend

import (
	"testing"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func Test_authGeoip_should_authenticate_AT_ips(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()

	if !AuthByGeoip("91.115.30.180") {
		t.Fatalf("Should allow IPs from AT")
	}
}
func Test_authGeoip_should_forbid_empty_string(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()

	if AuthByGeoip("") {
		t.Fatalf("Empty String is not allowed")
	}
}
func Test_authGeoip_should_forbid_random_ips(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()

	if AuthByGeoip("1.2.3.4") {
		t.Fatalf("Random IPs are not allowed")
	}
}
func Test_authGeoip_should_forbid_non_ips(t *testing.T) {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()

	if AuthByGeoip("host.example.com") {
		t.Fatalf("Non IPs are not allowed")
	}
}
