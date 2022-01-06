SRC= main.go build.go options.go tags.go
GOLANG=/usr/bin/go
CURL=/usr/bin/curl
GIT=/usr/bin/git
REPONAME=$(shell basename `git rev-parse --show-toplevel`)
SHA1=$(shell git rev-parse --short HEAD)
NOW=$(shell date +%Y-%m-%d_%T)

pws-exporter: ${SRC}
	echo ${REPONAME}
	${GOLANG} build -ldflags "-X main.sha1ver=${SHA1} -X main.buildTime=${NOW} -X main.repoName=${REPONAME}"

update-go:
	${GOLANG} get github.com/prometheus/client_golang/prometheus
	${GOLANG} get github.com/prometheus/client_golang/prometheus/promhttp
	${GOLANG} get"github.com/alecthomas/kong
	${GOLANG} get gopkg.in/yaml.v2

clean:
	rm -f pws-exporter

