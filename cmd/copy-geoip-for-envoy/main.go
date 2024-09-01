package main

import (
	"fmt"
	"os"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func copyGeoipDatabase(sourceFile, destinationFile string) {
	input, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return
	}
}

func main() {
	copyGeoipDatabase(common.GetEnvWithDefault("GEOIP_DATABASE", "/geoip/GeoLite2-City.mmdb"), common.GetEnvWithDefault("COPY_GEOIP_DATABASE", "/geoip-envoy/GeoLite2-City.mmdb"))
	fmt.Println("Successfully Copied Database")
	os.Exit(0)
}
