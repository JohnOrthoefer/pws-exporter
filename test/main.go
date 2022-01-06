package main

import (
   "fmt"
   "log"
   "gopkg.in/yaml.v2"
)

const yamlSimple = `
tags:
  foo:
    bar: other
    baz: what
  this: that
`

const yamlDefault = `
tags:
  winddir:
    name: Wind Direction
    help: 0-360 instantaneous wind direction
  windspeedmph:
    name: Wind Speed
    help: mph instantaneous wind speed
  windgustmph:
    name: Wind Gust
    help: mph current wind gust, using software specific time period
  windgustdir:
    name: Wind Gust Direction 
    help: 0-360 using software specific time period
  windspdmph_avg2m:
    name: Wind Average
    help: mph 2 minute average wind speed mph
  winddir_avg2m:
    name: Wind Direction
    help: 0-360 2 minute average wind direction
  windgustmph_10m:
    name: Wind Gust Average
    help: mph past 10 minutes wind gust mph 
  windgustdir_10m:
    name: Wind Gust Direction
    help: 0-360 past 10 minutes wind gust direction
  humidity:
    name: Humidity
    help: '% outdoor humidity 0-100%'
  dewptf:
    name: Dew Point
    help: F outdoor dewpoint F
  tempf:
    name: Temperature
    help: F outdoor temperature
  rainin:
    name: Rain
    help: rain inches over the past hour, the accumulated rainfall in the past 60 min
  dailyrainin:
    name: Rain Daily
    help: rain inches so far today in local time
  baromin:
    name: Pressure
    help: barometric pressure inches
  soiltempf:
    name: Soil Temp
    help: F soil temperature
  soilmoisture:
    name: Soil Moisture
    help: '%'
  leafwetness:
    name: Leaf Wetness
    help: '%'
  solarradiation:
    name: Radiation
    help: W/m^2
  UV:
    name: UV Index
    help: index
  visibility:
    name: Visibility
    help: nm visibility
  indoortempf:
    name: Indoor Temp
    help: F indoor temperature F
  indoorhumidity:
    name: Indoor Humifity
    help: '% indoor humidity 0-100'
  softwaretype:
    name: Software type
    help: '[text] ie: WeatherLink, VWS, WeatherDisplay'
`
 
type GaugeEntry struct {
   Name string `yaml:"name"`
   Help string `yaml:"help"`
}
type GaugeMap struct {
  Tags map[string]GaugeEntry `yaml:"tags,omitempty"`
}

func main() {

/*   t := make(GaugeMap) */
   var t GaugeMap

   err := yaml.Unmarshal([]byte(yamlDefault), &t)
   /* err := yaml.Unmarshal([]byte(yamlSimple), &t)*/
   if err != nil {
      log.Fatalf("error parsing (%s)\n%q", err, t)
   }
   fmt.Printf("%q", t)
}
