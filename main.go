package main

import (
   "io"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "strconv"
   "time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
   "github.com/alecthomas/kong"
   "gopkg.in/yaml.v2"
)

var (
   gauge GaugeMap
)

/* 
https://support.weather.com/s/article/PWS-Upload-Protocol

All the fields are defined on the Weather.com Website.
*/

var config struct {
   Upstream string  `help:"URL to Receive" default:"https://weatherstation.wunderground.com/weatherstation/updateweatherstation.php"`
   upURL    *url.URL 
   Forward  bool     `negatable:"" help:"Forward data to Upstream" default:"false"`
   Path     string   `help:"Path to scrape, defaults to Upstream Path"`
   Metrics  string   `help:"URL to be scraped" default:"/metrics"`
   Listen   string   `help:"Address to listen on" default:":8443"`
   ID       string   `help:"Overrides ID to be used Upstream and label data"`
   KEY      string   `help:"Overrides Key to be used Upstream"`
   Prefix   string   `help:"Prometheus prefix for metrics" default:"weather_"`
   Log      []string `help:"Tags that will be provided to Prometheus"`
   Filter   []string `help:"Filter tags, do not log or forward, superceeds --Log"`
}

func getFloat(s string) float64 {
   v, err := strconv.ParseFloat(s, 64)
   if err != nil {
      nan, _ := strconv.ParseFloat("NaN", 64)
      v = nan
   }
   return v
}

func getTimestamp(ts string) int64 {
   if ts == "now" {
      t := time.Now().Unix()
      return t
   } 
   t, err := time.Parse("2006-01-02 15:04:05", ts)
   if err != nil {
      t = time.Now()
   }
   return t.Unix()
}

func weather(w http.ResponseWriter, req *http.Request) {
   // get the Query options
   v := req.URL.Query()

   // check for the required tags
   // action=updateraw is required value
   if action, ok := v["action"]; !ok || action[0] != "updateraw" {
      fmt.Printf("Action is not 'updateraw', Required\n")
      return
   }
   // ID is requried
   if _, ok := v["ID"]; !ok {
      fmt.Printf("No wunderground.com ID/Login (REQUIRED)")
      return
   }
   id := v["ID"][0]
   // PASSWORD
   if _, ok := v["PASSWORD"]; !ok {
      fmt.Printf("No wunderground.com Password (REQUIRED)")
      return
   }
   //passwd := v["PASSWORD"][0]
   // date time is requried
   if _, ok := v["dateutc"]; !ok {
      fmt.Printf("No UTC Date (REQUIRED)")
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

// main
func main() {
   var err error
   ctx := kong.Parse(&config,
	   kong.Name("pws-exporter"),
	   kong.Description("Prometheus Exporter to log Weather Underground updates."))

   config.upURL, err = url.Parse(config.Upstream)
   ctx.FatalIfErrorf(err)

   if config.Path == "" {
      config.Path = config.upURL.Path
   }

   fmt.Printf("Listen: %s\n", config.Listen)
   fmt.Printf("Path: %s\n", config.Path)
   fmt.Printf("Metrics: %s\n", config.Metrics)
   fmt.Printf("Upload URL: %s\n", config.upURL.String())

   err = yaml.Unmarshal([]byte(yamlDefault), &gauge)
   if err != nil {
      log.Fatalf("YAML Parse Error, %s", err)
   }

   for key, v := range gauge.Tags {
      if v.Value != nil {
         continue
      }
      if v.Name == "" {
         gauge.Tags[key].Name = config.Prefix + key
      }
      gauge.Tags[key].Value = prometheus.NewGaugeVec(
                  prometheus.GaugeOpts{
                     Name: v.Name,
                     Help: v.Help,
                  },
                  []string{"id"},
               )
      prometheus.MustRegister(v.Value)
      for _, alias := range v.Alias {
         gauge.Tags[alias] = v
      }
   }

   http.HandleFunc(config.Path, weather)
   http.Handle(config.Metrics, promhttp.Handler())
   err = http.ListenAndServe(config.Listen, nil)
   ctx.FatalIfErrorf(err)
}

