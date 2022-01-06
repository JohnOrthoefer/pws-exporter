package main 

import (
   "github.com/prometheus/client_golang/prometheus"
)

const yamlDefault = `
tags:
  dateutc:
    help: second resolution Date-time stamp
  winddir:
    help: 0-360 instantaneous wind direction
  windspeedmph:
    help: mph instantaneous wind speed
  windgustmph:
    help: mph current wind gust, using software specific time period
  windgustdir:
    help: 0-360 using software specific time period
  windspdmph_avg2m:
    help: mph 2 minute average wind speed mph
    alias: [windspeedavgmph]
  winddir_avg2m:
    help: 0-360 2 minute average wind direction
  windgustmph_10m:
    help: mph past 10 minutes wind gust mph 
  windgustdir_10m:
    help: 0-360 past 10 minutes wind gust direction
  humidity:
    help: '% outdoor humidity 0-100%'
  dewptf:
    help: F outdoor dewpoint F
  tempf:
    help: F outdoor temperature
  rainin:
    help: rain inches over the past hour, the accumulated rainfall in the past 60 min
  dailyrainin:
    help: rain inches so far today in local time
  baromin:
    help: barometric pressure inches
  soiltempf:
    help: F soil temperature
  soilmoisture:
    help: '%'
  leafwetness:
    help: '%'
  solarradiation:
    help: W/m^2
  UV:
    help: index
  visibility:
    help: nm visibility
  indoortempf:
    help: F indoor temperature F
  indoorhumidity:
    help: '% indoor humidity 0-100'
`
 
type GaugeEntry struct {
   Name string `yaml:"name,omitempty"`
   Help string `yaml:"help,omitempty"`
   Alias []string `yaml:"alias,omitempty"`
   Value *prometheus.GaugeVec
}

type GaugeMap struct {
  Tags map[string]*GaugeEntry `yaml:"tags,omitempty"`
}
