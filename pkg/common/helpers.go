package common

import (
	"strings"
)

func IsIpAddressInternal(ip string) bool {
	if ip == "127.0.0.1" || ip == "::ffff:127.0.0.1" || ip == "::1" {
		return true
	}
	if strings.HasPrefix(ip, "10.0.0") && ip != "10.0.0.21" {
		return true
	}
	return false
}
