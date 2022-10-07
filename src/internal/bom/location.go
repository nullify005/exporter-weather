package bom

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type LocationResponse struct {
	Metadata LocationMetadata `json:"metadata"`
	Data     []LocationData   `json:"data"`
	Errors   []struct {
		Code   string `json:"code"`
		Title  string `json:"titls"`
		Status int    `json:"status"`
		Detail string `json:"detail"`
	} `json:"errors"`
}
type LocationMetadata struct {
	ResponseTimestamp time.Time `json:"response_timestamp"`
	Copyright         string    `json:"copyright"`
}
type LocationData struct {
	Geohash  string `json:"geohash"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Postcode string `json:"postcode"`
	State    string `json:"state"`
}

func Location(name string) (l []LocationData, err error) {
	uri := endpoint + locationApiPrefix
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return
	}
	q := url.Values{}
	q.Add(locationSearchQS, name)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusBadRequest:
	default:
		err = fmt.Errorf("%s invalid status: %v body: %s", uri, resp.StatusCode, body)
		return
	}
	lResp := &LocationResponse{}
	if err = json.Unmarshal(body, &lResp); err != nil {
		return
	}
	// when the resp.StatusCode == http.StatusBadRequest
	for _, v := range lResp.Errors {
		err = fmt.Errorf("%s %s %s", v.Code, v.Title, v.Detail)
		return
	}
	if lResp.Metadata.Copyright == "" {
		err = fmt.Errorf("missing copyright response in payload")
		return
	}
	l = lResp.Data
	return
}

func (l *LocationData) String() string {
	return fmt.Sprintf("(%s) id: %s postcode: %s state: %s geohash: %s", l.Name, l.ID, l.Postcode, l.State, l.Geohash)
}
