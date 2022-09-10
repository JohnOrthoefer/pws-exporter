package main

import (
   "fmt"
   "strings"
   "encoding/json"
   "net/url"
   "github.com/eclipse/paho.mqtt.golang"
)

// Strings because the requirements for Homebridge
//   see- https://github.com/arachnetech/homebridge-mqttthing/blob/master/docs/Accessories.md#weather-station
type (
   WeatherType struct {
      Temperature    string `json:"temperature,omitempty"`  // 0.0-100.0c
      Humidity       string `json:"humidity,omitempty"`     // 0-100
      AirPressure    string `json:"airPressure,omitempty"`  // 700-1100hPa
      Weather        string `json:"weather,omitempty"`      //
      Rain1h         string `json:"rain1h,omitempty"`       // mm
      Rain24h        string `json:"rain24h,omitempty"`      // mm
      UVIndex        string `json:"uvindex,omitempty"`      //
      Visibility     string `json:"visibility,omitempty"`   // km
      WindDirection  string `json:"winddir,omitempty"`      //
      WindSpeed      string `json:"windspeed,omitempty"`    // km/h
      Active         bool   `json:"active,omitempty"`       //
      Fault          bool   `json:"fault,omitempty"`        //
      Tampered       bool   `json:"tampered,omitempty"`     //
      Battery        bool   `json:"battery,omitempty"`      //
   }
)

var (
   WeatherValues WeatherType 
   mqttClient mqtt.Client
   mqttTopic  string
)


func deg2C(degF string)string {
   c := (getFloat(degF)-32.0) * 5/9
   if c < 0 {
      c = 0
   }
   if c > 100 {
      c = 100
   }
   return fmt.Sprintf("%.1f", c)
}

func inHg2hPa(inHgStr string)string {
   inHg := getFloat(inHgStr) 
   hPa := int64(inHg * 33.86389)
   if hPa < 700 {
      hPa = 700
   }
   if hPa > 1100 {
      hPa = 1100 
   }
   
   return fmt.Sprintf("%d", hPa)
}

func setTopic(t string) {
   mqttTopic = t
}

func getTopic(v url.Values) string {
   rtn := mqttTopic
   for _, val := range []string{"ID", "softwaretype", "id", "mt", "sensor"} {
      rtn = strings.Replace(rtn, fmt.Sprintf("%%%s%%", val), v.Get(val), -1)
   }
   //fmt.Printf("Topic: %s => %s\n", mqttTopic, rtn)
   return rtn
}

func mqttSetup(broker string, topic string) {
   opts := mqtt.NewClientOptions().AddBroker(broker)
   mqttClient = mqtt.NewClient(opts)
   if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
      mqttClient = nil
      fmt.Printf("%s\n", token.Error())
      return
   }
   setTopic(topic)
}

func batteryLow(s string) bool {
   if strings.ToLower(s) == "low" {
      return true
   }
   return false
}

func inch2mm(inch string) string {
   mmeters := getFloat(inch) * 25.4
   return fmt.Sprintf("%.1f", mmeters)
}

func mph2kph(mph string) string {
   kph := getFloat(mph) * 1.609
   return fmt.Sprintf("%.0f", kph)
}

func publish(v url.Values) {
   if mqttClient == nil {
      return
   }
   WeatherValues = WeatherType{
      Temperature:   deg2C(v.Get("tempf")),
      Humidity:      v.Get("humidity"),
      AirPressure:   inHg2hPa(v.Get("baromin")),
      Rain1h:        inch2mm(v.Get("rainin")),
      Rain24h:       inch2mm(v.Get("dailyrainin")),
      WindDirection: v.Get("winddir"),
      WindSpeed:     mph2kph(v.Get("windspeedmph")),
      Battery:       batteryLow(v.Get("sensorbattery")),
   }

   jsonOut, _ := json.Marshal(WeatherValues)
   //fmt.Printf("%s\n", jsonOut)
   mqttClient.Publish(getTopic(v), 0, false, jsonOut)
}
