package bom

import (
	"net/http"
	"testing"

	"github.com/nullify005/exporter-weather/internal/test"
	"github.com/stretchr/testify/assert"
)

const (
	_respObsError string  = "./assets/errorObsResponse.json"
	_respObsValid string  = "./assets/validObsResponse.json"
	_respNil      string  = "./assets/nilResponse.json"
	_respLocError string  = "./assets/errorLocResponse.json"
	_respLocValid string  = "./assets/validLocResponse.json"
	_stationID    string  = "012345"
	_stationName  string  = "Sydney"
	_temperature  float64 = 8.8
)

func TestObsErrors(t *testing.T) {
	tests := []struct {
		name       string
		httpStatus int
		payload    string
		errMsg     string
	}{
		{
			name:       "503",
			httpStatus: http.StatusBadGateway,
			payload:    _respNil,
			errMsg:     "",
		},
		{
			name:       "not found",
			httpStatus: http.StatusNotFound,
			payload:    _respObsError,
			errMsg:     "Invalid Geohash",
		},
		{
			name:       "bad payload",
			httpStatus: http.StatusOK,
			payload:    _respObsError,
			errMsg:     "station id missing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := test.MockHTTPServer(tt.httpStatus, tt.payload)
			if err != nil {
				t.Errorf(err.Error())
				return
			}
			SetEndpoint(s.URL)
			_, err = Observe(_stationID)
			if tt.errMsg == "" {
				assert.Error(t, err)
				return
			}
			assert.ErrorContains(t, err, tt.errMsg)
		})
	}
}

func TestObsValid(t *testing.T) {
	s, err := test.MockHTTPServer(http.StatusOK, _respObsValid)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	SetEndpoint(s.URL)
	d, err := Observe(_stationID)
	assert.NoError(t, err)
	assert.Equal(t, _temperature, d.TempFeelsLike)
	assert.Equal(t, _stationID, d.Station.BomID)
}

func TestLocErrors(t *testing.T) {
	tests := []struct {
		name       string
		httpStatus int
		payload    string
		errMsg     string
	}{
		{
			name:       "503",
			httpStatus: http.StatusBadGateway,
			payload:    _respNil,
			errMsg:     "",
		},
		{
			name:       "not found",
			httpStatus: http.StatusNotFound,
			payload:    _respLocError,
			errMsg:     "Invalid search query",
		},
		{
			name:       "bad payload",
			httpStatus: http.StatusOK,
			payload:    _respNil,
			errMsg:     "missing copyright response in payload",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := test.MockHTTPServer(tt.httpStatus, tt.payload)
			if err != nil {
				t.Errorf(err.Error())
				return
			}
			SetEndpoint(s.URL)
			_, err = Location(_stationName)
			if tt.errMsg == "" {
				assert.Error(t, err)
				return
			}
			assert.ErrorContains(t, err, tt.errMsg)
		})
	}
}

func TestLocValid(t *testing.T) {
	s, err := test.MockHTTPServer(http.StatusOK, _respLocValid)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	SetEndpoint(s.URL)
	d, err := Location(_stationName)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(d))
	assert.Equal(t, _stationName, d[1].Name)
}
