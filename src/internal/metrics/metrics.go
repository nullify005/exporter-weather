package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	m *metrics

	mTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "weather_temperature_celcius",
		Help: "Current temperature.",
	})
	mFeelsLike = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "weather_temperature_feelslike_celcius",
		Help: "Current feels like temperature.",
	})
	mWindSpeed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "weather_windspeed_kph",
		Help: "Current windspeed.",
	})
	mWindGust = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "weather_windspeed_gust_kph",
		Help: "Current windspeed gusts.",
	})
	mHumidity = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "weather_humidity",
		Help: "Current humidity.",
	})
	mRain = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "weather_rain_since_9am_mm",
		Help: "Rain since 9AM.",
	})
	mErrorState = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "weather_error",
		Help: "Error flag indicating an weather error.",
	})
	mBearing = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "weather_wind_bearing",
		Help: "Current wind direction.",
	})
)

const defaultValue float64 = 20

type metrics struct {
	mu sync.Mutex
}

// NOTE: the http handler still needs to be registered elsewhere
// TODO: write some tests for this
func New() *metrics {
	if m == nil {
		m = &metrics{}
		prometheus.MustRegister(mTemp)
		prometheus.MustRegister(mFeelsLike)
		prometheus.MustRegister(mWindSpeed)
		prometheus.MustRegister(mWindGust)
		prometheus.MustRegister(mHumidity)
		prometheus.MustRegister(mRain)
		prometheus.MustRegister(mErrorState)
		prometheus.MustRegister(mBearing)
		m.Temp(defaultValue)
		m.FeelsLike(defaultValue)
		m.WindSpeed(defaultValue)
		m.WindGust(defaultValue)
		m.Humidity(defaultValue)
		m.Rain(defaultValue)
		m.Bearing(0)
		m.Error(1)
	}
	return m
}

func (m *metrics) Temp(t float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mTemp.Set(t)
}

func (m *metrics) FeelsLike(t float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mFeelsLike.Set(t)
}

func (m *metrics) WindSpeed(t float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mWindSpeed.Set(t)
}

func (m *metrics) WindGust(t float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mWindGust.Set(t)
}

func (m *metrics) Humidity(t float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mHumidity.Set(t)
}

func (m *metrics) Rain(t float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mRain.Set(t)
}

func (m *metrics) Error(t float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mErrorState.Set(t)
}

func (m *metrics) Bearing(t float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mBearing.Set(t)
}
