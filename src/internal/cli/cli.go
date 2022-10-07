package cli

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nullify005/exporter-weather/internal/bom"
	"github.com/nullify005/exporter-weather/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var flagLocation = flag.String("location", "", "The geohash for the observation location (use lookup to find it)")
var flagInterval = flag.Int("interval", 30, "The observation Interval in Seconds")
var flagListen = flag.String("listen", "127.0.0.1:2112", "The HTTP port to listen on for metrics & health")
var flagHelp = flag.Bool("help", false, "Command line arguments")
var flagLookup = flag.String("lookup", "", "Lookup the geohash for a given location")

func watch(name string) {
	go func() {
		log.Print("setting up observation loop")
		m := metrics.New()
		for {
			time.Sleep(time.Duration(*flagInterval) * time.Second)
			res, err := bom.Observe(name)
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func Main() {
	log.SetPrefix("exporter-weather: ")
	log.SetFlags(log.LstdFlags)
	flag.Parse()
	if *flagLookup != "" {
		locations, err := bom.Location(*flagLookup)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if len(locations) <= 0 {
			fmt.Printf("no locations found for: %s", *flagLookup)
			os.Exit(1)
		}
		for _, v := range locations {
			fmt.Println(v.String())
		}
		os.Exit(0)
	}
	if *flagHelp || *flagLocation == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.Print("location: ", *flagLocation)
	log.Print("interval: ", *flagInterval)
	log.Print("listen: ", *flagListen)

	watch(*flagLocation)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", healthHandler)
	http.ListenAndServe(*flagListen, nil)
}
