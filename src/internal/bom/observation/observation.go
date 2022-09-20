package observation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var endpoint string = "https://api.weather.bom.gov.au/v1/locations/"
var maxGeo int = 6

type Response struct {
	Metadata Metadata `json:"metadata"`
	Data     Data     `json:"data"`
}
type Metadata struct {
	ResponseTimestamp time.Time `json:"response_timestamp"`
	IssueTime         time.Time `json:"issue_time"`
	ObservationTime   time.Time `json:"observation_time"`
	Copyright         string    `json:"copyright"`
}
type Wind struct {
	SpeedKilometre int    `json:"speed_kilometre"`
	SpeedKnot      int    `json:"speed_knot"`
	Direction      string `json:"direction"`
}
type Gust struct {
	SpeedKilometre int `json:"speed_kilometre"`
	SpeedKnot      int `json:"speed_knot"`
}
type MaxGust struct {
	SpeedKilometre int       `json:"speed_kilometre"`
	SpeedKnot      int       `json:"speed_knot"`
	Time           time.Time `json:"time"`
}
type MaxTemp struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}
type MinTemp struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}
type Station struct {
	BomID    string `json:"bom_id"`
	Name     string `json:"name"`
	Distance int    `json:"distance"`
}
type Data struct {
	Temp          float64 `json:"temp"`
	TempFeelsLike float64 `json:"temp_feels_like"`
	Wind          Wind    `json:"wind"`
	Gust          Gust    `json:"gust"`
	MaxGust       MaxGust `json:"max_gust"`
	MaxTemp       MaxTemp `json:"max_temp"`
	MinTemp       MinTemp `json:"min_temp"`
	RainSince9Am  float64 `json:"rain_since_9am"`
	Humidity      int     `json:"humidity"`
	Station       Station `json:"station"`
}

var bearing = map[string]float64{
	"N":   0,
	"NNE": 22.5,
	"NE":  45,
	"ENE": 67.5,
	"E":   90,
	"ESE": 112.5,
	"SE":  135,
	"SSE": 157.5,
	"S":   180,
	"SSW": 202.5,
	"SW":  225,
	"WSW": 247.5,
	"W":   270,
	"WNW": 292.5,
	"NW":  315,
	"NNW": 337.5,
}

func Observe(geo string) (Data, error) {
	data := Data{}
	if len(geo) > maxGeo {
		// fmt.Printf("truncating geohash: %s down to: %v\n", geo, maxGeo)
		geo = geo[:len(geo)-(len(geo)-maxGeo)]
		// fmt.Printf("result: %s\n", geo)
	}
	u := endpoint + geo + "/observations"
	resp, err := http.Get(u)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("invalid response from endpoint: %s status: %v content: %s", u, resp.StatusCode, body)
	}
	response := Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return data, err
	}
	data = response.Data
	return data, nil
}

func Bearing(direction string) float64 {
	return bearing[direction]
}
