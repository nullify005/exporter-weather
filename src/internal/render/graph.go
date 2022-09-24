package render

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/nullify005/exporter-weather/internal/data/temperature"
	"github.com/nullify005/exporter-weather/internal/data/wind"
)

type Content struct {
	Timestamp       string
	WindSpeed       string
	WindGust        string
	Temperature     string
	FeelsLike       string
	Label           Label
	Colour          Colour
	Title           string
	RefreshInterval int
}
type Label struct {
	WindSpeed   string
	WindGust    string
	Temperature string
	FeelsLike   string
}
type Colour struct {
	WindSpeed   string
	WindGust    string
	Temperature string
	FeelsLike   string
}

var label = &Label{
	WindSpeed:   "Wind Speed",
	WindGust:    "Wind Gust",
	Temperature: "Temp",
	FeelsLike:   "Feels Like",
}
var colour = &Colour{
	WindSpeed:   "rgb(75, 192, 192)",
	WindGust:    "rgb(255, 99, 71)",
	Temperature: "rgb(75, 192, 192)",
	FeelsLike:   "rgb(255, 99, 71)",
}
var refresh int = 30 * 1000

func tsString(ts []time.Time) string {
	var ret []string
	for _, v := range ts {
		// ret = append(ret, fmt.Sprintf("%v-%02v-%02v %02v:%02v", v.Year(), int(v.Month()), v.Day(), v.Hour(), v.Minute()))
		ret = append(ret, fmt.Sprintf("%02v:%02v", v.Hour(), v.Minute()))
	}
	return "\"" + strings.Join(ret, "\",\"") + "\""
}

func floatString(i []float64) string {
	var ret []string
	for _, v := range i {
		ret = append(ret, fmt.Sprint(v))
	}
	return strings.Join(ret, ",")
}

func Graphs(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/graph.html")
	if err != nil {
		handleError(w, err)
		return
	}
	obsWind := wind.Init()
	obsTemp := temperature.Init()
	c := &Content{
		Timestamp:       tsString(obsWind.Timestamps()),
		WindSpeed:       floatString(obsWind.WindSpeed()),
		WindGust:        floatString(obsWind.WindGust()),
		Temperature:     floatString(obsTemp.Temperature()),
		FeelsLike:       floatString(obsTemp.FeelsLike()),
		Label:           *label,
		Colour:          *colour,
		RefreshInterval: refresh,
	}
	err = tmpl.Execute(w, c)
	if err != nil {
		handleError(w, err)
		return
	}
}

func handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
