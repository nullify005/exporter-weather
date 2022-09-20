package wind

import (
	"sync"
	"time"
)

var lock = &sync.Mutex{}
var o *observations

type observations struct {
	observation []observation
}

type observation struct {
	timestamp time.Time
	windSpeed float64
	windGust  float64
}

func Init() *observations {
	lock.Lock()
	defer lock.Unlock()
	if o == nil {
		o = &observations{}
	}
	return o
}

func (*observations) Timestamps() []time.Time {
	var r []time.Time
	for _, v := range o.observation {
		r = append(r, v.timestamp)
	}
	return r
}

func (*observations) WindSpeed() []float64 {
	var r []float64
	for _, v := range o.observation {
		r = append(r, v.windSpeed)
	}
	return r
}

func (*observations) WindGust() []float64 {
	var r []float64
	for _, v := range o.observation {
		r = append(r, v.windGust)
	}
	return r
}

func (*observations) Observe(t time.Time, speed float64, gust float64) {
	o.observation = append(o.observation, observation{timestamp: t, windSpeed: speed, windGust: gust})
}
