package location

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type response struct {
	Metadata metadata `json:"metadata"`
	Data     []data   `json:"data"`
}
type metadata struct {
	ResponseTimestamp time.Time `json:"response_timestamp"`
	Copyright         string    `json:"copyright"`
}
type data struct {
	Geohash  string `json:"geohash"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Postcode string `json:"postcode"`
	State    string `json:"state"`
}

var endpoint string = "https://api.weather.bom.gov.au/v1/locations"

func Search(name string) (string, error) {
	resp, err := http.Get(endpoint + "?search=" + name)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid response from endpoint: %s status: %v content: %s", endpoint, resp.StatusCode, body)
	}
	if err != nil {
		return "", err
	}
	response := response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	if len(response.Data) > 1 {
		names := ""
		for _, v := range response.Data {
			s := fmt.Sprintf("   - %s (%s) [%s]\n", v.Name, v.State, v.Geohash)
			names += s
		}
		names = names[:len(names)-1] // strip the last char
		return "", fmt.Errorf("found multiple locations:\n%s\ntry narrowing the search", names)
	}
	return response.Data[0].Geohash, nil
}
