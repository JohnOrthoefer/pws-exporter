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
