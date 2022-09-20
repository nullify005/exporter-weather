package cli

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nullify005/exporter-weather/internal/bom/location"
	"github.com/nullify005/exporter-weather/internal/bom/observation"
	wind "github.com/nullify005/exporter-weather/internal/data"
	"github.com/nullify005/exporter-weather/internal/render"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var flagLocation = flag.String("location", "", "The geohash for the observation location (use lookup to find it)")
var flagInterval = flag.Int("interval", 30, "The observation Interval in Seconds")
var flagListenPort = flag.Int("port", 2112, "The HTTP port to listen on for metrics & health")
var flagHelp = flag.Bool("help", false, "Command line arguments")
var flagLookup = flag.String("lookup", "", "Lookup the geohash for a given location")

var (
	metricTemp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "observation_temperature_celcius",
		Help: "Current temperature.",
	})
	metricTempFeelsLike = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "observation_temperature_feelslike_celcius",
		Help: "Current feels like temperature.",
	})
	metricWindSpeed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "observation_windspeed_kph",
		Help: "Current windspeed.",
	})
	metricWindSpeedGust = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "observation_windspeed_gust_kph",
		Help: "Current windspeed gusts.",
	})
	metricHumidity = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "observation_humidity",
		Help: "Current humidity.",
	})
	metricRainSince9Am = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "observation_rain_since_9am_mm",
		Help: "Rain since 9AM.",
	})
	metricErrorState = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "observation_error",
		Help: "Error flag indicating an observation error.",
	})
	metricBearing = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "observation_wind_bearing",
		Help: "Current wind direction.",
	})
)

func observe(name string) {
	go func() {
		w := wind.Init()
		log.Print("setting up observation loop")
		for {
			time.Sleep(time.Duration(*flagInterval) * time.Second)
			res, err := observation.Observe(name)
			if err != nil {
				log.Println(err)
				metricErrorState.Set(1)
				continue
			}
			log.Printf("received response: %#v", res)
			metricTemp.Set(res.Temp)
			metricTempFeelsLike.Set(res.TempFeelsLike)
			metricWindSpeed.Set(float64(res.Wind.SpeedKilometre))
			metricWindSpeedGust.Set(float64(res.Gust.SpeedKilometre))
			metricHumidity.Set(float64(res.Humidity))
			metricRainSince9Am.Set(res.RainSince9Am)
			metricErrorState.Set(0)
			metricBearing.Set(observation.Bearing(res.Wind.Direction))
			w.Observe(time.Now(), float64(res.Wind.SpeedKilometre), float64(res.Gust.SpeedKilometre))
		}
	}()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func Main() {
	log.SetPrefix("exporter-weather: ")
	log.SetFlags(log.LstdFlags)
	flag.Parse()
	if *flagLookup != "" {
		geo, err := location.Search(*flagLookup)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("location: %s has geohash: %s", *flagLookup, geo)
		os.Exit(0)
	}
	if *flagHelp || *flagLocation == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.Print("location: ", *flagLocation)
	log.Print("interval: ", *flagInterval)
	log.Print("listen port: ", *flagListenPort)
	listenPort := fmt.Sprintf(":%d", *flagListenPort)

	observe(*flagLocation)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", render.WindSpeed)
	http.ListenAndServe(listenPort, nil)
}
