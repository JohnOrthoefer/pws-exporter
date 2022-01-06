package main

import (
   "io"
   "fmt"
   "net/http"
   "net/url"
   "strconv"
   "time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
   "github.com/alecthomas/kong"
)

/*
https://support.weather.com/s/article/PWS-Upload-Protocol?language=en_US

2021/12/28 03:07:32 map["ID":["KMAFRAMI83"] "PASSWORD":["LaC7PgHc"] "action":["updateraw"] "baromin":["29.75"] "dailyrainin":["0.00"] "dateutc":["now"] "dewptf":["31"] "hubbattery":["normal"] "humidity":["98"] "id":["24C86E0CE519"] "mt":["5N1"] "rainin":["0.00"] "realtime":["1"] "rssi":["31"] "rtfreq":["2.5"] "sensor":["06086M"] "sensorbattery":["normal"] "softwaretype":["myAcuRite"] "tempf":["31.3"] "winddir":["180"] "windgustmph":["0"] "windspeedavgmph":["0"] "windspeedmph":["0"]]

ID: Wunderground ID
action: Always "updateraw"
baromin: [barometric pressure inches]
dailyrainin: [rain inches so far today in local time]
dewptf: [F outdoor dewpoint F]
hubbattery: 
humidity: [% outdoor humidity 0-100%]
id: Units MAC Address
mt: 
rainin: [rain inches over the past hour)]
realtime: RapidFire Updates ["1" = enabled]
rssi: Radio power
rtfreq: RapidFire Updates Frequency [2.5s]
sensor:["06086M"] 
sensorbattery:["normal"] 
softwaretype:["myAcuRite"] 
tempf:[F outdoor temperature] ["31.3"] 
winddir: [0-360 instantaneous wind direction] ["180"] 
windgustmph:[mph current wind gust, using software specific time period] ["0"] 
windspeedavgmph: over 2 minutes Maybe? ["0"] 
windspeedmph: [mph instantaneous wind speed] ["0"]

*/

type Gauges_Config struct {
   Option struct {
      Name string
      Arg  string
      Value string
      Help string
   }
}

var (
   gauges map[string]interface{}
	dateutc = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_dateutc",
			Help: "second resolution Date-time stamp",
		},
		[]string{"id"},
	)
	winddir = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_winddir",
			Help: "0-360 instantaneous wind direction",
		},
		[]string{"id"},
	)
	windspeedmph = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_windspeedmph",
			Help: "mph instantaneous wind speed",
		},
		[]string{"id"},
	)
	windgustmph = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_windgustmph",
			Help: "mph current wind gust, using software specific time period",
		},
		[]string{"id"},
	)
	windgustdir = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_windgustdir",
			Help: "0-360 using software specific time period",
		},
		[]string{"id"},
	)
	windspeedavgmph = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_windspeedavgmph",
			Help: "mph 2 minute average wind speed mph",
		},
		[]string{"id"},
	)
	tempf = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_tempf",
			Help: "F outdoor temperature",
		},
		[]string{"id"},
	)
	humidity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_humidity",
			Help: "percent outdoor humidity 0-100",
		},
		[]string{"id"},
	)
	baromin = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_baromin",
			Help: "barometric pressure inches",
		},
		[]string{"id"},
	)
	dewptf = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_dewptf",
			Help: "F outdoor dewpoint F",
		},
		[]string{"id"},
	)
	rainin = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_rainin",
			Help: "rain inches over the past hour",
		},
		[]string{"id"},
	)
	dailyrainin = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_dailyrainin",
			Help: "rain inches so far today in local time",
		},
		[]string{"id"},
	)
)

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
   Log      []string `help:"Log tags" default:"all"`
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
   req.URL.Scheme = config.upURL.Scheme
   req.URL.Host = config.upURL.Host
   fmt.Printf("Source:%s, Recieved:%s\n", req.RemoteAddr, req.URL.String())
   resp, err := http.Get(req.URL.String())
   if err != nil {
      fmt.Printf("Upload Error code, %s\n", resp.Status)
      return
   }
   defer resp.Body.Close()
   body, err := io.ReadAll(resp.Body)

   fmt.Printf("Upload Return:%s Body:%s\n", resp.Status, strconv.Quote(string(body)))
   fmt.Fprintf(w, string(body))
   v := req.URL.Query()

   id := v["ID"][0]
   label := prometheus.Labels{"id":id}
   dateutc.With(label).Set(float64(getTimestamp(v["dateutc"][0])))
   tempf.With(label).Set(getFloat(v["tempf"][0]))
   humidity.With(label).Set(getFloat(v["humidity"][0]))
   baromin.With(label).Set(getFloat(v["baromin"][0]))
   dewptf.With(label).Set(getFloat(v["dewptf"][0]))
   rainin.With(label).Set(getFloat(v["rainin"][0]))
   dailyrainin.With(label).Set(getFloat(v["dailyrainin"][0]))
   winddir.With(label).Set(getFloat(v["winddir"][0]))
   windgustmph.With(label).Set(getFloat(v["windgustmph"][0]))
   windspeedavgmph.With(label).Set(getFloat(v["windspeedavgmph"][0]))
   windspeedmph.With(label).Set(getFloat(v["windspeedmph"][0]))
}

func init() {
	prometheus.MustRegister(dateutc)
	prometheus.MustRegister(tempf)
	prometheus.MustRegister(humidity)
	prometheus.MustRegister(baromin)
	prometheus.MustRegister(dewptf)
	prometheus.MustRegister(rainin)
	prometheus.MustRegister(dailyrainin)
	prometheus.MustRegister(winddir)
	prometheus.MustRegister(windgustmph)
	prometheus.MustRegister(windspeedavgmph)
	prometheus.MustRegister(windspeedmph)
}

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

   http.HandleFunc(config.Path, weather)
   http.Handle(config.Metrics, promhttp.Handler())
   err = http.ListenAndServe(config.Listen, nil)
   ctx.FatalIfErrorf(err)
}

