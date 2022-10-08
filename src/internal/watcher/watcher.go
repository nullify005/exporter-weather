package watcher

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nullify005/exporter-weather/internal/bom"
	"github.com/nullify005/exporter-weather/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Watcher struct {
	Interval time.Duration
	Geohash  string
	Listen   string
}

type Option func(w *Watcher)

func WithInterval(t time.Duration) Option {
	return func(w *Watcher) {
		w.Interval = t
	}
}

func WithListen(hostPort string) Option {
	return func(w *Watcher) {
		w.Listen = hostPort
	}
}

func New(geohash string, opts ...Option) Watcher {
	w := Watcher{
		Interval: 30 * time.Second,
		Geohash:  geohash,
		Listen:   "127.0.0.1:2112",
	}
	for _, opt := range opts {
		opt(&w)
	}
	return w
}

func (w *Watcher) Watch() {
	log.SetPrefix("exporter-weather: ")
	log.SetFlags(log.LstdFlags)
	log.Print("location: ", w.Geohash)
	log.Print("interval: ", w.Interval)
	log.Print("listen: ", w.Listen)
	watch(w)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", healthHandler)
	log.Fatal(http.ListenAndServe(w.Listen, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func watch(w *Watcher) {
	go func() {
		log.Print("setting up observation loop")
		m := metrics.New()
		for {
			time.Sleep(w.Interval)
			res, err := bom.Observe(w.Geohash)
			if err != nil {
				log.Println(err)
				m.Error(1)
				continue
			}
			log.Printf("received: %s", res.String())
			m.Temp(res.Temp)
			m.FeelsLike(res.TempFeelsLike)
			m.WindSpeed(float64(res.Wind.SpeedKilometre))
			m.WindGust(float64(res.Gust.SpeedKilometre))
			m.Humidity(float64(res.Humidity))
			m.Rain(float64(res.RainSince9Am))
			m.Error(0)
			m.Bearing(bom.WindBearing(res.Wind.Direction))
		}
	}()
}
