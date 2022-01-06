package main

import (
   "os"
   "fmt"
   "log"
   "io/ioutil"
   "net/http"
   "net/url"
   "github.com/prometheus/client_golang/prometheus"
   "github.com/prometheus/client_golang/prometheus/promhttp"
   "github.com/alecthomas/kong"
   "gopkg.in/yaml.v2"
)

var (
   gauge GaugeMap
)

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

   yamlTags := yamlDefault
   if config.Tags != "" {
      yamlTags, err = ioutil.ReadFile(config.Tags)
      if err != nil {
         fmt.Printf("Error: %s\n", err)
         os.Exit(1)
      }
   }

   if config.Version || config.Verbose {
      fmt.Printf("%s - %s @%s\n", repoName, sha1ver, buildTime)
   }

   if config.Verbose {
      fmt.Printf("\nConfig-\n")
      fmt.Printf("\tListen: %s\n", config.Listen)
      fmt.Printf("\tPath: %s\n", config.Path)
      fmt.Printf("\tMetrics: %s\n", config.Metrics)
      fmt.Printf("\tUpload URL: %s\n", config.upURL.String())
      tagType := "Built-in"
      if config.Tags != "" {
         tagType = config.Tags
      }
      fmt.Printf("\tTags: %s\n", tagType)
   }

   err = yaml.Unmarshal(yamlTags, &gauge)
   if err != nil {
      log.Fatalf("YAML Parse Error, %s", err)
   }
   // see if we can ingest the yaml
   if config.DumpYAML {
      fmt.Printf("%s\n", string(yamlTags))
      os.Exit(0)
   }


   for key, v := range gauge.Tags {
      if v.Value != nil {
         continue
      }

      if v.Name == "" {
         gauge.Tags[key].Name = config.Prefix + key
      }

      // in the future could support counter if needed
      if v.Type != "gauge" {
         fmt.Printf("Ignoring %s, only gauge is supported\n")
         continue
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

