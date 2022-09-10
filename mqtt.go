package main

import (
   "fmt"
   "encoding/json"
   "net/url"
   "github.com/eclipse/paho.mqtt.golang"
)

// Strings because the requirements for Homebridge are-
// Current temperature must be in the range 0 to 100 degrees Celsius to a maximum of 1dp.
// Current relative humidity must be in the range 0 to 100 percent with no decimal places.
// Air Pressure must be in the range 700 to 1100 hPa.
var (
   WeatherValues struct {
      Temperature string `json:"temperature"`
      Humidity string    `json:"humidity"`
      AirPressure string `json:"airPressure"`
   }
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

func mqttSetup(broker string, topic string) {
   opts := mqtt.NewClientOptions().AddBroker(broker)
   mqttClient = mqtt.NewClient(opts)
   if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
      mqttClient = nil
      fmt.Printf("%s\n", token.Error())
      return
   }
   mqttTopic = topic
}

func publish(v url.Values) {
   if mqttClient == nil {
      return
   }
   WeatherValues.Temperature = deg2C(v.Get("tempf"))
   WeatherValues.Humidity = v.Get("humidity")
   WeatherValues.AirPressure = inHg2hPa(v.Get("baromin"))

   jsonOut, _ := json.Marshal(WeatherValues)
   fmt.Printf("%s\n", jsonOut)
   mqttClient.Publish(mqttTopic, 0, false, jsonOut)
}
