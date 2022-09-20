package render

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	wind "github.com/nullify005/exporter-weather/internal/data"
)

type Content struct {
	Timestamp string
	WindSpeed string
	WindGust  string
	Label     Label
	Colour    Colour
}
type Label struct {
	WindSpeed string
	WindGust  string
}
type Colour struct {
	WindSpeed string
	WindGust  string
}

var label = &Label{
	WindSpeed: "Wind Speed kph",
	WindGust:  "Wind Gust kph",
}
var colour = &Colour{
	WindSpeed: "rgb(75, 192, 192)",
	WindGust:  "rgb(255, 99, 71)",
}

func tsString(ts []time.Time) string {
	var ret []string
	for _, v := range ts {
		ret = append(ret, fmt.Sprintf("%v-%02v-%02v %02v:%02v", v.Year(), int(v.Month()), v.Day(), v.Hour(), v.Minute()))
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

func WindSpeed(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/graph.html")
	if err != nil {
		handleError(w, err)
		return
	}
	obs := wind.Init()
	c := &Content{
		Timestamp: tsString(obs.Timestamps()),
		WindSpeed: floatString(obs.WindSpeed()),
		WindGust:  floatString(obs.WindGust()),
		Label:     *label,
		Colour:    *colour,
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
