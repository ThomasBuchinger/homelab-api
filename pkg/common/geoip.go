package common

import (
	"log"
	"net"
	"path"

	geoip "github.com/oschwald/geoip2-golang"
)

type GeoipFeature struct {
	Enabled      bool
	Failed       bool
	FailedError  error
	DatapasePath string
	Database     *geoip.Reader
}

var featureGeoip GeoipFeature = GeoipFeature{
	Enabled:      false,
	Failed:       false,
	FailedError:  nil,
	DatapasePath: getGeoipDatabasePath(),
	Database:     nil,
}

func getGeoipDatabasePath() string {
	default_database_path := "geoip/GeoLite2-City.mmdb"
	databse_path := GetEnvWithDefault(EnvGeoipDatabase, default_database_path)
	if !path.IsAbs(databse_path) {
		databse_path = path.Clean("../../" + databse_path)
	}
	return databse_path
}

func FeatureGeoipInit() bool {
	if featureGeoip.Failed {
		return featureGeoip.Enabled
	}

	db, err := geoip.Open(getGeoipDatabasePath())
	if err != nil {
		log.Println(err)
		featureGeoip.Failed = true
		featureGeoip.FailedError = err
	} else {
		featureGeoip.Enabled = true
		featureGeoip.Database = db
	}
	return featureGeoip.Enabled
}

func FeatureGeoipClose() {
	featureGeoip.Database.Close()
	featureGeoip.Enabled = false
	featureGeoip.Failed = false
	featureGeoip.FailedError = nil
	featureGeoip.DatapasePath = ""
	featureGeoip.Database = nil
}

func LookupIP(ip string) (string, error) {
	// If you are using strings that may be invalid, check that ip is not nil
	ip_addr := net.ParseIP(ip)
	record, err := featureGeoip.Database.City(ip_addr)
	if err != nil {
		log.Printf("Error parsing IP: %v\n", err)
		return "--", err
	}
	log.Printf("Geoip: lookup '%s' | %v - %v - %v/%v\n", ip, record.Country.IsoCode, record.City.Names["en"], record.Location.Latitude, record.Location.Longitude)

	return record.Country.IsoCode, nil
}
