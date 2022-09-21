package temperature

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
	timestamp   time.Time
	temperature float64
	feelslike   float64
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

func (*observations) Temperature() []float64 {
	var r []float64
	for _, v := range o.observation {
		r = append(r, v.temperature)
	}
	return r
}

func (*observations) FeelsLike() []float64 {
	var r []float64
	for _, v := range o.observation {
		r = append(r, v.feelslike)
	}
	return r
}

func (*observations) Observe(t time.Time, temp float64, feels float64) {
	o.observation = append(o.observation, observation{timestamp: t, temperature: temp, feelslike: feels})
}
