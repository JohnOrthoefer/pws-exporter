package main
import (
   "net/url"
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
   TagsYAML string   `help:"YAML file with the tags defined."`
   Log      []string `help:"Tags that will be provided to Prometheus"`
   Filter   []string `help:"Filter tags, do not log or forward, superceeds --Log"`
}
