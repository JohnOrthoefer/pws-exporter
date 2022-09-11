SRC= main.go build.go options.go tags.go weather.go mqtt.go
GOLANG=/usr/bin/go
CURL=/usr/bin/curl
GIT=/usr/bin/git
REPONAME=$(shell basename `git rev-parse --show-toplevel`)
DOCKERREPO=${REPONAME}
VERSION=v1.0-beta
SHA1=$(shell git rev-parse --short HEAD)
NOW=$(shell date +%Y-%m-%d_%T)

pws_exporter: ${SRC}
	echo ${REPONAME}
	${GOLANG} build \
		-o pws_exporter \
		-ldflags "-X main.sha1ver=${SHA1} \
		-X main.buildTime=${NOW} \
		-X main.repoName=${REPONAME}"

Docker: pws_exporter
	docker build -t ${DOCKERREPO}:${VERSION} .

update-go:
	${GOLANG} get github.com/prometheus/client_golang/prometheus
	${GOLANG} get github.com/prometheus/client_golang/prometheus/promhttp
	${GOLANG} get github.com/alecthomas/kong
	${GOLANG} get github.com/eclipse/paho.mqtt.golang
	${GOLANG} get gopkg.in/yaml.v2
	${GOLANG} mod tidy

clean:
	rm -f pws-exporter pws_exporter

