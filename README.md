# pws-exporter

## Overview
Prometheus exporter for wunderground Personal Weather Stations (PWS) protocol.

This works as a proxy for the PWS Gateway, man-in-the-middling the transaction before forwarding it up to the wunderground servers.  

Lots to add to this document about how to set it up.   That is coming soon.  I just wrote the code about a week ago and only tonight cleaned up the code enough to publish.

Right now some of the options are  just stubs and it only works as far as I know with the AcuRite 5-in-1 weather station.  

## Todo
- [ ] Add in the filter and log options
- [X] Add in an option for your own config.yaml file
- [ ] Make the station id and password options do what they should
- [ ] Write detailed directions about how to set this up
- [X] Add prometheus config
- [X] Add Dockerfile
- [X] Add makefile to clean up building
- [X] Add my Grafana dashboard
- [X] Reserve an actual prometheus port to run on
