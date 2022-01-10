# pws-exporter

## Overview
Prometheus exporter for wunderground Personal Weather Stations (PWS) protocol.

This works as a proxy for the PWS Gateway, man-in-the-middling the transaction before forwarding it up to the wunderground servers.  

Lots to add to this document about how to set it up.   That is coming soon.  I just wrote the code about a week ago and only tonight cleaned up the code enough to publish.

Right now some of the options are  just stubs and it only works as far as I know with the AcuRite 5-in-1 weather station.  

## Setup
* Make and Install  pws_exportor
```
   make
   sudo cp pws_exporter /usr/local/bin/pws_exporter
   sudo cp config/systemd/wuproxy.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable --now wuproxy.service
```

* Add the following to nginx with a self-signed SSL.
```
   cp config/nginx/wuproxy.conf /etc/nginx/conf.d/
   openssl req -x509 -newkey rsa:4096 \
      -keyout /etc/nginx/pki/key-enc.pem \
      -out /etc/nginx/pki/cert.pem \
      -subj '/CN=localhost' \
      -sha256 -days 3650
   openssl rsa -in /etc/nginx/pki/key-enc.pem \
      -out /etc/nginx/pki/key.pem
   sudo systemctl restart  nginx
```

* Add the following A record to the local DNS
```
   rtupdate.wunderground.com.  A  ${RPi_StaticIP}
```

## Run-time Options
* `--upstream=` Sets the address will accept the pass though post.  It also sets the path that is intercepted. Default https://weatherstation.wunderground.com/weatherstation/updateweatherstation.php
* `--forward` Flag allows the software to forward to the upstream address.  Default is do not forward.
* `--path=` the scrape path for Promethus.  Default is /metrics
* `--listen=` what port to listen in.  Default is all inteface on port 9874.
* `--id=` Station ID, Default is what is passed by the device.
* `--key=` Station Key, Default is what is passed by the device.
* `--prefix=` prefix that is prepended to the prometheus keys.  Default is "weather_".
* `--tags=` provide a file for external tags.yaml. Default is use the built-in.
* `--filter=` list of options that will be filtered out.  Default is to not filter anything.
* `--version` print the build information.
* `--verbose` print extra logging.
* `--dump-yaml` print the built-in Yaml.

## Todo
- [X] Add in the filter option
- [ ] Add in the log option
- [X] Add in an option for your own config.yaml file
- [X] Make the station id and password options do what they should
- [ ] Write detailed directions about how to set this up
- [X] Add prometheus config
- [X] Add Dockerfile
- [X] Add makefile to clean up building
- [X] Add my Grafana dashboard
- [X] Reserve an actual prometheus port to run on
