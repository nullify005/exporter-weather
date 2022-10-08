package bom

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

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

type ObservationResponse struct {
	Metadata struct {
		ResponseTimestamp time.Time `json:"response_timestamp"`
		IssueTime         time.Time `json:"issue_time"`
		ObservationTime   time.Time `json:"observation_time"`
		Copyright         string    `json:"copyright"`
	} `json:"metadata"`
	Data Observation `json:"data"`
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
type Observation struct {
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
	Errors        []struct {
		Code   string `json:"code"`
		Title  string `json:"titls"`
		Status int    `json:"status"`
		Detail string `json:"detail"`
	} `json:"errors"`
}

func Observe(location string) (data Observation, err error) {
	loc := shrink(location)
	url := endpoint + locationApiPrefix + "/" + loc + locationApiSuffix
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusBadRequest:
	default:
		err = fmt.Errorf("%s invalid status: %v body: %s", url, resp.StatusCode, body)
		return
	}
	response := &ObservationResponse{}
	if err = json.Unmarshal(body, &response); err != nil {
		return
	}
	// when the resp.StatusCode == http.StatusBadRequest
	for _, v := range response.Data.Errors {
		err = fmt.Errorf("%s %s %s", v.Code, v.Title, v.Detail)
		return
	}
	if response.Data.Station.BomID == "" {
		err = fmt.Errorf("station id missing in payload")
		return
	}
	data = response.Data
	return
}

func WindBearing(direction string) float64 {
	if _, ok := bearing[direction]; !ok {
		return -1
	}
	return bearing[direction]
}

func (o *Observation) String() string {
	return fmt.Sprintf("(%s) temp: %vC humidity: %v%% wind: %s@%vkph", o.Station.Name, o.Temp, o.Humidity, o.Wind.Direction, o.Wind.SpeedKilometre)
}

func shrink(geo string) (ret string) {
	ret = geo
	if len(ret) > geoSize {
		ret = ret[:len(ret)-(len(ret)-geoSize)]
	}
	return
}
