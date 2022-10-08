# exporter-weather

a little toy which collects BoM observations from a nominated location
and then outputs metrics which can be scraped by prometheus

# usage

```
./exporter-weather (search|watch) [-i time.Duration] [-l host:port]

where:
   - search <name>
     returns location & geohash data from BoM
   - watch <geohash>
     [-i time.Duration] /* the loop interval for observations */
     [-l host:port] /* the hostname & port to expose metrics & health on */
```

idea is to 1st lookup a location & get it's geohash

then you can iterate over fetching updates according to the interval
