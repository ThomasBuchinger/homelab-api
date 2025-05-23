package api

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/thomasbuchinger/homelab-api/pkg/geoip"
)

var metricGeoipRequestsTotal *prometheus.CounterVec = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "homelabapi_geoip_requests_total",
	Help: "Total Requests seen by the ExtAuthz handler by country",
}, []string{"latitude", "longitude", "country", "city"})

func LogMetricGeoip(loc geoip.GeoipLookup) {
	metricGeoipRequestsTotal.WithLabelValues(fmt.Sprint(loc.Latitude), fmt.Sprint(loc.Longitude), loc.Country, loc.CityName).Inc()
}
