package bom

const (
	locationApiPrefix string = "/v1/locations"
	locationApiSuffix string = "/observations"
	locationSearchQS  string = "search"
	defaultEndpoint   string = "https://api.weather.bom.gov.au"
	geoSize           int    = 6
)

var (
	endpoint string = defaultEndpoint
)

func SetEndpoint(name string) {
	endpoint = name
}

func ResetEndpoint() {
	endpoint = defaultEndpoint
}
