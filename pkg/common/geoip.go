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

type GeoipLookup struct {
	ip        string
	Country   string
	CityName  string
	Latitude  float64
	Longitude float64
}

func LookupIP(ip string) (GeoipLookup, error) {
	// If you are using strings that may be invalid, check that ip is not nil
	ip_addr := net.ParseIP(ip)
	record, err := featureGeoip.Database.City(ip_addr)
	if err != nil {
		log.Printf("Error parsing IP: %v\n", err)
		return GeoipLookup{ip: ip, Country: "--", CityName: "--", Latitude: 0, Longitude: 0}, err
	}
	ret := GeoipLookup{
		ip:        ip,
		Country:   record.Country.IsoCode,
		CityName:  record.City.Names["en"],
		Latitude:  record.Location.Latitude,
		Longitude: record.Location.Longitude,
	}

	GetServerConfig().RootLogger.Debugf("Geoip: lookup '%s' | %v - %v - %v/%v\n", ip, ret.Country, ret.CityName, ret.Latitude, ret.Longitude)

	return ret, nil

}
func LookupCountryCode(ip string) (string, error) {
	ret, err := LookupIP(ip)
	return ret.Country, err
}
