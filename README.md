# exporter-weather

a little toy which collects BOM observations from a nominated location
and then outputs metrics which can be scraped by prometheus

# usage

```
  -help
    	Command line arguments
  -interval int
    	The observation Interval in Seconds (default 30)
  -location string
    	The geohash for the observation location (use lookup to find it)
  -lookup string
    	Lookup the geohash for a given location
  -port int
    	The HTTP port to listen on for metrics & health (default 2112)
```

idea is to 1st lookup a location & get it's geohash

then you can iterate over fetching updates according to the interval

# todo

incorporate https://cobra.dev/ into the cli setup