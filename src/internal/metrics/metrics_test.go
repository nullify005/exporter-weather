package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

var requiredMetrics = []string{
	"weather_temperature_celcius 22.2",
	"weather_humidity 66",
}

const metricsPath string = "/metrics"

func TestMetrics(t *testing.T) {
	m := New()
	m.Humidity(66)
	m.Temp(22.2)
	request := httptest.NewRequest(http.MethodGet, metricsPath, nil)
	recorder := httptest.NewRecorder()
	promhttp.Handler().ServeHTTP(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusOK)
	body := recorder.Body.String()
	for _, s := range requiredMetrics {
		assert.Contains(t, body, s)
	}
}
