package main

import (
   "io"
   "fmt"
   "net/http"
   "strconv"
   "time"
   "github.com/prometheus/client_golang/prometheus"
)

/* 
https://support.weather.com/s/article/PWS-Upload-Protocol

All the fields are defined on the Weather.com Website.
*/

// turn a string into Float64
func getFloat(s string) float64 {
   v, err := strconv.ParseFloat(s, 64)
   if err != nil {
      nan, _ := strconv.ParseFloat("NaN", 64)
      v = nan
   }
   return v
}

// turn a timestamp into an int64
// -- this function always returns a time
func getTimestamp(ts string) int64 {
   // the spec allows for the uploader to not provide a time
   // if that is the case the the string "now" is used
   if ts == "now" {
      t := time.Now().Unix()
      return t
   } 
   t, err := time.Parse("2006-01-02 15:04:05", ts)
   if err != nil {
      // did not parse, just return now
      t = time.Now()
   }
   return t.Unix()
}

// http handler function for posts of data
func weather(w http.ResponseWriter, req *http.Request) {
   // get the Query options
   v := req.URL.Query()

   // check for the required tags
   // action=updateraw is required value
   if action, ok := v["action"]; !ok || action[0] != "updateraw" {
      fmt.Printf("Action is not 'updateraw' (REQUIRED)\n")
      return
   }
   // ID is requried
   if _, ok := v["ID"]; !ok {
      fmt.Printf("No wunderground.com ID/Login (REQUIRED)\n")
      return
   }
   id := v["ID"][0]
   // PASSWORD
   if _, ok := v["PASSWORD"]; !ok {
      fmt.Printf("No wunderground.com Password (REQUIRED)\n")
      return
   }
   //passwd := v["PASSWORD"][0]
   // date time is requried
   if _, ok := v["dateutc"]; !ok {
      fmt.Printf("No UTC Date (REQUIRED)\n")
      return
   }
   dateutc := v["dateutc"][0]

   // do you want to forward upstream
   if config.Forward {
      // forward up stream

      // update the target of the URL and send it
      req.URL.Scheme = config.upURL.Scheme
      req.URL.Host = config.upURL.Host
      fmt.Printf("Source:%s, Recieved:%s\n", req.RemoteAddr, req.URL.String())
      resp, err := http.Get(req.URL.String())
      if err != nil {
         // Toss the read
         fmt.Printf("Upload Error code, %s\n", resp.Status)
         return
      }
      defer resp.Body.Close()
      body, err := io.ReadAll(resp.Body)

      fmt.Printf("Upload Return:%s Body:%s\n", resp.Status, 
         strconv.Quote(string(body)))
      fmt.Fprintf(w, string(body))
   } else {
      fmt.Printf("Skipping forwarding.  Success!\n")
      fmt.Fprintf(w, "success\n")
   }

   // Label the data with the station login/name
   label := prometheus.Labels{"id":id}

   // time is a gauge and handled outside the loop
   gauge.Tags["dateutc"].Value.With(label).Set(float64(getTimestamp(dateutc)))
   delete(v, "dateutc")

   // loop though what has been provided and get them ready to be scraped
   for k, val := range v {
      if _, ok := gauge.Tags[k]; !ok {
         continue
      }
      gauge.Tags[k].Value.With(label).Set(float64(getFloat(val[0])))
   }
}

