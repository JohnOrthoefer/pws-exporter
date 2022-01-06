package main 
// embedded yaml for the default wu rules

import (
   "github.com/prometheus/client_golang/prometheus"
)

var yamlDefault = []byte (`
#
# Default YAML built in
# 
tags:
  dateutc:
    help: second resolution Date-time stamp
    type: gauge
  winddir:
    help: 0-360 instantaneous wind direction
    type: gauge
  windspeedmph:
    help: mph instantaneous wind speed
    type: gauge
  windgustmph:
    help: mph current wind gust, using software specific time period
    type: gauge
  windgustdir:
    help: 0-360 using software specific time period
    type: gauge
  windspdmph_avg2m:
    help: mph 2 minute average wind speed mph
    type: gauge
    # the Acurite 5-in-1 provides this option named differently
    alias: 
      - windspeedavgmph
  winddir_avg2m:
    help: 0-360 2 minute average wind direction
    type: gauge
  windgustmph_10m:
    help: mph past 10 minutes wind gust mph 
    type: gauge
  windgustdir_10m:
    help: 0-360 past 10 minutes wind gust direction
    type: gauge
  humidity:
    help: '% outdoor humidity 0-100%'
    type: gauge
  dewptf:
    help: F outdoor dewpoint F
    type: gauge
  tempf:
    help: F outdoor temperature
    type: gauge
  rainin:
    help: rain inches over the past hour, the accumulated rainfall in the past 60 min
    type: gauge
  dailyrainin:
    help: rain inches so far today in local time
    type: gauge
  baromin:
    help: barometric pressure inches
    type: gauge
  soiltempf:
    help: F soil temperature
    type: gauge
  soilmoisture:
    help: '%'
    type: gauge
  leafwetness:
    help: '%'
    type: gauge
  solarradiation:
    help: W/m^2
    type: gauge
  UV:
    help: index
    type: gauge
  visibility:
    help: nm visibility
    type: gauge
  indoortempf:
    help: F indoor temperature F
    type: gauge
  indoorhumidity:
    help: '% indoor humidity 0-100'
    type: gauge
`)
 
type GaugeEntry struct {
   Name string `yaml:"name,omitempty"`
   Help string `yaml:"help,omitempty"`
   Type string `yaml:"type,omitempty"`
   Alias []string `yaml:"alias,omitempty"`
   Value *prometheus.GaugeVec
}

type GaugeMap struct {
  Tags map[string]*GaugeEntry `yaml:"tags,omitempty"`
}
