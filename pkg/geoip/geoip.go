package geoip

import (
	"log"
	"net"
	"os"
	"path"

	geoip "github.com/oschwald/geoip2-golang"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"go.uber.org/zap"
)

var database *geoip.Reader
var geoipLogger *zap.SugaredLogger

type GeoipLookup struct {
	ip        string
	Country   string
	CityName  string
	Latitude  float64
	Longitude float64
}

func LookupIP(ip string) (GeoipLookup, error) {
	clodeFunc, err := openGeoipDatabase()
	defer clodeFunc()
	if err != nil {
		geoipLogger.Error("Cannot open GeoIP Database: ")
		return GeoipLookup{ip: ip, Country: "--", CityName: "--", Latitude: 0, Longitude: 0}, err
	}

	// If you are using strings that may be invalid, check that ip is not nil
	ip_addr := net.ParseIP(ip)
	record, err := database.City(ip_addr)
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

	geoipLogger.Debugf("Geoip: lookup '%s' | %v - %v - %v/%v\n", ip, ret.Country, ret.CityName, ret.Latitude, ret.Longitude)

	return ret, nil

}
func LookupCountryCode(ip string) (string, error) {
	ret, err := LookupIP(ip)
	return ret.Country, err
}


func openGeoipDatabase() (func() error, error) {
	geoipLogger = common.GetServerConfig().RootLogger.Named("GeoIP")

	database_path := getGeoipDatabasePath()
	var err error
	database, err = geoip.Open(database_path)
	if err != nil {
		return nil, err
	}

	return database.Close, nil
}

func FeatureGeoipEnabled() bool {
	_, err := os.Stat(getGeoipDatabasePath())
	return err == nil
}

func getGeoipDatabasePath() string {
	database_path := "geoip/GeoLite2-City.mmdb"
	if !path.IsAbs(database_path) {
		database_path = path.Clean("../../" + database_path)
	}
	return database_path
}
