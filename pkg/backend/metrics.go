package backend

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

var MetricGeoipRequestsTotal *prometheus.CounterVec = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "homelabapi_geoip_requests_total",
	Help: "Total Requests seen by the ExtAuthz handler by country",
}, []string{"Latitude", "Longitude", "Country", "City"})

func LogGeoipRequest(loc common.GeoipLookup) {
	MetricGeoipRequestsTotal.WithLabelValues(fmt.Sprint(loc.Latitude), fmt.Sprint(loc.Longitude), loc.Country, loc.CityName)
}
