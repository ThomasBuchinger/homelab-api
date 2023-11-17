package common_test

import (
	"testing"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func Test_IsIPAddressInternal_should_be_true_for_localhost(t *testing.T) {
	if !common.IsIpAddressInternal("127.0.0.1") {
		t.Fatalf("127.0.0.1 shoudl be considered internal")
	}
}
func Test_IsIPAddressInternal_should_be_true_for_ipv6_localhost(t *testing.T) {
	if !common.IsIpAddressInternal("::ffff:127.0.0.1") {
		t.Fatalf("::ffff:127.0.0.1 shoudl be considered internal")
	}
}
